package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
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

// Upper bounds on the KDF parameters accepted from a keystore file.
const (
	maxScryptN     = 1 << 22
	maxScryptR     = 32
	maxScryptP     = 16
	maxPBKDF2Count = 10_000_000
)

func encryptV3Key(key crypto.PrivateKey, passphrase string, scryptN, scryptP int) (*jsonKey, error) {
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
	cipherText, err := aesCTRXOR(derivedKey[:16], key.Bytes(), iv)
	if err != nil {
		return nil, err
	}

	// Calculate the MAC of the encrypted key.
	mac := crypto.Keccak256(derivedKey[16:32], cipherText)

	// Generate a random UUID for the keyfile.
	id, err := randUUID()
	if err != nil {
		return nil, err
	}

	// Assemble and return the key JSON.
	return &jsonKey{
		Version: 3,
		ID:      id,
		Address: types.Address(crypto.ECPublicKeyToAddress(crypto.ECPrivateKeyToPublicKey(key))),
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
			MAC: mac[:],
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

	// Verify that the derived key matches the key in the JSON. If not, the
	// passphrase is incorrect.
	calculatedMAC := crypto.Keccak256(derivedKey[16:32], cryptoJson.CipherText)
	if subtle.ConstantTimeCompare(calculatedMAC[:], cryptoJson.MAC) != 1 {
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
	if cryptoJSON.KDFParams.DKLen != scryptDKLen {
		return nil, fmt.Errorf("invalid KDF key length: got %d, want %d", cryptoJSON.KDFParams.DKLen, scryptDKLen)
	}
	switch cryptoJSON.KDF {
	case "scrypt":
		switch {
		case cryptoJSON.KDFParams.N <= 1 || cryptoJSON.KDFParams.N > maxScryptN:
			return nil, fmt.Errorf("invalid scrypt N: %d", cryptoJSON.KDFParams.N)
		case cryptoJSON.KDFParams.R <= 0 || cryptoJSON.KDFParams.R > maxScryptR:
			return nil, fmt.Errorf("invalid scrypt r: %d", cryptoJSON.KDFParams.R)
		case cryptoJSON.KDFParams.P <= 0 || cryptoJSON.KDFParams.P > maxScryptP:
			return nil, fmt.Errorf("invalid scrypt p: %d", cryptoJSON.KDFParams.P)
		}
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
		if cryptoJSON.KDFParams.C <= 0 || cryptoJSON.KDFParams.C > maxPBKDF2Count {
			return nil, fmt.Errorf("invalid PBKDF2 iteration count: %d", cryptoJSON.KDFParams.C)
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

func randUUID() (string, error) {
	var uuid [16]byte
	var text [36]byte
	if _, err := io.ReadFull(rand.Reader, uuid[:]); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	hex.Encode(text[:8], uuid[:4])
	text[8] = '-'
	hex.Encode(text[9:13], uuid[4:6])
	text[13] = '-'
	hex.Encode(text[14:18], uuid[6:8])
	text[18] = '-'
	hex.Encode(text[19:23], uuid[8:10])
	text[23] = '-'
	hex.Encode(text[24:], uuid[10:])
	return string(text[:]), nil
}
