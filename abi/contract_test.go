package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
