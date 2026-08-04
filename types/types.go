package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/hexutil"
)

var (
	// ForceAddressChecksum is a global flag that forces the use of checksummed
	// addresses in text representations.
	//
	// Note: This is a global flag, so it affects all packages that use the
	// Address type.
	ForceAddressChecksum = false
)

// Pad is a padding type.
type Pad uint8

const (
	PadNone  Pad = 0 // PadNone does not allow padding.
	PadLeft  Pad = 1 // PadLeft pads the input with zeros on the left.
	PadRight Pad = 2 // PadRight pads the input with zeros on the right.
)

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
	if err != nil {
		return ZeroAddress, err
	}
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

// AddressFromChecksum parses a checksummed address and returns an Address type.
func AddressFromChecksum(h string) (Address, error) {
	a, err := AddressFromHex(h)
	if err != nil {
		return ZeroAddress, err
	}
	if a.Checksum() != h {
		return ZeroAddress, fmt.Errorf("invalid checksum: expected %s, got %s", a.Checksum(), h)
	}
	return a, nil
}

// MustAddressFromChecksum parses a checksummed address and returns an Address type.
// It panics if the address is invalid.
func MustAddressFromChecksum(h string) Address {
	a, err := AddressFromChecksum(h)
	if err != nil {
		panic(err)
	}
	return a
}

// AddressFromChecksumPtr parses a checksummed address and returns an Address type.
// It returns nil if the address is invalid.
func AddressFromChecksumPtr(h string) *Address {
	a, err := AddressFromChecksum(h)
	if err != nil {
		return nil
	}
	return &a
}

// Bytes returns the byte representation of the address.
func (t Address) Bytes() []byte {
	return t[:]
}

// String returns the hex representation of the address.
func (t Address) String() string {
	if ForceAddressChecksum {
		return t.Checksum()
	}
	return hexutil.BytesToHex(t[:])
}

// Checksum returns the address with the checksum calculated according to
// EIP-55.
func (t Address) Checksum() string {
	hex := []byte(hexutil.BytesToHex(t[:])[2:])
	hash := crypto.Keccak256(hex)
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

// MarshalJSON implements the json.Marshaler interface.
func (t Address) MarshalJSON() ([]byte, error) {
	if ForceAddressChecksum {
		return naiveQuote([]byte(t.Checksum())), nil
	}
	return bytesMarshalJSON(t[:]), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Address) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t Address) MarshalText() ([]byte, error) {
	if ForceAddressChecksum {
		return []byte(t.Checksum()), nil
	}
	return bytesMarshalText(t[:]), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Address) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

// EncodeRLP implements the rlp.Encoder interface.
func (t Address) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

// DecodeRLP implements the rlp.Decoder interface.
func (t *Address) DecodeRLP(data []byte) (int, error) {
	r, n, err := rlp.DecodeLazy(data)
	if err != nil {
		return 0, err
	}
	b, err := r.Bytes()
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		*t = ZeroAddress
		return n, nil
	}
	if len(b) != AddressLength {
		return 0, fmt.Errorf("invalid address length %d", len(b))
	}
	copy(t[:], b)
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

// HashKeccak256 calculates the Keccak256 hash of the given data.
func HashKeccak256(data ...[]byte) Hash {
	return Hash(crypto.Keccak256(data...))
}

// HashFromHex parses a hash in hex format and returns a Hash type.
// If hash is longer than 32 bytes, it returns an error.
func HashFromHex(h string, pad Pad) (Hash, error) {
	b, err := hexutil.HexToBytes(h)
	if err != nil {
		return ZeroHash, err
	}
	return HashFromBytes(b, pad)
}

// HashFromHexPtr parses a hash in hex format and returns a *Hash type.
// If hash is longer than 32 bytes, it returns an error.
// It returns nil if the hash is invalid.
func HashFromHexPtr(h string, pad Pad) *Hash {
	hash, err := HashFromHex(h, pad)
	if err != nil {
		return nil
	}
	return &hash
}

// MustHashFromHex parses a hash in hex format and returns a Hash type.
// If hash is longer than 32 bytes, it returns an error.
// It panics if the hash is invalid.
func MustHashFromHex(h string, pad Pad) Hash {
	hash, err := HashFromHex(h, pad)
	if err != nil {
		panic(err)
	}
	return hash
}

// MustHashFromHexPtr parses a hash in hex format and returns a *Hash type.
// If hash is longer than 32 bytes, it returns an error.
// It panics if the hash is invalid.
func MustHashFromHexPtr(h string, pad Pad) *Hash {
	hash := MustHashFromHex(h, pad)
	return &hash
}

// HashFromBytes converts a byte slice to a Hash type.
// If bytes is longer than 32 bytes, it returns an error.
func HashFromBytes(b []byte, pad Pad) (Hash, error) {
	var h Hash
	if len(b) > HashLength {
		return ZeroHash, fmt.Errorf("hash too long %d", len(b))
	}
	switch pad {
	case PadLeft:
		copy(h[HashLength-len(b):], b)
	case PadRight:
		copy(h[:], b)
	case PadNone:
		if len(b) != HashLength {
			return ZeroHash, fmt.Errorf("invalid hash length %d", len(b))
		}
		copy(h[:], b)
	}
	return h, nil
}

// HashFromBytesPtr converts a byte slice to a *Hash type.
// If bytes is longer than 32 bytes, it returns an error.
// It returns nil if the hash is invalid.
func HashFromBytesPtr(b []byte, pad Pad) *Hash {
	h, err := HashFromBytes(b, pad)
	if err != nil {
		return nil
	}
	return &h
}

// MustHashFromBytes converts a byte slice to a Hash type.
// If bytes is longer than 32 bytes, it returns an error.
// It panics if the hash is invalid.
func MustHashFromBytes(b []byte, pad Pad) Hash {
	h, err := HashFromBytes(b, pad)
	if err != nil {
		panic(err)
	}
	return h
}

// MustHashFromBytesPtr converts a byte slice to a *Hash type.
// If bytes is longer than 32 bytes, it returns an error.
// It panics if the hash is invalid.
func MustHashFromBytesPtr(b []byte, pad Pad) *Hash {
	h := MustHashFromBytes(b, pad)
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
	return HashFromBytes(b, PadLeft)
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

// MarshalJSON implements the json.Marshaler interface.
func (t Hash) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Hash) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t Hash) MarshalText() ([]byte, error) {
	return bytesMarshalText(t[:]), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Hash) UnmarshalText(input []byte) error {
	return fixedBytesUnmarshalText(input, t[:])
}

// EncodeRLP implements the rlp.Encoder interface.
func (t Hash) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

// DecodeRLP implements the rlp.Decoder interface.
func (t *Hash) DecodeRLP(data []byte) (int, error) {
	r, n, err := rlp.DecodeLazy(data)
	if err != nil {
		return 0, err
	}
	b, err := r.Bytes()
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		*t = ZeroHash
		return n, nil
	}
	if len(b) != HashLength {
		return 0, fmt.Errorf("invalid hash length %d", len(b))
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
	earliestBlockNumber  = -1
	latestBlockNumber    = -2
	pendingBlockNumber   = -3
	safeBlockNumber      = -4
	finalizedBlockNumber = -5
)

var (
	EarliestBlockNumber  = BlockNumber{x: *new(big.Int).SetInt64(earliestBlockNumber)}
	LatestBlockNumber    = BlockNumber{x: *new(big.Int).SetInt64(latestBlockNumber)}
	PendingBlockNumber   = BlockNumber{x: *new(big.Int).SetInt64(pendingBlockNumber)}
	SafeBlockNumber      = BlockNumber{x: *new(big.Int).SetInt64(safeBlockNumber)}
	FinalizedBlockNumber = BlockNumber{x: *new(big.Int).SetInt64(finalizedBlockNumber)}
)

// BlockNumberFromHex converts a string to a BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "safe", "finalized", "pending".
// If the string is not a valid block number, it returns an error.
func BlockNumberFromHex(h string) (BlockNumber, error) {
	b := &BlockNumber{}
	err := b.UnmarshalText([]byte(h))
	return *b, err
}

// BlockNumberFromHexPtr converts a string to a *BlockNumber type.
// The string can be a hex number or one of the following strings:
// "earliest", "latest", "safe", "finalized", "pending".
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
// "earliest", "latest", "safe", "finalized", "pending".
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
// "earliest", "latest", "safe", "finalized", "pending".
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

// IsSafe returns true if the block tag is "safe".
func (t *BlockNumber) IsSafe() bool {
	return t.Big().Int64() == safeBlockNumber
}

// IsFinalized returns true if the block tag is "finalized".
func (t *BlockNumber) IsFinalized() bool {
	return t.Big().Int64() == finalizedBlockNumber
}

// IsTag returns true if the block tag is used.
func (t *BlockNumber) IsTag() bool {
	return t.Big().Sign() < 0
}

// Big returns the big.Int representation of the block number.
// It returns a negative number if the block tag is used:
//   - earliest: -1
//   - latest: -2
//   - pending: -3
//   - safe: -4
//   - finalized: -5
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
	case t.IsSafe():
		return "safe"
	case t.IsFinalized():
		return "finalized"
	default:
		return "0x" + t.x.Text(16)
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (t BlockNumber) MarshalJSON() ([]byte, error) {
	b, err := t.MarshalText()
	if err != nil {
		return nil, err
	}
	return naiveQuote(b), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *BlockNumber) UnmarshalJSON(input []byte) error {
	input, ok := naiveUnquote(input)
	if !ok {
		return fmt.Errorf("invalid JSON string: %s", input)
	}
	return t.UnmarshalText(input)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t BlockNumber) MarshalText() ([]byte, error) {
	switch {
	case t.IsEarliest():
		return []byte("earliest"), nil
	case t.IsLatest():
		return []byte("latest"), nil
	case t.IsPending():
		return []byte("pending"), nil
	case t.IsSafe():
		return []byte("safe"), nil
	case t.IsFinalized():
		return []byte("finalized"), nil
	default:
		return []byte(hexutil.BigIntToHex(&t.x)), nil
	}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
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
	case "safe":
		*t = BlockNumber{x: *new(big.Int).SetInt64(safeBlockNumber)}
		return nil
	case "finalized":
		*t = BlockNumber{x: *new(big.Int).SetInt64(finalizedBlockNumber)}
		return nil
	default:
		u, err := hexutil.HexToBigInt(string(input))
		if err != nil {
			return err
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

var _ = crypto.Signature(Signature{})

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

// Equal returns true if the signature is equal to the given signature.
//
// Nil values are considered as zero.
func (s Signature) Equal(c Signature) bool {
	sv, sr, ss := s.V, s.R, s.S
	cv, cr, cs := c.V, c.R, c.S
	if sv == nil {
		sv = new(big.Int)
	}
	if sr == nil {
		sr = new(big.Int)
	}
	if ss == nil {
		ss = new(big.Int)
	}
	if cv == nil {
		cv = new(big.Int)
	}
	if cr == nil {
		cr = new(big.Int)
	}
	if cs == nil {
		cs = new(big.Int)
	}
	return sv.Cmp(cv) == 0 && sr.Cmp(cr) == 0 && ss.Cmp(cs) == 0
}

// Copy returns a deep copy of the signature.
func (s Signature) Copy() *Signature {
	cpy := &Signature{}
	if s.V != nil {
		cpy.V = new(big.Int).Set(s.V)
	}
	if s.R != nil {
		cpy.R = new(big.Int).Set(s.R)
	}
	if s.S != nil {
		cpy.S = new(big.Int).Set(s.S)
	}
	return cpy
}

// MarshalJSON implements the json.Marshaler interface.
func (s Signature) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(s.Bytes()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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

// MarshalText implements the encoding.TextMarshaler interface.
func (s Signature) MarshalText() ([]byte, error) {
	return bytesMarshalText(s.Bytes()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
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
	return Number{x: *copyBigInt(x)}
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

// Bytes returns the absolute value of a number as a big-endian byte slice.
func (t *Number) Bytes() []byte {
	return t.x.Bytes()
}

// String returns the hex representation of the number.
func (t *Number) String() string {
	return hexutil.BigIntToHex(&t.x)
}

// MarshalJSON implements the json.Marshaler interface.
func (t Number) MarshalJSON() ([]byte, error) {
	return numberMarshalJSON(t.Big()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Number) UnmarshalJSON(input []byte) error {
	return numberUnmarshalJSON(input, &t.x)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (t Number) MarshalText() ([]byte, error) {
	return numberMarshalText(t.Big()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
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

// MarshalJSON implements the json.Marshaler interface.
func (b Bytes) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(b), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	return bytesUnmarshalJSON(input, (*[]byte)(b))
}

// MarshalText implements the encoding.TextMarshaler interface.
func (b Bytes) MarshalText() ([]byte, error) {
	return bytesMarshalText(b), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *Bytes) UnmarshalText(input []byte) error {
	return bytesUnmarshalText(input, (*[]byte)(b))
}

//
// Internal types:
//

const (
	bloomLength = 256
	nonceLength = 8
)

// oneOrList is a type that can marshal and unmarshal a single element or a list
// of elements.
type oneOrList[T any] []T

func (l oneOrList[T]) MarshalJSON() ([]byte, error) {
	if len(l) == 1 {
		return json.Marshal(l[0])
	}
	return json.Marshal([]T(l))
}

func (l *oneOrList[T]) UnmarshalJSON(input []byte) error {
	if len(input) >= 1 && input[0] == '[' {
		var list []T
		if err := json.Unmarshal(input, &list); err != nil {
			return err
		}
		*l = list
		return nil
	}
	var i T
	if err := json.Unmarshal(input, &i); err != nil {
		return err
	}
	*l = oneOrList[T]{i}
	return nil
}

// kzgHash is a fixed-length byte array used for KZG hash.
type kzgHash [crypto.KZGHashSize]byte

func (t kzgHash) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *kzgHash) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t kzgHash) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

func (t *kzgHash) DecodeRLP(data []byte) (int, error) {
	return fixedBytesDecodeRLP(data, t[:])
}

// kzgBlob is a fixed-length byte array used for KZG blob.
type kzgBlob [crypto.KZGBlobSize]byte

func (t kzgBlob) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *kzgBlob) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t kzgBlob) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

func (t *kzgBlob) DecodeRLP(data []byte) (int, error) {
	return fixedBytesDecodeRLP(data, t[:])
}

// kzgCommitment is a fixed-length byte array used for KZG commitment.
type kzgCommitment [crypto.KZGCommitmentSize]byte

func (t kzgCommitment) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *kzgCommitment) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t kzgCommitment) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

func (t *kzgCommitment) DecodeRLP(data []byte) (int, error) {
	return fixedBytesDecodeRLP(data, t[:])
}

// kzgProof is a fixed-length byte array used for KZG proof.
type kzgProof [crypto.KZGProofSize]byte

func (t kzgProof) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *kzgProof) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func (t kzgProof) EncodeRLP() ([]byte, error) {
	return rlp.Encode(rlp.Bytes(t[:]))
}

func (t *kzgProof) DecodeRLP(data []byte) (int, error) {
	return fixedBytesDecodeRLP(data, t[:])
}

// bloom is a fixed-length byte array used for bloom filter.
type bloom [bloomLength]byte

func bloomFromBytes(x []byte) bloom {
	var b [bloomLength]byte
	if len(x) > len(b) {
		return b
	}
	copy(b[bloomLength-len(x):], x)
	return b
}

func (t *bloom) Bytes() []byte {
	return t[:]
}

func (t bloom) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *bloom) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}

func nonceFromBigInt(x *big.Int) (n nonce) {
	if x == nil {
		return nonce{}
	}
	b := x.Bytes()
	copy(n[nonceLength-len(b):], b)
	return n
}

// nonce is a fixed-length byte array used for nonce.
type nonce [nonceLength]byte

func (t *nonce) Big() *big.Int {
	return new(big.Int).SetBytes(t[:])
}

func (t nonce) MarshalJSON() ([]byte, error) {
	return bytesMarshalJSON(t[:]), nil
}

func (t *nonce) UnmarshalJSON(input []byte) error {
	return fixedBytesUnmarshalJSON(input, t[:])
}
