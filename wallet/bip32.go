package wallet

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// hardenedKeyStart is the index at and above which BIP-32 child derivation is
// hardened.
const hardenedKeyStart uint32 = 0x80000000

// bip32MasterKey is the HMAC key used to derive the BIP-32 master key from a
// seed, as defined by BIP-32.
var bip32MasterKey = []byte("Bitcoin seed")

// hdKey is a minimal BIP-32 private extended key. It supports only private
// (CKDpriv) child derivation, which is all that is needed to derive Ethereum
// keys from a seed; it deliberately omits public derivation and the xprv/xpub
// serialization (and therefore any network parameters such as chaincfg).
//
// Reference: BIP-32
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
type hdKey struct {
	key       [32]byte // secp256k1 private scalar, big-endian.
	chainCode [32]byte
}

// newMasterKey derives the BIP-32 master key from the given seed.
func newMasterKey(seed []byte) (*hdKey, error) {
	// I = HMAC-SHA512(Key = "Bitcoin seed", Data = seed); IL is the master
	// key and IR is the master chain code.
	sum := hmacSHA512(bip32MasterKey, seed)

	// The master key must be a valid, non-zero secp256k1 scalar.
	var scalar secp256k1.ModNScalar
	if scalar.SetByteSlice(sum[:32]) {
		return nil, errors.New("bip32: invalid master key (>= curve order)")
	}
	if scalar.IsZero() {
		return nil, errors.New("bip32: invalid master key (zero)")
	}

	k := &hdKey{}
	copy(k.key[:], sum[:32])
	copy(k.chainCode[:], sum[32:])
	return k, nil
}

// derive derives the child key for the given index (BIP-32 CKDpriv). Indices
// at or above hardenedKeyStart perform hardened derivation.
func (k *hdKey) derive(index uint32) (*hdKey, error) {
	// Assemble the HMAC data.
	data := make([]byte, 0, 37)
	if index >= hardenedKeyStart {
		// Hardened: 0x00 || ser256(k_par) || ser32(index).
		data = append(data, 0x00)
		data = append(data, k.key[:]...)
	} else {
		// Normal: serP(point(k_par)) || ser32(index), i.e. the compressed
		// public key of the parent.
		priv := secp256k1.PrivKeyFromBytes(k.key[:])
		data = append(data, priv.PubKey().SerializeCompressed()...)
	}
	var idx [4]byte
	binary.BigEndian.PutUint32(idx[:], index)
	data = append(data, idx[:]...)

	// I = HMAC-SHA512(Key = c_par, Data = data); split into IL and IR.
	sum := hmacSHA512(k.chainCode[:], data)

	// child_key = (parse256(IL) + k_par) mod n. Per BIP-32 the key is invalid
	// if IL >= n or the result is zero (probability ~2^-127).
	var il, par secp256k1.ModNScalar
	if il.SetByteSlice(sum[:32]) {
		return nil, errors.New("bip32: invalid derived key (IL >= curve order)")
	}
	par.SetByteSlice(k.key[:])
	il.Add(&par)
	if il.IsZero() {
		return nil, errors.New("bip32: invalid derived key (zero)")
	}

	child := &hdKey{key: il.Bytes()}
	copy(child.chainCode[:], sum[32:])
	return child, nil
}

func hmacSHA512(key, data []byte) []byte {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
