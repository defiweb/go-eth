package wallet

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"

	"github.com/defiweb/go-eth/crypto"
)

// The code below is based on github.com/ethereum/go-ethereum/tree/master/accounts/keystore

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

	// Verify the derived key matches the key in the JSON. If not, the
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
	salt := cryptoJSON.KDFParams.Salt
	dkLen := cryptoJSON.KDFParams.DKLen

	switch cryptoJSON.KDF {
	case "scrypt":
		n := cryptoJSON.KDFParams.N
		r := cryptoJSON.KDFParams.R
		p := cryptoJSON.KDFParams.P
		return scrypt.Key(passphrase, salt, n, r, p, dkLen)
	case "pbkdf2":
		c := cryptoJSON.KDFParams.C
		prf := cryptoJSON.KDFParams.PRF
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(passphrase, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}

// aesCTRXOR performs AES-128-CTR decryption on the given cipher text with the
// given key and IV.
func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}
