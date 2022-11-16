package types

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/hexutil"
)

//
// Address type:
//

// AddressLength is the length of an Ethereum address in bytes.
const AddressLength = 20

// Address represents an Ethereum address encoded as a 20 byte array.
type Address [AddressLength]byte

// HexToAddress parses an address in hex format and returns an Address type.
func HexToAddress(address string) (a Address, err error) {
	err = a.UnmarshalText([]byte(address))
	return a, err
}

// HexToAddressPtr parses an address in hex format and returns an *Address type.
// It returns nil if the address is invalid.
func HexToAddressPtr(address string) *Address {
	a, err := HexToAddress(address)
	if err != nil {
		return nil
	}
	return &a
}

// MustHexToAddress parses an address in hex format and returns an Address type.
// It panics if the address is invalid.
func MustHexToAddress(address string) Address {
	a, err := HexToAddress(address)
	if err != nil {
		panic(err)
	}
	return a
}

// MustHexToAddressPtr parses an address in hex format and returns an *Address type.
// It panics if the address is invalid.
func MustHexToAddressPtr(address string) *Address {
	a := MustHexToAddress(address)
	return &a
}

// BytesToAddress converts a byte slice to an Address type.
func BytesToAddress(b []byte) (Address, error) {
	var a Address
	if len(b) != len(a) {
		return a, fmt.Errorf("invalid address length %d", len(b))
	}
	copy(a[:], b)
	return a, nil
}

// BytesToAddressPtr converts a byte slice to an *Address type.
// It returns nil if the address is invalid.
func BytesToAddressPtr(b []byte) *Address {
	a, err := BytesToAddress(b)
	if err != nil {
		return nil
	}
	return &a
}

// MustBytesToAddress converts a byte slice to an Address type.
// It panics if the address is invalid.
func MustBytesToAddress(b []byte) Address {
	a, err := BytesToAddress(b)
	if err != nil {
		panic(err)
	}
	return a
}

// MustBytesToAddressPtr converts a byte slice to an *Address type.
// It panics if the address is invalid.
func MustBytesToAddressPtr(b []byte) *Address {
	a := MustBytesToAddress(b)
	return &a
}

// Bytes returns the byte representation of the address.
func (t Address) Bytes() []byte {
	return t[:]
}

func (t Address) String() string {
	return hexutil.BytesToHex(t[:])
}

func (t Address) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *Address) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t Address) MarshalText() ([]byte, error) {
	return bytesMarshalText(t[:]), nil
}

func (t *Address) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

func (t Address) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.NewBytes(t[:]))
}

func (t *Address) DecodeRLP(data []byte) (int, error) {
	r, n, err := rlp.Decode(data)
	if err != nil {
		return 0, err
	}
	a, err := r.GetBytes()
	if err != nil {
		return 0, err
	}
	if len(a) != AddressLength {
		return 0, fmt.Errorf("invalid address length %d", len(a))
	}
	copy(t[:], a)
	return n, nil
}

//
// Hash type:
//

const HashLength = 32

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// HexToHash parses a hash in hex format and returns a Hash type.
func HexToHash(x string) (h Hash, err error) {
	err = h.UnmarshalText([]byte(x))
	return
}

// HexToHashPtr parses a hash in hex format and returns a *Hash type.
// It returns nil if the hash is invalid.
func HexToHashPtr(x string) *Hash {
	h, err := HexToHash(x)
	if err != nil {
		return nil
	}
	return &h
}

// MustHexToHash parses a hash in hex format and returns a Hash type.
// It panics if the hash is invalid.
func MustHexToHash(x string) Hash {
	h, err := HexToHash(x)
	if err != nil {
		panic(err)
	}
	return h
}

// MustHexToHashPtr parses a hash in hex format and returns a *Hash type.
// It panics if the hash is invalid.
func MustHexToHashPtr(x string) *Hash {
	h := MustHexToHash(x)
	return &h
}

// BytesToHash converts a byte slice to a Hash type.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it returns an error.
func BytesToHash(x []byte) (Hash, error) {
	var h Hash
	if len(x) > len(h) {
		return h, fmt.Errorf("invalid hash length %d", len(x))
	}
	copy(h[HashLength-len(x):], x)
	return h, nil
}

// BytesToHashPtr converts a byte slice to a *Hash type.
// It returns nil if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it returns nil.
func BytesToHashPtr(x []byte) *Hash {
	h, err := BytesToHash(x)
	if err != nil {
		return nil
	}
	return &h
}

// MustBytesToHash converts a byte slice to a Hash type.
// It panics if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it panics.
func MustBytesToHash(x []byte) Hash {
	h, err := BytesToHash(x)
	if err != nil {
		panic(err)
	}
	return h
}

// MustBytesToHashPtr converts a byte slice to a *Hash type.
// It panics if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it panics.
func MustBytesToHashPtr(x []byte) *Hash {
	h := MustBytesToHash(x)
	return &h
}

func (t Hash) Bytes() []byte {
	return t[:]
}

func (t Hash) String() string {
	return hexutil.BytesToHex(t[:])
}

func (t Hash) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *Hash) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t Hash) MarshalText() ([]byte, error) {
	return bytesMarshalText(t[:]), nil
}

func (t *Hash) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

func (t Hash) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.NewBytes(t[:]))
}

func (t *Hash) DecodeRLP(data []byte) (int, error) {
	r, n, err := rlp.Decode(data)
	if err != nil {
		return 0, err
	}
	b, err := r.GetBytes()
	if err != nil {
		return 0, err
	}
	if len(b) != HashLength {
		return 0, fmt.Errorf("invalid hash length %d", len(t))
	}
	copy(t[:], b)
	return n, nil
}

//
// BlockNumber type:
//

// BlockNumber is a type that can hold a block number or a tag.
type BlockNumber struct{ x big.Int }

const (
	earliestBlockNumber = -1
	latestBlockNumber   = -2
	pendingBlockNumber  = -3
)

var (
	EarliestBlockNumber = BlockNumber{x: *new(big.Int).SetInt64(earliestBlockNumber)}
	LatestBlockNumber   = BlockNumber{x: *new(big.Int).SetInt64(latestBlockNumber)}
	PendingBlockNumber  = BlockNumber{x: *new(big.Int).SetInt64(pendingBlockNumber)}
)

// HexToBlockNumber converts a string to a BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// If the string is not a valid block number, it returns an error.
func HexToBlockNumber(x string) (BlockNumber, error) {
	b := &BlockNumber{}
	err := b.UnmarshalText([]byte(x))
	return *b, err
}

// HexToBlockNumberPtr converts a string to a *BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// If the string is not a valid block number, it returns nil.
func HexToBlockNumberPtr(x string) *BlockNumber {
	b, err := HexToBlockNumber(x)
	if err != nil {
		return nil
	}
	return &b
}

// MustHexToBlockNumber converts a string to a BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// It panics if the string is not a valid block number.
func MustHexToBlockNumber(x string) BlockNumber {
	b, err := HexToBlockNumber(x)
	if err != nil {
		panic(err)
	}
	return b
}

// Uint64ToBlockNumber converts an uint64 to a BlockNumber type.
func Uint64ToBlockNumber(x uint64) BlockNumber {
	return BlockNumber{x: *new(big.Int).SetUint64(x)}
}

// Uint64ToBlockNumberPtr converts an uint64 to a *BlockNumber type.
func Uint64ToBlockNumberPtr(x uint64) *BlockNumber {
	b := Uint64ToBlockNumber(x)
	return &b
}

// BigIntToBlockNumber converts a big.Int to a BlockNumber type.
func BigIntToBlockNumber(x *big.Int) BlockNumber {
	if x == nil {
		return BlockNumber{}
	}
	return BlockNumber{x: *new(big.Int).Set(x)}
}

// BigIntToBlockNumberPtr converts a big.Int to a *BlockNumber type.
func BigIntToBlockNumberPtr(x *big.Int) *BlockNumber {
	b := BigIntToBlockNumber(x)
	return &b
}

// IsEarliest returns true if the block tag is "earliest".
func (t *BlockNumber) IsEarliest() bool {
	return t.Big().Int64() == earliestBlockNumber
}

// IsLatest returns true if the block tag is "latest".
func (t *BlockNumber) IsLatest() bool {
	return t.Big().Int64() == latestBlockNumber
}

// IsPending returns true if the block tag is "pending".
func (t *BlockNumber) IsPending() bool {
	return t.Big().Int64() == pendingBlockNumber
}

// IsTag returns true if the block tag is used.
func (t *BlockNumber) IsTag() bool {
	return t.Big().Sign() < 0
}

// Big returns the big.Int representation of the block number.
func (t *BlockNumber) Big() *big.Int {
	return new(big.Int).Set(&t.x)
}

func (t *BlockNumber) String() string {
	switch {
	case t.IsEarliest():
		return "earliest"
	case t.IsLatest():
		return "latest"
	case t.IsPending():
		return "pending"
	default:
		return "0x" + t.x.Text(16)
	}
}

func (t BlockNumber) MarshalJSON() ([]byte, error) {
	b, err := t.MarshalText()
	if err != nil {
		return nil, err
	}
	return naiveQuote(b), nil
}

func (t *BlockNumber) UnmarshalJSON(input []byte) error {
	return t.UnmarshalText(naiveUnquote(input))
}

func (t BlockNumber) MarshalText() ([]byte, error) {
	switch {
	case t.IsEarliest():
		return []byte("earliest"), nil
	case t.IsLatest():
		return []byte("latest"), nil
	case t.IsPending():
		return []byte("pending"), nil
	default:
		return []byte(hexutil.BigIntToHex(&t.x)), nil
	}
}

func (t *BlockNumber) UnmarshalText(input []byte) error {
	switch strings.ToLower(strings.TrimSpace(string(input))) {
	case "earliest":
		*t = BlockNumber{x: *new(big.Int).SetInt64(earliestBlockNumber)}
		return nil
	case "latest":
		*t = BlockNumber{x: *new(big.Int).SetInt64(latestBlockNumber)}
		return nil
	case "pending":
		*t = BlockNumber{x: *new(big.Int).SetInt64(pendingBlockNumber)}
		return nil
	default:
		u, err := hexutil.HexToBigInt(string(input))
		if err != nil {
			return err
		}
		if u.Cmp(big.NewInt(math.MaxInt64)) > 0 {
			return fmt.Errorf("block number larger than int64")
		}
		*t = BlockNumber{x: *u}
		return nil
	}
}

//
// Signature type:
//

// SignatureLength is the expected length of the Signature.
const SignatureLength = 65

// Signature represents the 65 byte signature.
type Signature [SignatureLength]byte

// HexToSignature parses a hex string into a Signature.
func HexToSignature(x string) Signature {
	var s Signature
	_ = s.UnmarshalText([]byte(x))
	return s
}

// BytesToSignature returns Signature from bytes.
func BytesToSignature(b []byte) Signature {
	var sig Signature
	if len(b) != SignatureLength {
		return sig
	}
	copy(sig[:], b)
	return sig
}

// VRSToSignature returns Signature from VRS values.
func VRSToSignature(v uint8, r [32]byte, s [32]byte) Signature {
	return BytesToSignature(append(append(append([]byte{}, r[:]...), s[:]...), v))
}

func BigIntToSignature(v, r, s *big.Int) (Signature, error) {
	var sig Signature
	vb := v.Bytes()
	rb := r.Bytes()
	sb := s.Bytes()
	if len(vb) > 1 || len(rb) > 32 || len(sb) > 32 {
		return sig, fmt.Errorf("invalid signature")
	}
	copy(sig[64:65], vb)
	copy(sig[32-len(rb):32], rb)
	copy(sig[64-len(sb):64], sb)
	return sig, nil
}

// VRS returns the V, R, S values of the signature.
func (s Signature) VRS() (sv uint8, sr [32]byte, ss [32]byte) {
	copy(sr[:], s[:32])
	copy(ss[:], s[32:64])
	sv = s[64]
	return
}

func (s Signature) V() uint8 {
	return s[64]
}

func (s Signature) R() (r [32]byte) {
	copy(r[:], s[:32])
	return r
}

func (s Signature) S() (s2 [32]byte) {
	copy(s2[:], s[32:64])
	return s2
}

func (s Signature) BigV() *big.Int {
	return big.NewInt(int64(s[64]))
}

func (s Signature) BigR() *big.Int {
	return new(big.Int).SetBytes(s[:32])
}

func (s Signature) BigS() *big.Int {
	return new(big.Int).SetBytes(s[32:64])
}

// Bytes returns the byte representation of the signature. .
func (s Signature) Bytes() []byte {
	return s[:]
}

// String returns the hex representation of the signature.
func (s Signature) String() string {
	return hexutil.BytesToHex(s[:])
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(s[:]), nil
}

func (s *Signature) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, s[:])
}

func (s Signature) MarshalText() ([]byte, error) {
	return bytesMarshalText(s[:]), nil
}

func (s *Signature) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, s[:])
}

//
// Number type:
//

// Number represents a hex-encoded number. This type is used for marshaling
// and unmarshalling JSON numbers. When possible, use big.Int or regular integers
// instead.
type Number struct{ x big.Int }

// HexToNumber converts a hex string to a Number type.
func HexToNumber(x string) Number {
	u, _ := hexutil.HexToBigInt(x)
	return Number{x: *u}
}

// HexToNumberPtr converts a hex string to a *Number type.
func HexToNumberPtr(x string) *Number {
	n := HexToNumber(x)
	return &n
}

// Uint64ToNumber converts an uint64 to a Number type.
func Uint64ToNumber(x uint64) Number {
	return Number{x: *new(big.Int).SetUint64(x)}
}

// Uint64ToNumberPtr converts an uint64 to a *Number type.
func Uint64ToNumberPtr(x uint64) *Number {
	n := Uint64ToNumber(x)
	return &n
}

// BigIntToNumber converts a big.Int to a Number type.
func BigIntToNumber(x *big.Int) Number {
	if x == nil {
		return Number{}
	}
	return Number{x: *x}
}

// BigIntToNumberPtr converts a big.Int to a *Number type.
func BigIntToNumberPtr(x *big.Int) *Number {
	n := BigIntToNumber(x)
	return &n
}

func (t *Number) Big() *big.Int {
	return new(big.Int).Set(&t.x)
}

func (t *Number) Bytes() []byte {
	return t.x.Bytes()
}

func (t *Number) String() string {
	return hexutil.BigIntToHex(&t.x)
}

func (t Number) MarshalJSON() ([]byte, error) {
	return numberMarshalJSON(t.Big()), nil
}

func (t *Number) UnmarshalJSON(input []byte) error {
	return numberUnmarshalJSON(input, &t.x)
}

func (t Number) MarshalText() ([]byte, error) {
	return numberMarshalText(t.Big()), nil
}

func (t *Number) UnmarshalText(input []byte) error {
	return numberUnmarshalText(input, &t.x)
}

//
// Bytes type:
//

// Bytes represents a hex-encoded byte slice. This type is used for marshaling
// and unmarshalling JSON numbers. When possible, use byte slices instead.
type Bytes []byte

func (t *Bytes) Bytes() []byte {
	return *t
}

func (t *Bytes) String() string {
	if t == nil {
		return ""
	}
	return hexutil.BytesToHex(*t)
}

func (t Bytes) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t), nil
}

func (t *Bytes) UnmarshalJSON(input []byte) error {
	return bytesUnmarshalJSON(input, (*[]byte)(t))
}

func (t Bytes) MarshalText() ([]byte, error) {
	return bytesMarshalText(t), nil
}

func (t *Bytes) UnmarshalText(input []byte) error {
	return bytesUnmarshalText(input, (*[]byte)(t))
}

//
// Internal types:
//

const bloomLength = 256

type hexBloom [bloomLength]byte

func bloomFromBytes(x []byte) hexBloom {
	var b [bloomLength]byte
	if len(x) > len(b) {
		return b
	}
	copy(b[bloomLength-len(x):], x)
	return b
}

func (t *hexBloom) Bytes() []byte {
	return t[:]
}

func (t *hexBloom) String() string {
	if t == nil {
		return ""
	}
	return hexutil.BytesToHex(t[:])
}

func (t hexBloom) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *hexBloom) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t hexBloom) MarshalText() ([]byte, error) {
	return bytesMarshalText(t[:]), nil
}

func (t *hexBloom) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

const nonceLength = 8

type hexNonce [nonceLength]byte

func nonceFromBigInt(x *big.Int) hexNonce {
	if x == nil {
		return hexNonce{}
	}
	return nonceFromBytes(x.Bytes())
}

func nonceFromBytes(x []byte) hexNonce {
	var n hexNonce
	if len(x) > len(n) {
		return n
	}
	copy(n[nonceLength-len(x):], x)
	return n
}

func (t *hexNonce) Big() *big.Int {
	return new(big.Int).SetBytes(t[:])
}

func (t *hexNonce) String() string {
	if t == nil {
		return ""
	}
	return hexutil.BytesToHex(t[:])
}

func (t hexNonce) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *hexNonce) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t hexNonce) MarshalText() ([]byte, error) {
	return bytesMarshalText(t[:]), nil
}

func (t *hexNonce) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

type hashList []Hash

func (b hashList) MarshalJSON() ([]byte, error) {
	if len(b) == 1 {
		return json.Marshal(b[0])
	}
	return json.Marshal([]Hash(b))
}

func (b *hashList) UnmarshalJSON(input []byte) error {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		*b = hashList{{}}
		return json.Unmarshal(input, &((*b)[0]))
	}
	return json.Unmarshal(input, (*[]Hash)(b))
}

type addressList []Address

func (t addressList) MarshalJSON() ([]byte, error) {
	if len(t) == 1 {
		return json.Marshal(t[0])
	}
	return json.Marshal([]Address(t))
}

func (t *addressList) UnmarshalJSON(input []byte) error {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		*t = addressList{{}}
		return json.Unmarshal(input, &((*t)[0]))
	}
	return json.Unmarshal(input, (*[]Address)(t))
}
