package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	"github.com/defiweb/go-eth/crypto"
)

// The code below is based on:
// github.com/ethereum/go-ethereum/tree/master/accounts/keystore

const (
	StandardScryptN = 1 << 18
	StandardScryptP = 1
	LightScryptN    = 1 << 12
	LightScryptP    = 6
	scryptR         = 8
	scryptDKLen     = 32
)

func encryptV3Key(key *ecdsa.PrivateKey, passphrase string, scryptN, scryptP int) (*jsonKey, error) {
	// Generate a random salt.
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive the key from the passphrase.
	derivedKey, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}

	// Generate a random IV.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// Encrypt the key with AES-128-CTR.
	d := key.D.Bytes()
	data := make([]byte, 32)
	copy(data[32-len(d):], d)
	cipherText, err := aesCTRXOR(derivedKey[:16], data, iv)
	if err != nil {
		return nil, err
	}

	// Calculate the MAC of the encrypted key.
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	// Generate a random UUID for the keyfile.
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Assemble and return the key JSON.
	return &jsonKey{
		Version: 3,
		ID:      id.String(),
		Address: crypto.ECPublicKeyToAddress(&key.PublicKey),
		Crypto: jsonKeyCrypto{
			Cipher: "aes-128-ctr",
			CipherParams: jsonKeyCipherParams{
				IV: iv,
			},
			CipherText: cipherText,
			KDF:        "scrypt",
			KDFParams: jsonKeyKDFParams{
				DKLen: scryptDKLen,
				N:     scryptN,
				P:     scryptP,
				R:     scryptR,
				Salt:  salt,
			},
			MAC: mac.Bytes(),
		},
	}, nil
}

// decryptKey decrypts the given V3 key with the given passphrase.
func decryptV3Key(cryptoJson jsonKeyCrypto, passphrase []byte) ([]byte, error) {
	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", cryptoJson.Cipher)
	}

	// Derive the key from the passphrase.
	derivedKey, err := deriveKey(cryptoJson, passphrase)
	if err != nil {
		return nil, err
	}

	// VerifyHash the derived key matches the key in the JSON. If not, the
	// passphrase is incorrect.
	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cryptoJson.CipherText)
	if !bytes.Equal(calculatedMAC.Bytes(), cryptoJson.MAC) {
		return nil, fmt.Errorf("invalid passphrase or keyfile")
	}

	// Decrypt the key with AES-128-CTR.
	plainText, err := aesCTRXOR(derivedKey[:16], cryptoJson.CipherText, cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	return plainText, err
}

// deriveKey returns the derived key from the JSON keyfile.
func deriveKey(cryptoJSON jsonKeyCrypto, passphrase []byte) ([]byte, error) {
	switch cryptoJSON.KDF {
	case "scrypt":
		return scrypt.Key(
			passphrase,
			cryptoJSON.KDFParams.Salt,
			cryptoJSON.KDFParams.N,
			cryptoJSON.KDFParams.R,
			cryptoJSON.KDFParams.P,
			cryptoJSON.KDFParams.DKLen,
		)
	case "pbkdf2":
		if cryptoJSON.KDFParams.PRF != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", cryptoJSON.KDFParams.PRF)
		}
		key := pbkdf2.Key(
			passphrase,
			cryptoJSON.KDFParams.Salt,
			cryptoJSON.KDFParams.C,
			cryptoJSON.KDFParams.DKLen,
			sha256.New,
		)
		return key, nil
	}
	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}

// aesCTRXOR performs AES-128-CTR decryption on the given cipher text with the
// given key and IV.
func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}
