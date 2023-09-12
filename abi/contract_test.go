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

	assert.NotNil(t, abi.Methods["Foo"])
	assert.NotNil(t, abi.Methods["Bar"])
	assert.NotNil(t, abi.Constructor)
	assert.NotNil(t, abi.Events["EventA"])
	assert.NotNil(t, abi.Events["EventB"])
	assert.NotNil(t, abi.Errors["ErrorA"])

	assert.Equal(t, "function Foo(uint256 a) returns (uint256)", abi.Methods["Foo"].String())
	assert.Equal(t, "function Bar((bytes32 A, bytes32 B)[2][2] a) returns (uint256[2][2])", abi.Methods["Bar"].String())
	assert.Equal(t, "constructor(uint256 a)", abi.Constructor.String())
	assert.Equal(t, "event EventA(uint256 indexed a, uint256 b)", abi.Events["EventA"].String())
	assert.Equal(t, "event EventB(uint256 indexed a, uint256 b) anonymous", abi.Events["EventB"].String())
	assert.Equal(t, "error ErrorA(uint256 a, uint256 b)", abi.Errors["ErrorA"].String())
}

func TestABI_ParseSignatures(t *testing.T) {
	c, err := ParseSignatures(
		"foo(uint256)",
		"function bar(uint256) returns (uint256)",
		"constructor(uint256)",
		"event baz(uint256)",
		"error qux(uint256)",
	)
	require.NoError(t, err)
	assert.NotNil(t, c.Methods["foo"])
	assert.NotNil(t, c.Methods["bar"])
	assert.NotNil(t, c.MethodsBySignature["foo(uint256)"])
	assert.NotNil(t, c.MethodsBySignature["bar(uint256)"])
	assert.NotNil(t, c.Constructor)
	assert.NotNil(t, c.Events["baz"])
	assert.NotNil(t, c.Errors["qux"])
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
