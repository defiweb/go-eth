package abi

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

func Test_mappingRules(t *testing.T) {
	// Test if mapping rules specified in the README work as expected.

	tests := []struct {
		name       string
		goTyp      any
		solTyp     string
		src        any
		wantDst    any
		wantEncErr bool
		wantDecErr bool
	}{
		// intX <=> intX
		{
			name:    "intX<=>intX",
			goTyp:   new(int),
			solTyp:  "int",
			src:     1,
			wantDst: 1,
		},
		{
			name:    "intX<=>intX#to-smaller-type",
			goTyp:   new(int8),
			solTyp:  "int256",
			src:     1,
			wantDst: int8(1),
		},
		{
			name:       "intX<=>intX#encode-error",
			goTyp:      new(int),
			solTyp:     "int8",
			src:        256,
			wantEncErr: true,
		},
		{
			name:       "intX<=>intX#decode-error",
			goTyp:      new(int8),
			solTyp:     "int",
			src:        256,
			wantDecErr: true,
		},

		// intX <=> uintX
		{
			name:    "intX<=>uintX",
			goTyp:   new(int),
			solTyp:  "uint",
			src:     1,
			wantDst: 1,
		},
		{
			name:    "intX<=>uintX#to-smaller-type",
			goTyp:   new(int8),
			solTyp:  "uint256",
			src:     1,
			wantDst: int8(1),
		},
		{
			name:       "intX<=>uintX#encode-error",
			goTyp:      new(int),
			solTyp:     "uint8",
			src:        256,
			wantEncErr: true,
		},
		{
			name:       "intX<=>uintX#decode-error",
			goTyp:      new(int8),
			solTyp:     "uint",
			src:        256,
			wantDecErr: true,
		},

		// intX <=> bool
		{
			name:       "intX<=>bool#encode-error",
			goTyp:      new(int),
			solTyp:     "bool",
			src:        1,
			wantEncErr: true,
		},
		{
			name:       "intX<=>bool#decode-error",
			goTyp:      new(int),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// intX <=> string
		{
			name:       "intX<=>string#encode-error",
			goTyp:      new(int),
			solTyp:     "string",
			src:        1,
			wantEncErr: true,
		},
		{
			name:       "intX<=>string#decode-error",
			goTyp:      new(int),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// intX <=> bytes
		{
			name:       "intX<=>bytes#encode-error",
			goTyp:      new(int),
			solTyp:     "bytes",
			src:        1,
			wantEncErr: true,
		},
		{
			name:       "intX<=>bytes#decode-error",
			goTyp:      new(int),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// intX <=> bytesX
		{
			name:    "intX<=>bytesX",
			goTyp:   new(int),
			solTyp:  "bytes32",
			src:     1,
			wantDst: 1,
		},
		{
			name:    "intX<=>bytesX#negative",
			goTyp:   new(int),
			solTyp:  "bytes32",
			src:     -1,
			wantDst: -1,
		},
		{
			name:       "intX<=>bytesX#encode-error",
			goTyp:      new(int64),
			solTyp:     "bytes8",
			src:        1,
			wantEncErr: true,
		},
		{
			name:       "intX<=>bytesX#decode-error",
			goTyp:      new(int64),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// intX <=> address
		{
			name:       "intX<=>address#encode-error",
			goTyp:      new(int),
			solTyp:     "address",
			src:        1,
			wantEncErr: true,
		},
		{
			name:       "intX<=>address#decode-error",
			goTyp:      new(int),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// uintX <=> intX
		{
			name:    "uintX<=>intX",
			goTyp:   new(uint),
			solTyp:  "int",
			src:     uint(1),
			wantDst: uint(1),
		},
		{
			name:    "uintX<=>intX#to-smaller-type",
			goTyp:   new(uint8),
			solTyp:  "int256",
			src:     uint(1),
			wantDst: uint8(1),
		},
		{
			name:       "uintX<=>intX#encode-error",
			goTyp:      new(uint),
			solTyp:     "int8",
			src:        uint(256),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>intX#decode-error",
			goTyp:      new(uint8),
			solTyp:     "int",
			src:        uint(256),
			wantDecErr: true,
		},

		// uintX <=> uintX
		{
			name:    "uintX<=>uintX",
			goTyp:   new(uint),
			solTyp:  "uint",
			src:     uint(1),
			wantDst: uint(1),
		},
		{
			name:    "uintX<=>uintX#to-smaller-type",
			goTyp:   new(uint8),
			solTyp:  "uint256",
			src:     uint(1),
			wantDst: uint8(1),
		},
		{
			name:       "uintX<=>uintX#encode-error",
			goTyp:      new(uint),
			solTyp:     "uint8",
			src:        uint(256),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>uintX#decode-error",
			goTyp:      new(uint8),
			solTyp:     "uint",
			src:        uint(256),
			wantDecErr: true,
		},

		// uintX <=> bool
		{
			name:       "uintX<=>bool#encode-error",
			goTyp:      new(uint),
			solTyp:     "bool",
			src:        uint(1),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>bool#decode-error",
			goTyp:      new(uint),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// uintX <=> string
		{
			name:       "uintX<=>string#encode-error",
			goTyp:      new(uint),
			solTyp:     "string",
			src:        uint(1),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>string#decode-error",
			goTyp:      new(uint),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// uintX <=> bytes
		{
			name:       "uintX<=>bytes#encode-error",
			goTyp:      new(uint),
			solTyp:     "bytes",
			src:        uint(1),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>bytes#decode-error",
			goTyp:      new(uint),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// uintX <=> bytesX
		{
			name:    "uintX<=>bytesX",
			goTyp:   new(uint),
			solTyp:  "bytes32",
			src:     uint(1),
			wantDst: uint(1),
		},
		{
			name:       "uintX<=>bytesX#encode-error",
			goTyp:      new(uint64),
			solTyp:     "bytes8",
			src:        uint64(1),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>bytesX#decode-error",
			goTyp:      new(uint64),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// uintX <=> address
		{
			name:       "uintX<=>address#encode-error",
			goTyp:      new(uint),
			solTyp:     "address",
			src:        uint(1),
			wantEncErr: true,
		},
		{
			name:       "uintX<=>address#decode-error",
			goTyp:      new(uint),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// bool <=> intX
		{
			name:       "bool<=>intX#encode-error",
			goTyp:      new(bool),
			solTyp:     "int",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>intX#decode-error",
			goTyp:      new(bool),
			solTyp:     "int",
			src:        1,
			wantDecErr: true,
		},

		// bool <=> uintX
		{
			name:       "bool<=>uintX#encode-error",
			goTyp:      new(bool),
			solTyp:     "uint",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>uintX#decode-error",
			goTyp:      new(bool),
			solTyp:     "uint",
			src:        uint(1),
			wantDecErr: true,
		},

		// bool <=> bool
		{
			name:    "bool<=>bool",
			goTyp:   new(bool),
			solTyp:  "bool",
			src:     true,
			wantDst: true,
		},

		// bool <=> string
		{
			name:       "bool<=>string#encode-error",
			goTyp:      new(bool),
			solTyp:     "string",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>string#decode-error",
			goTyp:      new(bool),
			solTyp:     "string",
			src:        "true",
			wantDecErr: true,
		},

		// bool <=> bytes
		{
			name:       "bool<=>bytes#encode-error",
			goTyp:      new(bool),
			solTyp:     "bytes",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>bytes#decode-error",
			goTyp:      new(bool),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// bool <=> bytesX
		{
			name:       "bool<=>bytesX#encode-error",
			goTyp:      new(bool),
			solTyp:     "bytes8",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>bytesX#decode-error",
			goTyp:      new(bool),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// bool <=> address
		{
			name:       "bool<=>address#encode-error",
			goTyp:      new(bool),
			solTyp:     "address",
			src:        true,
			wantEncErr: true,
		},
		{
			name:       "bool<=>address#decode-error",
			goTyp:      new(bool),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// string <=> intX
		{
			name:    "string<=>intX",
			goTyp:   new(string),
			solTyp:  "int",
			src:     "0x1",
			wantDst: "0x1",
		},
		{
			name:    "string<=>intX#negative",
			goTyp:   new(string),
			solTyp:  "int",
			src:     "-0x1",
			wantDst: "-0x1",
		},
		{
			name:       "string<=>intX#too-small-dest",
			goTyp:      new(int8),
			solTyp:     "int8",
			src:        "0xff",
			wantEncErr: true,
		},
		{
			name:    "string<=>intX#skipped-0x-prefix",
			goTyp:   new(string),
			solTyp:  "int",
			src:     "100",
			wantDst: "0x100",
		},
		{
			name:       "string<=>intX#invalid-format",
			goTyp:      new(string),
			solTyp:     "int",
			src:        "foo",
			wantEncErr: true,
		},

		// string <=> uintX
		{
			name:    "string<=>uintX",
			goTyp:   new(string),
			solTyp:  "uint",
			src:     "0x1",
			wantDst: "0x1",
		},
		{
			name:       "string<=>uintX#negative",
			goTyp:      new(string),
			solTyp:     "uint",
			src:        "-0x1",
			wantEncErr: true,
		},
		{
			name:       "string<=>uintX#too-small-dest",
			goTyp:      new(int8),
			solTyp:     "uint8",
			src:        "0x100",
			wantEncErr: true,
		},
		{
			name:    "string<=>uintX#skipped-0x-prefix",
			goTyp:   new(string),
			solTyp:  "uint",
			src:     "100",
			wantDst: "0x100",
		},
		{
			name:       "string<=>uintX#invalid-format",
			goTyp:      new(string),
			solTyp:     "uint",
			src:        "foo",
			wantEncErr: true,
		},

		// string <=> bool
		{
			name:       "string<=>bool#encode-error",
			goTyp:      new(string),
			solTyp:     "bool",
			src:        "true",
			wantEncErr: true,
		},
		{
			name:       "string<=>bool#decode-error",
			goTyp:      new(string),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// string <=> string
		{
			name:    "string<=>string",
			goTyp:   new(string),
			solTyp:  "string",
			src:     "foo",
			wantDst: "foo",
		},

		// string <=> bytes
		{
			name:    "string<=>bytes",
			goTyp:   new(string),
			solTyp:  "bytes",
			src:     "0x666f6f",
			wantDst: "0x666f6f",
		},
		{
			name:    "string<=>bytes#skipped-0x-prefix",
			goTyp:   new(string),
			solTyp:  "bytes",
			src:     "666f6f",
			wantDst: "0x666f6f",
		},
		{
			name:       "string<=>bytes#invalid-format",
			goTyp:      new(string),
			solTyp:     "bytes",
			src:        "foo",
			wantEncErr: true,
		},

		// string <=> bytesX
		{
			name:    "string<=>bytesX",
			goTyp:   new(string),
			solTyp:  "bytes32",
			src:     "0x666f6f0000000000000000000000000000000000000000000000000000000000",
			wantDst: "0x666f6f0000000000000000000000000000000000000000000000000000000000",
		},

		// string <=> address
		{
			name:    "string<=>address",
			goTyp:   new(string),
			solTyp:  "address",
			src:     "0x1234567890123456789012345678901234567890",
			wantDst: "0x1234567890123456789012345678901234567890",
		},
		{
			name:    "string<=>address#skipped-0x-prefix",
			goTyp:   new(string),
			solTyp:  "address",
			src:     "1234567890123456789012345678901234567890",
			wantDst: "0x1234567890123456789012345678901234567890",
		},
		{
			name:       "string<=>address#too-short",
			goTyp:      new(string),
			solTyp:     "address",
			src:        "0x123456789012345678901234567890123456789",
			wantEncErr: true,
		},
		{
			name:       "string<=>address#too-long",
			goTyp:      new(string),
			solTyp:     "address",
			src:        "0x12345678901234567890123456789012345678900",
			wantEncErr: true,
		},

		// []byte <=> intX
		{
			name:       "[]byte<=>intX#encode-error",
			goTyp:      new([]byte),
			solTyp:     "int",
			src:        []byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[]byte<=>intX#decode-error",
			goTyp:      new([]byte),
			solTyp:     "int",
			src:        "0x1",
			wantDecErr: true,
		},

		// []byte <=> uintX
		{
			name:       "[]byte<=>uintX#encode-error",
			goTyp:      new([]byte),
			solTyp:     "uint",
			src:        []byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[]byte<=>uintX#decode-error",
			goTyp:      new([]byte),
			solTyp:     "uint",
			src:        "0x1",
			wantDecErr: true,
		},

		// []byte <=> bool
		{
			name:       "[]byte<=>bool#encode-error",
			goTyp:      new([]byte),
			solTyp:     "bool",
			src:        []byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[]byte<=>bool#decode-error",
			goTyp:      new([]byte),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// []byte <=> string
		{
			name:    "[]byte<=>string",
			goTyp:   new([]byte),
			solTyp:  "string",
			src:     []byte{0x1},
			wantDst: []byte{0x1},
		},

		// []byte <=> bytes
		{
			name:    "[]byte<=>bytes",
			goTyp:   new([]byte),
			solTyp:  "bytes",
			src:     []byte{0x66, 0x6f, 0x6f},
			wantDst: []byte{0x66, 0x6f, 0x6f},
		},

		// []byte <=> bytesX
		{
			name:    "[]byte<=>bytesX",
			goTyp:   new([]byte),
			solTyp:  "bytes4",
			src:     []byte{0x66, 0x6f, 0x6f, 0x00},
			wantDst: []byte{0x66, 0x6f, 0x6f, 0x00},
		},
		{
			name:       "[]byte<=>bytesX#too-short",
			goTyp:      new([]byte),
			solTyp:     "bytes2",
			src:        []byte{0x66, 0x6f, 0x6f},
			wantEncErr: true,
		},

		// []byte <=> address
		{
			name:    "[]byte<=>address",
			goTyp:   new([]byte),
			solTyp:  "address",
			src:     []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90},
			wantDst: []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90},
		},
		{
			name:       "[]byte<=>address#too-short",
			goTyp:      new([]byte),
			solTyp:     "address",
			src:        []byte{0x12},
			wantEncErr: true,
		},
		{
			name:       "[]byte<=>address#too-long",
			goTyp:      new([]byte),
			solTyp:     "address",
			src:        []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12},
			wantEncErr: true,
		},

		// [X]byte <=> intX
		{
			name:       "[X]byte<=>uintX#encode-error",
			goTyp:      new([1]byte),
			solTyp:     "uint",
			src:        []byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[X]byte<=>uintX#decode-error",
			goTyp:      new([1]byte),
			solTyp:     "uint",
			src:        "0x1",
			wantDecErr: true,
		},

		// [X]byte <=> bool
		{
			name:       "[X]byte<=>bool#encode-error",
			goTyp:      new([1]byte),
			solTyp:     "bool",
			src:        []byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[X]byte<=>bool#decode-error",
			goTyp:      new([1]byte),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// [X]byte <=> string
		{
			name:       "[X]byte<=>string#encode-error",
			goTyp:      new([1]byte),
			solTyp:     "string",
			src:        [1]byte{0x1},
			wantEncErr: true,
		},
		{
			name:       "[X]byte<=>string#decode-error",
			goTyp:      new([1]byte),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// [X]byte <=> bytes
		{
			name:    "[X]byte<=>bytes",
			goTyp:   new([3]byte),
			solTyp:  "bytes",
			src:     []byte{0x66, 0x6f, 0x6f},
			wantDst: [3]byte{0x66, 0x6f, 0x6f},
		},

		// [X]byte <=> bytesX
		{
			name:    "[X]byte<=>bytesX",
			goTyp:   new([3]byte),
			solTyp:  "bytes3",
			src:     [3]byte{0x66, 0x6f, 0x6f},
			wantDst: [3]byte{0x66, 0x6f, 0x6f},
		},
		{
			name:       "[X]byte<=>bytesX#too-short",
			goTyp:      new([3]byte),
			solTyp:     "bytes2",
			src:        [3]byte{0x66, 0x6f, 0x6f},
			wantEncErr: true,
		},

		// [X]byte <=> address
		{
			name:    "[X]byte<=>address",
			goTyp:   new([20]byte),
			solTyp:  "address",
			src:     [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90},
			wantDst: [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90},
		},
		{
			name:       "[X]byte<=>address#too-short",
			goTyp:      new([19]byte),
			solTyp:     "address",
			src:        [19]byte{0x12},
			wantEncErr: true,
		},
		{
			name:       "[X]byte<=>address#too-long",
			goTyp:      new([21]byte),
			solTyp:     "address",
			src:        [21]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12},
			wantEncErr: true,
		},

		// big.Int <=> intX
		{
			name:    "big.Int<=>int",
			goTyp:   new(big.Int),
			solTyp:  "int",
			src:     big.NewInt(0x1234),
			wantDst: big.NewInt(0x1234),
		},
		{
			name:    "big.Int<=>int#negative",
			goTyp:   new(big.Int),
			solTyp:  "int",
			src:     big.NewInt(-0x1234),
			wantDst: big.NewInt(-0x1234),
		},
		{
			name:    "big.Int<=>int#too-smaller-type",
			goTyp:   new(big.Int),
			solTyp:  "int8",
			src:     big.NewInt(42),
			wantDst: big.NewInt(42),
		},
		{
			name:       "big.Int<=>int#overflow",
			goTyp:      new(big.Int),
			solTyp:     "int8",
			src:        big.NewInt(0x1234),
			wantEncErr: true,
		},

		// big.Int <=> uintX
		{
			name:    "big.Int<=>uint",
			goTyp:   new(big.Int),
			solTyp:  "uint",
			src:     big.NewInt(0x1234),
			wantDst: big.NewInt(0x1234),
		},
		{
			name:    "big.Int<=>uint#too-smaller-type",
			goTyp:   new(big.Int),
			solTyp:  "uint8",
			src:     big.NewInt(42),
			wantDst: big.NewInt(42),
		},
		{
			name:       "big.Int<=>uint#overflow",
			goTyp:      new(big.Int),
			solTyp:     "uint8",
			src:        big.NewInt(0x1234),
			wantEncErr: true,
		},

		// big.Int <=> bool
		{
			name:       "big.Int<=>bool#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "bool",
			src:        big.NewInt(0x1234),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>bool#decode-error",
			goTyp:      new(big.Int),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// big.Int <=> bytes
		{
			name:       "big.Int<=>bytes#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes",
			src:        big.NewInt(0x1234),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>bytes#decode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes",
			src:        []byte{0x12, 0x34},
			wantDecErr: true,
		},

		// big.Int <=> string
		{
			name:       "big.Int<=>string#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "string",
			src:        big.NewInt(1),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>string#decode-error",
			goTyp:      new(int),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// big.Int <=> bytes
		{
			name:       "big.Int<=>bytes#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes",
			src:        big.NewInt(1),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>bytes#decode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// big.Int <=> bytesX
		{
			name:    "big.Int<=>bytesX",
			goTyp:   new(big.Int),
			solTyp:  "bytes32",
			src:     big.NewInt(1),
			wantDst: big.NewInt(1),
		},
		{
			name:    "intX<=>bytesX#negative",
			goTyp:   new(big.Int),
			solTyp:  "bytes32",
			src:     big.NewInt(-1),
			wantDst: big.NewInt(-1),
		},
		{
			name:       "big.Int<=>bytesX#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes8",
			src:        big.NewInt(1),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>bytesX#decode-error",
			goTyp:      new(big.Int),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// big.Int <=> address
		{
			name:       "big.Int<=>address#encode-error",
			goTyp:      new(big.Int),
			solTyp:     "address",
			src:        new(big.Int),
			wantEncErr: true,
		},
		{
			name:       "big.Int<=>address#decode-error",
			goTyp:      new(big.Int),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// types.Address <=> intX
		{
			name:       "types.Address<=>int#encode-error",
			goTyp:      new(types.Address),
			solTyp:     "int",
			src:        new(types.Address),
			wantEncErr: true,
		},
		{
			name:       "types.Address<=>int#decode-error",
			goTyp:      new(types.Address),
			solTyp:     "int",
			src:        1,
			wantDecErr: true,
		},

		// types.Address <=> uintX
		{
			name:       "types.Address<=>uint#encode-error",
			goTyp:      new(types.Address),
			solTyp:     "uint",
			src:        new(types.Address),
			wantEncErr: true,
		},
		{
			name:       "types.Address<=>uint#decode-error",
			goTyp:      new(types.Address),
			solTyp:     "uint",
			src:        1,
			wantDecErr: true,
		},

		// types.Address <=> bool
		{
			name:       "types.Address<=>bool#encode-error",
			goTyp:      new(types.Address),
			solTyp:     "bool",
			src:        new(types.Address),
			wantEncErr: true,
		},
		{
			name:       "types.Address<=>bool#decode-error",
			goTyp:      new(types.Address),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// types.Address <=> string
		{
			name:       "types.Address<=>string#encode-error",
			goTyp:      new(types.Address),
			solTyp:     "string",
			src:        new(types.Address),
			wantEncErr: true,
		},
		{
			name:       "types.Address<=>string#decode-error",
			goTyp:      new(types.Address),
			solTyp:     "string",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// types.Address <=> bytes
		{
			name:    "types.Address<=>bytes",
			goTyp:   new(types.Address),
			solTyp:  "bytes",
			src:     types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantDst: types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
		},

		// types.Address <=> bytesX
		{
			name:    "types.Address<=>bytesX",
			goTyp:   new(types.Address),
			solTyp:  "bytes20",
			src:     types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantDst: types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
		},
		{
			name:       "types.Address<=>bytesX#too-long",
			goTyp:      new(types.Address),
			solTyp:     "bytes32",
			src:        types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantEncErr: true,
		},
		{
			name:       "types.Address<=>bytesX#too-short",
			goTyp:      new(types.Address),
			solTyp:     "bytes8",
			src:        types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantEncErr: true,
		},

		// types.Address <=> address
		{
			name:    "types.Address<=>address",
			goTyp:   new(types.Address),
			solTyp:  "address",
			src:     types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantDst: types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
		},

		// types.Hash <=> intX
		{
			name:       "types.Hash<=>int#encode-error",
			goTyp:      new(types.Hash),
			solTyp:     "int",
			src:        new(types.Hash),
			wantEncErr: true,
		},
		{
			name:       "types.Hash<=>int#decode-error",
			goTyp:      new(types.Hash),
			solTyp:     "int",
			src:        1,
			wantDecErr: true,
		},

		// types.Hash <=> uintX
		{
			name:       "types.Hash<=>uint#encode-error",
			goTyp:      new(types.Hash),
			solTyp:     "uint",
			src:        new(types.Hash),
			wantEncErr: true,
		},
		{
			name:       "types.Hash<=>uint#decode-error",
			goTyp:      new(types.Hash),
			solTyp:     "uint",
			src:        1,
			wantDecErr: true,
		},

		// types.Hash <=> bool
		{
			name:       "types.Hash<=>bool#encode-error",
			goTyp:      new(types.Hash),
			solTyp:     "bool",
			src:        new(types.Hash),
			wantEncErr: true,
		},
		{
			name:       "types.Hash<=>bool#decode-error",
			goTyp:      new(types.Hash),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// types.Hash <=> string
		{
			name:       "types.Hash<=>string#encode-error",
			goTyp:      new(types.Hash),
			solTyp:     "string",
			src:        new(types.Hash),
			wantEncErr: true,
		},

		// types.Hash <=> bytes
		{
			name:    "types.Hash<=>bytes",
			goTyp:   new(types.Hash),
			solTyp:  "bytes",
			src:     types.MustHexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
			wantDst: types.MustHexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		},

		// types.Hash <=> bytesX
		{
			name:    "types.Hash<=>bytesX",
			goTyp:   new(types.Hash),
			solTyp:  "bytes32",
			src:     types.MustHexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
			wantDst: types.MustHexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		},
		{
			name:       "types.Hash<=>bytesX#too-short",
			goTyp:      new(types.Hash),
			solTyp:     "bytes20",
			src:        types.MustHexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
			wantEncErr: true,
		},

		// types.Hash <=> address
		{
			name:       "types.Hash<=>address#encode-error",
			goTyp:      new(types.Hash),
			solTyp:     "address",
			src:        new(types.Hash),
			wantEncErr: true,
		},
		{
			name:       "types.Hash<=>address#decode-error",
			goTyp:      new(types.Hash),
			solTyp:     "address",
			src:        types.MustHexToAddress("0x1234567890123456789012345678901234567890"),
			wantDecErr: true,
		},

		// types.Bytes <=> intX
		{
			name:       "types.Bytes<=>intX#encode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "int",
			src:        types.Bytes([]byte{0x1}),
			wantEncErr: true,
		},
		{
			name:       "types.Bytes<=>intX#decode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "int",
			src:        "0x1",
			wantDecErr: true,
		},

		// types.Bytes <=> uintX
		{
			name:       "types.Bytes<=>uintX#encode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "uint",
			src:        types.Bytes([]byte{0x1}),
			wantEncErr: true,
		},
		{
			name:       "types.Bytes<=>uintX#decode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "uint",
			src:        "0x1",
			wantDecErr: true,
		},

		// types.Bytes <=> bool
		{
			name:       "types.Bytes<=>bool#encode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "bool",
			src:        types.Bytes([]byte{0x1}),
			wantEncErr: true,
		},
		{
			name:       "types.Bytes<=>bool#decode-error",
			goTyp:      new(types.Bytes),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// types.Bytes <=> string
		{
			name:    "types.Bytes<=>string",
			goTyp:   new(types.Bytes),
			solTyp:  "string",
			src:     types.Bytes([]byte{0x1}),
			wantDst: types.Bytes([]byte{0x1}),
		},

		// types.Bytes <=> bytes
		{
			name:    "types.Bytes<=>bytes",
			goTyp:   new(types.Bytes),
			solTyp:  "bytes",
			src:     types.Bytes([]byte{0x66, 0x6f, 0x6f}),
			wantDst: types.Bytes([]byte{0x66, 0x6f, 0x6f}),
		},

		// types.Bytes <=> bytesX
		{
			name:    "types.Bytes<=>bytesX",
			goTyp:   new(types.Bytes),
			solTyp:  "bytes4",
			src:     types.Bytes([]byte{0x66, 0x6f, 0x6f, 0x00}),
			wantDst: types.Bytes([]byte{0x66, 0x6f, 0x6f, 0x00}),
		},
		{
			name:       "types.Bytes<=>bytesX#too-short",
			goTyp:      new(types.Bytes),
			solTyp:     "bytes2",
			src:        types.Bytes([]byte{0x66, 0x6f, 0x6f}),
			wantEncErr: true,
		},

		// types.Bytes <=> address
		{
			name:    "types.Bytes<=>address",
			goTyp:   new(types.Bytes),
			solTyp:  "address",
			src:     types.Bytes([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}),
			wantDst: types.Bytes([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}),
		},
		{
			name:       "types.Bytes<=>address#too-short",
			goTyp:      new(types.Bytes),
			solTyp:     "address",
			src:        types.Bytes([]byte{0x12}),
			wantEncErr: true,
		},
		{
			name:       "types.Bytes<=>address#too-long",
			goTyp:      new(types.Bytes),
			solTyp:     "address",
			src:        types.Bytes([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12}),
			wantEncErr: true,
		},

		// types.Number <=> intX
		{
			name:    "types.Number<=>int",
			goTyp:   new(types.Number),
			solTyp:  "int",
			src:     types.BigIntToNumber(big.NewInt(0x1234)),
			wantDst: types.BigIntToNumber(big.NewInt(0x1234)),
		},
		{
			name:    "types.Number<=>int#negative",
			goTyp:   new(types.Number),
			solTyp:  "int",
			src:     types.BigIntToNumber(big.NewInt(-0x1234)),
			wantDst: types.BigIntToNumber(big.NewInt(-0x1234)),
		},
		{
			name:    "types.Number<=>int#too-smaller-type",
			goTyp:   new(types.Number),
			solTyp:  "int8",
			src:     types.BigIntToNumber(big.NewInt(42)),
			wantDst: types.BigIntToNumber(big.NewInt(42)),
		},
		{
			name:       "types.Number<=>int#overflow",
			goTyp:      new(types.Number),
			solTyp:     "int8",
			src:        types.BigIntToNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},

		// types.Number <=> uintX
		{
			name:    "types.Number<=>uint",
			goTyp:   new(types.Number),
			solTyp:  "uint",
			src:     types.BigIntToNumber(big.NewInt(0x1234)),
			wantDst: types.BigIntToNumber(big.NewInt(0x1234)),
		},
		{
			name:    "types.Number<=>uint#too-smaller-type",
			goTyp:   new(types.Number),
			solTyp:  "uint8",
			src:     types.BigIntToNumber(big.NewInt(42)),
			wantDst: types.BigIntToNumber(big.NewInt(42)),
		},
		{
			name:       "types.Number<=>uint#overflow",
			goTyp:      new(types.Number),
			solTyp:     "uint8",
			src:        types.BigIntToNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},

		// types.Number <=> bool
		{
			name:       "types.Number<=>bool#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "bool",
			src:        types.BigIntToNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>bool#decode-error",
			goTyp:      new(types.Number),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// types.Number <=> bytes
		{
			name:       "types.Number<=>bytes#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes",
			src:        types.BigIntToNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>bytes#decode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes",
			src:        []byte{0x12, 0x34},
			wantDecErr: true,
		},

		// types.Number <=> string
		{
			name:       "types.Number<=>string#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "string",
			src:        types.BigIntToNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>string#decode-error",
			goTyp:      new(int),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// types.Number <=> bytes
		{
			name:       "types.Number<=>bytes#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes",
			src:        types.BigIntToNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>bytes#decode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// types.Number <=> bytesX
		{
			name:    "types.Number<=>bytesX",
			goTyp:   new(types.Number),
			solTyp:  "bytes32",
			src:     types.BigIntToNumber(big.NewInt(1)),
			wantDst: types.BigIntToNumber(big.NewInt(1)),
		},
		{
			name:    "types.Number<=>bytesX#negative",
			goTyp:   new(types.Number),
			solTyp:  "bytes32",
			src:     types.BigIntToNumber(big.NewInt(-1)),
			wantDst: types.BigIntToNumber(big.NewInt(-1)),
		},
		{
			name:       "types.Number<=>bytesX#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes8",
			src:        types.BigIntToNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>bytesX#decode-error",
			goTyp:      new(types.Number),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// types.Number <=> address
		{
			name:       "types.Number<=>address#encode-error",
			goTyp:      new(types.Number),
			solTyp:     "address",
			src:        types.BigIntToNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.Number<=>address#decode-error",
			goTyp:      new(types.Number),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},

		// types.BlockNumber <=> intX
		{
			name:    "types.BlockNumber<=>int",
			goTyp:   new(types.BlockNumber),
			solTyp:  "int",
			src:     types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantDst: types.BigIntToBlockNumber(big.NewInt(0x1234)),
		},
		{
			name:       "types.BlockNumber<=>int#negative",
			goTyp:      new(types.BlockNumber),
			solTyp:     "int",
			src:        types.BigIntToBlockNumber(big.NewInt(-0x1234)),
			wantEncErr: true,
		},
		{
			name:    "types.BlockNumber<=>int#too-smaller-type",
			goTyp:   new(types.BlockNumber),
			solTyp:  "int8",
			src:     types.BigIntToBlockNumber(big.NewInt(42)),
			wantDst: types.BigIntToBlockNumber(big.NewInt(42)),
		},
		{
			name:       "types.BlockNumber<=>int#overflow",
			goTyp:      new(types.BlockNumber),
			solTyp:     "int8",
			src:        types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>int#Pending",
			goTyp:      new(types.BlockNumber),
			solTyp:     "int",
			src:        types.PendingBlockNumber,
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>int#Latest",
			goTyp:      new(types.BlockNumber),
			solTyp:     "int",
			src:        types.LatestBlockNumber,
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>int#Earliest",
			goTyp:      new(types.BlockNumber),
			solTyp:     "int",
			src:        types.EarliestBlockNumber,
			wantEncErr: true,
		},

		// types.BlockNumber <=> uintX
		{
			name:    "types.BlockNumber<=>uint",
			goTyp:   new(types.BlockNumber),
			solTyp:  "uint",
			src:     types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantDst: types.BigIntToBlockNumber(big.NewInt(0x1234)),
		},
		{
			name:    "types.BlockNumber<=>uint#too-smaller-type",
			goTyp:   new(types.BlockNumber),
			solTyp:  "uint8",
			src:     types.BigIntToBlockNumber(big.NewInt(42)),
			wantDst: types.BigIntToBlockNumber(big.NewInt(42)),
		},
		{
			name:       "types.BlockNumber<=>uint#overflow",
			goTyp:      new(types.BlockNumber),
			solTyp:     "uint8",
			src:        types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},

		// types.BlockNumber <=> bool
		{
			name:       "types.BlockNumber<=>bool#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bool",
			src:        types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>bool#decode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bool",
			src:        true,
			wantDecErr: true,
		},

		// types.BlockNumber <=> bytes
		{
			name:       "types.BlockNumber<=>bytes#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes",
			src:        types.BigIntToBlockNumber(big.NewInt(0x1234)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>bytes#decode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes",
			src:        []byte{0x12, 0x34},
			wantDecErr: true,
		},

		// types.BlockNumber <=> string
		{
			name:       "types.BlockNumber<=>string#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "string",
			src:        types.BigIntToBlockNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>string#decode-error",
			goTyp:      new(int),
			solTyp:     "string",
			src:        "0x1",
			wantDecErr: true,
		},

		// types.BlockNumber <=> bytes
		{
			name:       "types.BlockNumber<=>bytes#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes",
			src:        types.BigIntToBlockNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>bytes#decode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes",
			src:        []byte{0x1},
			wantDecErr: true,
		},

		// types.BlockNumber <=> bytesX
		{
			name:    "types.BlockNumber<=>bytesX",
			goTyp:   new(types.BlockNumber),
			solTyp:  "bytes32",
			src:     types.BigIntToBlockNumber(big.NewInt(1)),
			wantDst: types.BigIntToBlockNumber(big.NewInt(1)),
		},
		{
			name:       "types.BlockNumber<=>bytesX#negative",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes32",
			src:        types.BigIntToBlockNumber(big.NewInt(-1)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>bytesX#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes8",
			src:        types.BigIntToBlockNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>bytesX#decode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "bytes8",
			src:        []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantDecErr: true,
		},

		// types.BlockNumber <=> address
		{
			name:       "types.BlockNumber<=>address#encode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "address",
			src:        types.BigIntToBlockNumber(big.NewInt(1)),
			wantEncErr: true,
		},
		{
			name:       "types.BlockNumber<=>address#decode-error",
			goTyp:      new(types.BlockNumber),
			solTyp:     "address",
			src:        "0x1234567890123456789012345678901234567890",
			wantDecErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := MustParseType(tt.solTyp)
			cd, err := EncodeValue(typ, tt.src)
			if tt.wantEncErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			err = DecodeValue(typ, cd, &tt.goTyp)
			if tt.wantDecErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			switch tt.goTyp.(type) {
			case *big.Int:
				require.Equal(t, tt.goTyp.(*big.Int).String(), tt.wantDst.(*big.Int).String())
			default:
				assert.Equal(t, reflect.ValueOf(tt.goTyp).Elem().Interface(), tt.wantDst)
			}
		})
	}
}

func Test_fieldMapper(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"aB", "aB"},
		{"AB", "ab"},
		{"Ab", "ab"},
		{"abc", "abc"},
		{"ABC", "abc"},
		{"Abc", "abc"},
		{"ABc", "aBc"},
		{"Abcd", "abcd"},
		{"ABcd", "aBcd"},
		{"ABCd", "abCd"},
		{"ID", "id"},
		{"Id", "id"},
		{"UserID", "userID"},
		{"UserId", "userId"},
		{"DAPP", "dapp"},
		{"Dapp", "dapp"},
		{"DAPPName", "dappName"},
		{"DappName", "dappName"},
		{"DAPP1Name", "dapp1Name"},
		{"DAPP_Name", "dapp_Name"},
		{"I18NCode", "i18nCode"},
		{"Int32Num", "int32Num"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fieldMapper(tt.name))
		})
	}
}
