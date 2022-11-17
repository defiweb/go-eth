package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fieldMapper(tt.name))
		})
	}
}
