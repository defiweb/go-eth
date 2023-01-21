package abi

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestABI_ParseJSON(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/abi.json")
	require.NoError(t, err)

	abi, err := ParseJSON(data)
	require.NoError(t, err)

	assert.NotNil(t, abi.Method("Foo"))
	assert.NotNil(t, abi.Method("Bar"))
	assert.NotNil(t, abi.Constructor())
	assert.NotNil(t, abi.Event("EventA"))
	assert.NotNil(t, abi.Event("EventB"))
	assert.NotNil(t, abi.Error("ErrorA"))

	assert.Equal(t, "function Foo(uint256 a) returns (uint256)", abi.Method("Foo").String())
	assert.Equal(t, "function Bar((bytes32 A, bytes32 B)[2][2] a) returns (uint256[2][2])", abi.Method("Bar").String())
	assert.Equal(t, "constructor(uint256 a)", abi.Constructor().String())
	assert.Equal(t, "event EventA(uint256 indexed a, uint256 b)", abi.Event("EventA").String())
	assert.Equal(t, "event EventB(uint256 indexed a, uint256 b) anonymous", abi.Event("EventB").String())
	assert.Equal(t, "error ErrorA(uint256 a, uint256 b)", abi.Error("ErrorA").String())
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
	assert.NotNil(t, c.Method("foo"))
	assert.NotNil(t, c.Method("bar"))
	assert.NotNil(t, c.MethodBySignature("foo(uint256)"))
	assert.NotNil(t, c.MethodBySignature("bar(uint256)"))
	assert.NotNil(t, c.Constructor())
	assert.NotNil(t, c.Event("baz"))
	assert.NotNil(t, c.Error("qux"))
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
		parseArrays(typ)
	})
}
