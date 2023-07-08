package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseType(t *testing.T) {
	tests := []struct {
		sig     string
		want    string
		wantErr bool
	}{
		{sig: "uint256", want: "uint256"},
		{sig: "uint256[]", want: "uint256[]"},
		{sig: "(uint256 a, uint256 b)", want: "(uint256 a, uint256 b)"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := ParseType(tt.sig)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestParseStruct(t *testing.T) {
	tests := []struct {
		sig     string
		want    string
		wantErr bool
	}{
		{sig: "struct { uint256 a; }", want: "(uint256 a)"},
		{sig: "struct test { uint256 a; }", want: "(uint256 a)"}, // name is ignored
		{sig: "struct { uint256[] a; }", want: "(uint256[] a)"},
		{sig: "struct { uint256 a; uint256 b; }", want: "(uint256 a, uint256 b)"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := ParseStruct(tt.sig)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}
