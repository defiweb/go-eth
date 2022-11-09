package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/hexutil"
)

// bytesMarshalJSON encodes the given bytes as a JSON string where each byte is
// represented by a two-digit hex number. The hex string is always even-length
// and prefixed with "0x".
func bytesMarshalJSON(input []byte) []byte {
	return naiveQuote(bytesMarshalText(input))
}

// bytesMarshalText encodes the given bytes as a string where each byte is
// represented by a two-digit hex number. The hex string is always even-length
// and prefixed with "0x".
func bytesMarshalText(input []byte) []byte {
	return []byte(hexutil.BytesToHex(input))
}

// bytesUnmarshalJSON decodes the given JSON string where each byte is
// represented by a two-digit hex number. The hex string may be prefixed with
// "0x". If the hex string is odd-length, it is padded with a leading zero.
func bytesUnmarshalJSON(input []byte, output *[]byte) error {
	if bytes.Equal(input, []byte("null")) {
		return nil
	}
	return bytesUnmarshalText(naiveUnquote(input), output)
}

// bytesUnmarshalText decodes the given string where each byte is represented by
// a two-digit hex number. The hex string may be prefixed with "0x". If the hex
// string is odd-length, it is padded with a leading zero.
func bytesUnmarshalText(input []byte, output *[]byte) error {
	var err error
	*output, err = hexutil.HexToBytes(string(input))
	return err
}

// fixedBytesUnmarshalJSON works like bytesUnmarshalJSON, but it is designed to
// be used with fixed-size byte arrays. The given byte array must be large
// enough to hold the decoded data.
func fixedBytesUnmarshalJSON(input, output []byte) error {
	if bytes.Equal(input, []byte("null")) {
		return nil
	}
	return fixedBytesUnmarshalText(naiveUnquote(input), output)
}

// fixedBytesUnmarshalText works like bytesUnmarshalText, but it is designed to
// be used with fixed-size byte arrays. The given byte array must be large
// enough to hold the decoded data.
func fixedBytesUnmarshalText(input, output []byte) error {
	data, err := hexutil.HexToBytes(string(input))
	if err != nil {
		return err
	}
	if len(data) > len(output) {
		return fmt.Errorf("hex string has length %d, want %d", len(data), len(output))
	}
	copy(output[len(output)-len(data):], data)
	return nil
}

// numberMarshalJSON encodes the given big integer as JSON string where number
// is resented in hexadecimal format. The hex string is prefixed with "0x".
// Negative numbers are prefixed with "-0x".
func numberMarshalJSON(input *big.Int) []byte {
	return naiveQuote(numberMarshalText(input))
}

// numberMarshalText encodes the given big integer as string where number is
// resented in hexadecimal format. The hex string is prefixed with "0x".
// Negative numbers are prefixed with "-0x".
func numberMarshalText(input *big.Int) []byte {
	return []byte(hexutil.BigIntToHex(input))
}

// numberUnmarshalJSON decodes the given JSON string where number is resented in
// hexadecimal format. The hex string may be prefixed with "0x". Negative numbers
// must start with minus sign.
func numberUnmarshalJSON(input []byte, output *big.Int) error {
	return numberUnmarshalText(naiveUnquote(input), output)
}

// numberUnmarshalText decodes the given string where number is resented in
// hexadecimal format. The hex string may be prefixed with "0x". Negative numbers
// must start with minus sign.
func numberUnmarshalText(input []byte, output *big.Int) error {
	data, err := hexutil.HexToBigInt(string(input))
	if err != nil {
		return err
	}
	output.Set(data)
	return nil
}

// naiveQuote returns a double-quoted string. It does not perform any escaping.
func naiveQuote(i []byte) []byte {
	b := make([]byte, len(i)+2)
	b[0] = '"'
	b[len(b)-1] = '"'
	copy(b[1:], i)
	return b
}

// naiveUnquote returns the string inside the quotes. It does not perform any
// unescaping.
func naiveUnquote(i []byte) []byte {
	if len(i) >= 2 && i[0] == '"' && i[len(i)-1] == '"' {
		return i[1 : len(i)-1]
	}
	return i
}
