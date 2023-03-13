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

// HashFunc returns the hash for the given input.
type HashFunc func(data ...[]byte) Hash

//
// Address type:
//

// AddressLength is the length of an Ethereum address in bytes.
const AddressLength = 20

// Address represents an Ethereum address encoded as a 20 byte array.
type Address [AddressLength]byte

// ZeroAddress is an address with all zeros.
var ZeroAddress = Address{}

// AddressFromHex parses an address in hex format and returns an Address type.
func AddressFromHex(h string) (a Address, err error) {
	err = a.UnmarshalText([]byte(h))
	return a, err
}

// AddressFromHexPtr parses an address in hex format and returns an *Address type.
// It returns nil if the address is invalid.
func AddressFromHexPtr(h string) *Address {
	a, err := AddressFromHex(h)
	if err != nil {
		return nil
	}
	return &a
}

// MustAddressFromHex parses an address in hex format and returns an Address type.
// It panics if the address is invalid.
func MustAddressFromHex(h string) Address {
	a, err := AddressFromHex(h)
	if err != nil {
		panic(err)
	}
	return a
}

// MustAddressFromHexPtr parses an address in hex format and returns an *Address type.
// It panics if the address is invalid.
func MustAddressFromHexPtr(h string) *Address {
	a := MustAddressFromHex(h)
	return &a
}

// AddressFromBytes converts a byte slice to an Address type.
func AddressFromBytes(b []byte) (Address, error) {
	var a Address
	if len(b) != len(a) {
		return a, fmt.Errorf("invalid address length %d", len(b))
	}
	copy(a[:], b)
	return a, nil
}

// AddressFromBytesPtr converts a byte slice to an *Address type.
// It returns nil if the address is invalid.
func AddressFromBytesPtr(b []byte) *Address {
	a, err := AddressFromBytes(b)
	if err != nil {
		return nil
	}
	return &a
}

// MustAddressFromBytes converts a byte slice to an Address type.
// It panics if the address is invalid.
func MustAddressFromBytes(b []byte) Address {
	a, err := AddressFromBytes(b)
	if err != nil {
		panic(err)
	}
	return a
}

// MustAddressFromBytesPtr converts a byte slice to an *Address type.
// It panics if the address is invalid.
func MustAddressFromBytesPtr(b []byte) *Address {
	a := MustAddressFromBytes(b)
	return &a
}

// Bytes returns the byte representation of the address.
func (t Address) Bytes() []byte {
	return t[:]
}

// String returns the hex representation of the address.
func (t Address) String() string {
	return hexutil.BytesToHex(t[:])
}

// Checksum returns the address with the checksum calculated according to
// EIP-55.
func (t Address) Checksum(h HashFunc) string {
	hex := []byte(hexutil.BytesToHex(t[:])[2:])
	hash := h(hex)
	for i, c := range hex {
		if c >= '0' && c <= '9' {
			continue
		}
		if hash[i/2]&(uint8(1)<<(((i+1)%2)*4+3)) != 0 {
			hex[i] = c ^ 0x20
		}
	}
	return "0x" + string(hex)
}

// IsZero returns true if the address is the zero address.
func (t Address) IsZero() bool {
	return t == ZeroAddress
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
	if len(a) == 0 {
		*t = ZeroAddress
		return n, nil
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

// ZeroHash is a hash with all zeros.
var ZeroHash = Hash{}

// HashFromHex parses a hash in hex format and returns a Hash type.
func HashFromHex(h string) (Hash, error) {
	var hash Hash
	err := hash.UnmarshalText([]byte(h))
	return hash, err
}

// HashFromHexPtr parses a hash in hex format and returns a *Hash type.
// It returns nil if the hash is invalid.
func HashFromHexPtr(h string) *Hash {
	hash, err := HashFromHex(h)
	if err != nil {
		return nil
	}
	return &hash
}

// MustHashFromHex parses a hash in hex format and returns a Hash type.
// It panics if the hash is invalid.
func MustHashFromHex(h string) Hash {
	hash, err := HashFromHex(h)
	if err != nil {
		panic(err)
	}
	return hash
}

// MustHashFromHexPtr parses a hash in hex format and returns a *Hash type.
// It panics if the hash is invalid.
func MustHashFromHexPtr(h string) *Hash {
	hash := MustHashFromHex(h)
	return &hash
}

// HashFromBytes converts a byte slice to a Hash type.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it returns an error.
func HashFromBytes(b []byte) (Hash, error) {
	var h Hash
	if len(b) > len(h) {
		return h, fmt.Errorf("invalid hash length %d", len(b))
	}
	copy(h[HashLength-len(b):], b)
	return h, nil
}

// HashFromBytesPtr converts a byte slice to a *Hash type.
// It returns nil if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it returns nil.
func HashFromBytesPtr(b []byte) *Hash {
	h, err := HashFromBytes(b)
	if err != nil {
		return nil
	}
	return &h
}

// MustHashFromBytes converts a byte slice to a Hash type.
// It panics if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it panics.
func MustHashFromBytes(b []byte) Hash {
	h, err := HashFromBytes(b)
	if err != nil {
		panic(err)
	}
	return h
}

// MustHashFromBytesPtr converts a byte slice to a *Hash type.
// It panics if the hash is invalid.
// If bytes is shorter than 32 bytes, it left-pads Hash with zeros.
// If bytes is longer than 32 bytes, it panics.
func MustHashFromBytesPtr(b []byte) *Hash {
	h := MustHashFromBytes(b)
	return &h
}

// HashFromBigInt converts a big.Int to a Hash type.
// Negative numbers are represented as two's complement.
func HashFromBigInt(i *big.Int) (Hash, error) {
	var b []byte
	if i.Sign() >= 0 {
		b = i.Bytes()
	} else {
		m := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(HashLength*8)), big.NewInt(1))
		x := new(big.Int).Set(i).And(i, m)
		b = x.Bytes()
		if len(b) != HashLength || b[0]&0x80 == 0 {
			return Hash{}, fmt.Errorf("number too large to convert to hash")
		}
	}
	if len(b) > HashLength {
		return Hash{}, fmt.Errorf("number too large to convert to hash")
	}
	return HashFromBytes(b)
}

// HashFromBigIntPtr converts a big.Int to a *Hash type.
// Negative numbers are represented as two's complement.
// It returns nil if the hash is invalid.
func HashFromBigIntPtr(i *big.Int) *Hash {
	h, err := HashFromBigInt(i)
	if err != nil {
		return nil
	}
	return &h
}

// MustHashFromBigInt converts a big.Int to a Hash type.
// Negative numbers are represented as two's complement.
// It panics if the hash is invalid.
func MustHashFromBigInt(i *big.Int) Hash {
	h, err := HashFromBigInt(i)
	if err != nil {
		panic(err)
	}
	return h
}

// MustHashFromBigIntPtr converts a big.Int to a *Hash type.
// Negative numbers are represented as two's complement.
// It panics if the hash is invalid.
func MustHashFromBigIntPtr(i *big.Int) *Hash {
	h := MustHashFromBigInt(i)
	return &h
}

// Bytes returns hash as a byte slice.
func (t Hash) Bytes() []byte {
	return t[:]
}

// String returns the hex string representation of the hash.
func (t Hash) String() string {
	return hexutil.BytesToHex(t[:])
}

// IsZero returns true if the hash is the zero hash.
func (t Hash) IsZero() bool {
	return t == ZeroHash
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

// BlockNumberFromHex converts a string to a BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// If the string is not a valid block number, it returns an error.
func BlockNumberFromHex(h string) (BlockNumber, error) {
	b := &BlockNumber{}
	err := b.UnmarshalText([]byte(h))
	return *b, err
}

// BlockNumberFromHexPtr converts a string to a *BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// If the string is not a valid block number, it returns nil.
func BlockNumberFromHexPtr(h string) *BlockNumber {
	b, err := BlockNumberFromHex(h)
	if err != nil {
		return nil
	}
	return &b
}

// MustBlockNumberFromHex converts a string to a BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// It panics if the string is not a valid block number.
func MustBlockNumberFromHex(h string) BlockNumber {
	b, err := BlockNumberFromHex(h)
	if err != nil {
		panic(err)
	}
	return b
}

// MustBlockNumberFromHexPtr converts a string to a *BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "pending".
// It panics if the string is not a valid block number.
func MustBlockNumberFromHexPtr(h string) *BlockNumber {
	b := MustBlockNumberFromHex(h)
	return &b
}

// BlockNumberFromUint64 converts an uint64 to a BlockNumber type.
func BlockNumberFromUint64(x uint64) BlockNumber {
	return BlockNumber{x: *new(big.Int).SetUint64(x)}
}

// BlockNumberFromUint64Ptr converts an uint64 to a *BlockNumber type.
func BlockNumberFromUint64Ptr(x uint64) *BlockNumber {
	b := BlockNumberFromUint64(x)
	return &b
}

// BlockNumberFromBigInt converts a big.Int to a BlockNumber type.
func BlockNumberFromBigInt(x *big.Int) BlockNumber {
	if x == nil {
		return BlockNumber{}
	}
	return BlockNumber{x: *new(big.Int).Set(x)}
}

// BlockNumberFromBigIntPtr converts a big.Int to a *BlockNumber type.
func BlockNumberFromBigIntPtr(x *big.Int) *BlockNumber {
	b := BlockNumberFromBigInt(x)
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

// String returns the string representation of the block number.
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

// Signature represents the transaction signature.
type Signature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}

// SignatureFromHex parses a hex string into a Signature.
// Hex representation of the signature is hex([R || S || V]).
func SignatureFromHex(h string) (Signature, error) {
	b, err := hexutil.HexToBytes(h)
	if err != nil {
		return Signature{}, err
	}
	return SignatureFromBytes(b)
}

// SignatureFromHexPtr parses a hex string into a *Signature.
// Hex representation of the signature is hex([R || S || V]).
// It returns nil if the string is not a valid signature.
func SignatureFromHexPtr(h string) *Signature {
	sig, err := SignatureFromHex(h)
	if err != nil {
		return nil
	}
	return &sig
}

// MustSignatureFromHex parses a hex string into a Signature.
// Hex representation of the signature is hex([R || S || V]).
// It panics if the string is not a valid signature.
func MustSignatureFromHex(h string) Signature {
	sig, err := SignatureFromHex(h)
	if err != nil {
		panic(err)
	}
	return sig
}

// MustSignatureFromHexPtr parses a hex string into a *Signature.
// Hex representation of the signature is hex([R || S || V]).
// It panics if the string is not a valid signature.
func MustSignatureFromHexPtr(h string) *Signature {
	sig, err := SignatureFromHex(h)
	if err != nil {
		panic(err)
	}
	return &sig
}

// SignatureFromBytes returns Signature from bytes.
// Byte representation of the signature is [R || S || V].
func SignatureFromBytes(b []byte) (Signature, error) {
	if len(b) < 65 {
		return Signature{}, fmt.Errorf("signature too short")
	}
	return Signature{
		V: new(big.Int).SetBytes(b[64:]),
		R: new(big.Int).SetBytes(b[:32]),
		S: new(big.Int).SetBytes(b[32:64]),
	}, nil
}

// SignatureFromBytesPtr returns *Signature from bytes.
// Byte representation of the signature is [R || S || V].
// It returns nil if the length of the bytes is not 65.
func SignatureFromBytesPtr(b []byte) *Signature {
	sig, err := SignatureFromBytes(b)
	if err != nil {
		return nil
	}
	return &sig
}

// MustSignatureFromBytes returns Signature from bytes.
// Byte representation of the signature is [R || S || V].
// It panics if the length of the bytes is not 65.
func MustSignatureFromBytes(b []byte) Signature {
	sig, err := SignatureFromBytes(b)
	if err != nil {
		panic(err)
	}
	return sig
}

// MustSignatureFromBytesPtr returns *Signature from bytes.
// Byte representation of the signature is [R || S || V].
// It panics if the length of the bytes is not 65.
func MustSignatureFromBytesPtr(b []byte) *Signature {
	sig, err := SignatureFromBytes(b)
	if err != nil {
		panic(err)
	}
	return &sig
}

// SignatureFromVRS returns Signature from V, R, S values.
func SignatureFromVRS(v, r, s *big.Int) Signature {
	return Signature{
		V: v,
		R: r,
		S: s,
	}
}

// SignatureFromVRSPtr returns *Signature from V, R, S values.
func SignatureFromVRSPtr(v, r, s *big.Int) *Signature {
	sig := SignatureFromVRS(v, r, s)
	return &sig
}

// Bytes returns the byte representation of the signature.
// The byte representation is [R || S || V].
func (s Signature) Bytes() []byte {
	sv, sr, ss := s.V, s.R, s.S
	if sv == nil {
		sv = new(big.Int)
	}
	if sr == nil {
		sr = new(big.Int)
	}
	if ss == nil {
		ss = new(big.Int)
	}
	vb := sv.Bytes()
	if len(vb) == 0 {
		vb = []byte{0}
	}
	b := make([]byte, 64+len(vb))
	sr.FillBytes(b[:32])
	ss.FillBytes(b[32:64])
	copy(b[64:], vb)
	return b
}

// String returns the hex representation of the signature.
// The hex representation is hex([R || S || V]).
func (s Signature) String() string {
	return hexutil.BytesToHex(s.Bytes())
}

// IsZero returns true if the signature is zero.
func (s Signature) IsZero() bool {
	if s.V != nil && s.V.Sign() != 0 {
		return false
	}
	if s.R != nil && s.R.Sign() != 0 {
		return false
	}
	if s.S != nil && s.S.Sign() != 0 {
		return false
	}
	return true
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(s.Bytes()), nil
}

func (s *Signature) UnmarshalJSON(input []byte) error {
	var b []byte
	if err := bytesUnmarshalJSON(input, &b); err != nil {
		return err
	}
	sig, err := SignatureFromBytes(b)
	if err != nil {
		return err
	}
	*s = sig
	return nil
}

func (s Signature) MarshalText() ([]byte, error) {
	return bytesMarshalText(s.Bytes()), nil
}

func (s *Signature) UnmarshalText(input []byte) error {
	var b []byte
	if err := bytesUnmarshalText(input, &b); err != nil {
		return err
	}
	sig, err := SignatureFromBytes(b)
	if err != nil {
		return err
	}
	*s = sig
	return nil
}

//
// Number type:
//

// Number represents a hex-encoded number. This type is used for marshaling
// and unmarshalling JSON numbers. When possible, use big.Int or regular integers
// instead.
type Number struct{ x big.Int }

// NumberFromHex converts a hex string to a Number type.
func NumberFromHex(h string) (Number, error) {
	u, err := hexutil.HexToBigInt(h)
	if err != nil {
		return Number{}, err
	}
	return Number{x: *u}, nil
}

// NumberFromHexPtr converts a hex string to a *Number type.
func NumberFromHexPtr(h string) *Number {
	n, err := NumberFromHex(h)
	if err != nil {
		return nil
	}
	return &n
}

// MustNumberFromHex converts a hex string to a Number type. It panics if the
// conversion fails.
func MustNumberFromHex(h string) Number {
	n, err := NumberFromHex(h)
	if err != nil {
		panic(err)
	}
	return n
}

// MustNumberFromHexPtr converts a hex string to a *Number type. It panics if the
// conversion fails.
func MustNumberFromHexPtr(h string) *Number {
	n, err := NumberFromHex(h)
	if err != nil {
		panic(err)
	}
	return &n
}

// NumberFromBytes converts a byte slice to a Number type.
func NumberFromBytes(b []byte) Number {
	return Number{x: *new(big.Int).SetBytes(b)}
}

// NumberFromBytesPtr converts a byte slice to a *Number type.
func NumberFromBytesPtr(b []byte) *Number {
	n := NumberFromBytes(b)
	return &n
}

// NumberFromUint64 converts an uint64 to a Number type.
func NumberFromUint64(x uint64) Number {
	return Number{x: *new(big.Int).SetUint64(x)}
}

// NumberFromUint64Ptr converts an uint64 to a *Number type.
func NumberFromUint64Ptr(x uint64) *Number {
	n := NumberFromUint64(x)
	return &n
}

// NumberFromBigInt converts a big.Int to a Number type.
func NumberFromBigInt(x *big.Int) Number {
	if x == nil {
		return Number{}
	}
	return Number{x: *x}
}

// NumberFromBigIntPtr converts a big.Int to a *Number type.
func NumberFromBigIntPtr(x *big.Int) *Number {
	n := NumberFromBigInt(x)
	return &n
}

// Big returns the big.Int representation of the number.
func (t *Number) Big() *big.Int {
	return new(big.Int).Set(&t.x)
}

// Bytes returns the byte representation of the number.
func (t *Number) Bytes() []byte {
	return t.x.Bytes()
}

// String returns the hex representation of the number.
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

// BytesFromHex converts a hex string to a Bytes type.
func BytesFromHex(h string) (Bytes, error) {
	return hexutil.HexToBytes(h)
}

// BytesFromHexPtr converts a hex string to a *Bytes type.
// If the input is not a valid hex string, it returns nil.
func BytesFromHexPtr(h string) *Bytes {
	b, err := BytesFromHex(h)
	if err != nil {
		return nil
	}
	return &b
}

// MustBytesFromHex converts a hex string to a Bytes type. It panics if the
// input is not a valid hex string.
func MustBytesFromHex(h string) Bytes {
	b, err := BytesFromHex(h)
	if err != nil {
		panic(err)
	}
	return b
}

// MustBytesFromHexPtr converts a hex string to a *Bytes type. It panics if the
// input is not a valid hex string.
func MustBytesFromHexPtr(h string) *Bytes {
	b := MustBytesFromHex(h)
	return &b
}

// BytesFromString converts a string to a Bytes type.
func BytesFromString(s string) Bytes {
	return Bytes(s)
}

// BytesFromStringPtr converts a string to a *Bytes type.
func BytesFromStringPtr(s string) *Bytes {
	b := BytesFromString(s)
	return &b
}

// PadLeft returns a new byte slice padded with zeros to the given length.
// If the byte slice is longer than the given length, it is truncated leaving
// the rightmost bytes.
func (b Bytes) PadLeft(n int) Bytes {
	cp := make([]byte, n)
	if len(b) > n {
		copy(cp, b[len(b)-n:])
	} else {
		copy(cp[n-len(b):], b)
	}
	return cp
}

// PadRight returns a new byte slice padded with zeros to the given length.
// If the byte slice is longer than the given length, it is truncated leaving
// the leftmost bytes.
func (b Bytes) PadRight(n int) Bytes {
	cp := make([]byte, n)
	copy(cp, b)
	return cp
}

// Bytes returns the byte slice.
func (b *Bytes) Bytes() []byte {
	return *b
}

// String returns the hex-encoded string representation of the byte slice.
func (b *Bytes) String() string {
	if b == nil {
		return ""
	}
	return hexutil.BytesToHex(*b)
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(b), nil
}

func (b *Bytes) UnmarshalJSON(input []byte) error {
	return bytesUnmarshalJSON(input, (*[]byte)(b))
}

func (b Bytes) MarshalText() ([]byte, error) {
	return bytesMarshalText(b), nil
}

func (b *Bytes) UnmarshalText(input []byte) error {
	return bytesUnmarshalText(input, (*[]byte)(b))
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
