package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestABI_LoadJSON(t *testing.T) {
	abi, err := LoadJSON("testdata/abi.json")

	require.NoError(t, err)

	require.NotNil(t, abi.Types["Status"])
	require.NotNil(t, abi.Types["Struct"])
	require.NotNil(t, abi.Types["CustomUint"])
	require.NotNil(t, abi.Events["EventA"])
	require.NotNil(t, abi.Events["EventB"])
	require.NotNil(t, abi.Events["EventC"])
	require.NotNil(t, abi.Errors["ErrorA"])
	require.NotNil(t, abi.Constructor)
	require.NotNil(t, abi.Methods["Foo"])
	require.NotNil(t, abi.Methods["Bar"])
	require.NotNil(t, abi.Methods["structField"])
	require.NotNil(t, abi.Methods["structsMapping"])
	require.NotNil(t, abi.Methods["structsArray"])

	assert.Equal(t, "Status", abi.Types["Status"].String())
	assert.Equal(t, "Struct", abi.Types["Struct"].String())
	assert.Equal(t, "CustomUint", abi.Types["CustomUint"].String())
	assert.Equal(t, "event EventA(uint256 indexed a, string b)", abi.Events["EventA"].String())
	assert.Equal(t, "event EventB(uint256 indexed a, string indexed b)", abi.Events["EventB"].String())
	assert.Equal(t, "event EventC(uint256 indexed a, string b) anonymous", abi.Events["EventC"].String())
	assert.Equal(t, "error ErrorA(uint256 a, uint256 b)", abi.Errors["ErrorA"].String())
	assert.Equal(t, "constructor(CustomUint a)", abi.Constructor.String())
	assert.Equal(t, "function Foo(CustomUint a) nonpayable returns (CustomUint)", abi.Methods["Foo"].String())
	assert.Equal(t, "function Bar(Struct[2][2] a) nonpayable returns (uint8[2][2])", abi.Methods["Bar"].String())
	assert.Equal(t, "function structField() view returns (bytes32 A, bytes32 B, Status status)", abi.Methods["structField"].String())

	assert.Equal(t, "uint8", abi.Types["Status"].CanonicalType())
	assert.Equal(t, "(bytes32,bytes32,uint8)", abi.Types["Struct"].CanonicalType())
	assert.Equal(t, "uint256", abi.Types["CustomUint"].CanonicalType())
}

func TestABI_ParseSignatures(t *testing.T) {
	abi, err := ParseSignatures(
		`uint8 Status`,
		`struct Struct { bytes32 A; bytes32 B; Status status;}`,
		`uint256 CustomUint`,
		`event EventA(uint256 indexed a, string b)`,
		`event EventB(uint256 indexed a, string indexed b)`,
		`event EventC(uint256 indexed a, string b) anonymous`,
		`error ErrorA(uint256 a, uint256 b)`,
		`constructor(CustomUint a)`,
		`function Foo(CustomUint a) nonpayable returns (CustomUint)`,
		`function Bar(Struct[2][2] a) nonpayable returns (uint8[2][2])`,
	)

	require.NoError(t, err)
	require.NotNil(t, abi.Types["Status"])
	require.NotNil(t, abi.Types["Struct"])
	require.NotNil(t, abi.Types["CustomUint"])
	require.NotNil(t, abi.Events["EventA"])
	require.NotNil(t, abi.Events["EventB"])
	require.NotNil(t, abi.Events["EventC"])
	require.NotNil(t, abi.Errors["ErrorA"])
	require.NotNil(t, abi.Constructor)
	require.NotNil(t, abi.Methods["Foo"])
	require.NotNil(t, abi.Methods["Bar"])

	assert.Equal(t, "Status", abi.Types["Status"].String())
	assert.Equal(t, "Struct", abi.Types["Struct"].String())
	assert.Equal(t, "CustomUint", abi.Types["CustomUint"].String())
	assert.Equal(t, "event EventA(uint256 indexed a, string b)", abi.Events["EventA"].String())
	assert.Equal(t, "event EventB(uint256 indexed a, string indexed b)", abi.Events["EventB"].String())
	assert.Equal(t, "event EventC(uint256 indexed a, string b) anonymous", abi.Events["EventC"].String())
	assert.Equal(t, "error ErrorA(uint256 a, uint256 b)", abi.Errors["ErrorA"].String())
	assert.Equal(t, "constructor(CustomUint a)", abi.Constructor.String())
	assert.Equal(t, "function Foo(CustomUint a) nonpayable returns (CustomUint)", abi.Methods["Foo"].String())
	assert.Equal(t, "function Bar(Struct[2][2] a) nonpayable returns (uint8[2][2])", abi.Methods["Bar"].String())

	assert.Equal(t, "uint8", abi.Types["Status"].CanonicalType())
	assert.Equal(t, "(bytes32,bytes32,uint8)", abi.Types["Struct"].CanonicalType())
	assert.Equal(t, "uint256", abi.Types["CustomUint"].CanonicalType())
}

func TestContract_IsError(t *testing.T) {
	c, err := ParseSignatures(
		"error foo(uint256)",
	)
	require.NoError(t, err)

	assert.True(t, c.IsError(hexutil.MustHexToBytes("0x08c379a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003666f6f0000000000000000000000000000000000000000000000000000000000")))
	assert.True(t, c.IsError(hexutil.MustHexToBytes("0x4e487b710000000000000000000000000000000000000000000000000000000000000020")))
	assert.True(t, c.IsError(hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c")))
	assert.False(t, c.IsError(hexutil.MustHexToBytes("0xaabbccdd000000000000000000000000000000000000000000000000000000000000012c")))
}

func TestContract_ToError(t *testing.T) {
	c, err := ParseSignatures("error foo(uint256)")
	require.NoError(t, err)

	// Revert
	revertErr := c.ToError(hexutil.MustHexToBytes("0x08c379a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003666f6f0000000000000000000000000000000000000000000000000000000000"))
	require.NotNil(t, revertErr)
	assert.Equal(t, "revert: foo", revertErr.Error())

	// Panic
	panicErr := c.ToError(hexutil.MustHexToBytes("0x4e487b710000000000000000000000000000000000000000000000000000000000000020"))
	require.NotNil(t, panicErr)
	assert.Equal(t, "panic: 32", panicErr.Error())

	// Custom error
	customErr := c.ToError(hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c"))
	require.NotNil(t, customErr)
	assert.Equal(t, "error: foo", customErr.Error())
}

func TestContract_RegisterTypes(t *testing.T) {
	abi := NewABI()

	c, err := abi.ParseSignatures(
		`uint8 Status`,
		`struct Struct { bytes32 A; bytes32 B; Status status;}`,
	)

	require.NoError(t, err)

	c.RegisterTypes(abi)
	assert.Equal(t, "Status", abi.Types["Status"].String())
	assert.Equal(t, "Struct", abi.Types["Struct"].String())
}

func Test_parseArrays(t *testing.T) {
	tests := []struct {
		typ        string
		wantName   string
		wantArrays []int
		wantErr    assert.ErrorAssertionFunc
	}{
		{typ: "uint256", wantName: "uint256", wantArrays: nil, wantErr: assert.NoError},
		{typ: "uint256[]", wantName: "uint256", wantArrays: []int{-1}, wantErr: assert.NoError},
		{typ: "uint256[][]", wantName: "uint256", wantArrays: []int{-1, -1}, wantErr: assert.NoError},
		{typ: "uint256[2]", wantName: "uint256", wantArrays: []int{2}, wantErr: assert.NoError},
		{typ: "uint256[2][3]", wantName: "uint256", wantArrays: []int{2, 3}, wantErr: assert.NoError},
		{typ: "uint256[][3]", wantName: "uint256", wantArrays: []int{-1, 3}, wantErr: assert.NoError},
		{typ: "uint256[2][]", wantName: "uint256", wantArrays: []int{2, -1}, wantErr: assert.NoError},
		{typ: "uint256[", wantName: "", wantArrays: nil, wantErr: assert.Error},     // missing ]
		{typ: "uint256[2", wantName: "", wantArrays: nil, wantErr: assert.Error},    // missing ]
		{typ: "uint256[2][", wantName: "", wantArrays: nil, wantErr: assert.Error},  // missing ]
		{typ: "uint256[2][3", wantName: "", wantArrays: nil, wantErr: assert.Error}, // missing ]
		{typ: "uint256[]]", wantName: "", wantArrays: nil, wantErr: assert.Error},   // missing [
		{typ: "uint256[]]]", wantName: "", wantArrays: nil, wantErr: assert.Error},  // invalid syntax
		{typ: "uint256[[[]", wantName: "", wantArrays: nil, wantErr: assert.Error},  // invalid syntax
		{typ: "uint256[-1]", wantName: "", wantArrays: nil, wantErr: assert.Error},  // negative size
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			gotName, gotArrays, err := parseArrays(tt.typ)
			if !tt.wantErr(t, err, fmt.Sprintf("parseArrays(%v)", tt.typ)) {
				return
			}
			assert.Equalf(t, tt.wantName, gotName, "parseArrays(%v)", tt.typ)
			assert.Equalf(t, tt.wantArrays, gotArrays, "parseArrays(%v)", tt.typ)
		})
	}
}

func Fuzz_parseArrays(f *testing.F) {
	for _, typ := range []string{
		"uint256",
		"[",
		"]",
		"[]",
		"1",
		"-",
	} {
		f.Add(typ)
	}
	f.Fuzz(func(t *testing.T, typ string) {
		_, _, _ = parseArrays(typ)
	})
}
