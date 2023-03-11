package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/defiweb/go-eth/types"
)

var ErrKeyNotFound = errors.New("key not found")

// NewKeyFromJSON loads an Ethereum key from a JSON file.
func NewKeyFromJSON(path string, passphrase string) (*PrivateKey, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewKeyFromJSONContent(content, passphrase)
}

// NewKeyFromJSONContent returns a new key from a JSON.
func NewKeyFromJSONContent(content []byte, passphrase string) (*PrivateKey, error) {
	var jKey jsonKey
	if err := json.Unmarshal(content, &jKey); err != nil {
		return nil, err
	}
	if jKey.Version != 3 {
		return nil, errors.New("only V3 keys are supported")
	}
	prv, err := decryptV3Key(jKey.Crypto, []byte(passphrase))
	if err != nil {
		return nil, err
	}
	key := NewKeyFromBytes(prv)
	if !jKey.Address.IsZero() && jKey.Address != key.Address() {
		return nil, errors.New("decrypted key address does not match address in file")
	}
	return key, nil
}

// NewKeyFromDirectory returns a new key from a directory containing JSON
// files.
func NewKeyFromDirectory(path string, passphrase string, address types.Address) (*PrivateKey, error) {
	items, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IsDir() {
			// Skip directories.
			continue
		}
		if item.Size() == 0 || item.Size() > 1<<20 {
			// Skip empty files and files larger than 1MB.
			continue
		}
		key, err := NewKeyFromJSON(filepath.Join(path, item.Name()), passphrase)
		if err != nil {
			// Skip files that are not keys or have invalid content.
			continue
		}
		if address == key.Address() {
			return key, nil
		}
	}
	return nil, ErrKeyNotFound
}

type jsonKey struct {
	ID      string        `json:"id"`
	Version int64         `json:"version"`
	Address types.Address `json:"address"`
	Crypto  jsonKeyCrypto `json:"crypto"`
}

type jsonKeyCrypto struct {
	Cipher       string              `json:"cipher"`
	CipherText   jsonHex             `json:"ciphertext"`
	CipherParams jsonKeyCipherParams `json:"cipherparams"`
	KDF          string              `json:"kdf"`
	KDFParams    jsonKeyKDFParams    `json:"kdfparams"`
	MAC          jsonHex             `json:"mac"`
}

type jsonKeyCipherParams struct {
	IV jsonHex `json:"iv"`
}

type jsonKeyKDFParams struct {
	DKLen int     `json:"dklen"`
	Salt  jsonHex `json:"salt"`

	// Scrypt params:
	N int `json:"n"`
	P int `json:"p"`
	R int `json:"r"`

	// PBKDF2 params:
	C   int    `json:"c"`
	PRF string `json:"prf"`
}

type jsonHex []byte

func (h jsonHex) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString(h) + `"`), nil
}

func (h *jsonHex) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid hex string")
	}
	*h, err = hex.DecodeString(string(data[1 : len(data)-1]))
	return
}
