package wallet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet_Mnemonic(t *testing.T) {
	tests := []struct {
		account uint32
		index   uint32
		addr    string
	}{
		{0, 0, "0x02941ca660485ba7dc196b510d9a6192c2648709"},
		{0, 1, "0xd050d1f66eb5ed560079754f3c1623b369a1a5ee"},
		{1, 0, "0x7931220c3f0ee7efb9e323de4b9053e8aba3ff30"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			key, err := NewKeyFromMnemonic(
				"gravity trophy shrimp suspect sheriff avocado label trust dove tragic pitch title network myself spell task protect smooth sword diary brain blossom under bulb",
				"fJF*(SDF*(*@J!)(SU*(D*F&^&TYSDFHL#@HO*&O",
				tt.account,
				tt.index,
			)
			require.NoError(t, err)
			assert.Equal(t, tt.addr, key.Address().String())

		})
	}
}

func TestParseDerivationPath(t *testing.T) {
	// Based on test cases from github.com/ethereum/go-ethereum/blob/master/accounts/hd_test.go
	tests := []struct {
		input  string
		output DerivationPath
	}{
		// Plain absolute derivation paths.
		{"m/44'/60'/0'/0", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}},
		{"m/44'/60'/0'/128", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 128}},
		{"m/44'/60'/0'/0'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 0}},
		{"m/44'/60'/0'/128'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 128}},
		{"m/2147483692/2147483708/2147483648/0", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}},
		{"m/2147483692/2147483708/2147483648/2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 0}},

		// Plain relative derivation paths.
		{"0", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}},
		{"128", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 128}},
		{"0'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 0}},
		{"128'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 128}},
		{"2147483648", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 0}},

		// Hexadecimal absolute derivation paths.
		{"m/0x2C'/0x3c'/0x00'/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}},
		{"m/0x2C'/0x3c'/0x00'/0x80", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 128}},
		{"m/0x2C'/0x3c'/0x00'/0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 0}},
		{"m/0x2C'/0x3c'/0x00'/0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 128}},
		{"m/0x8000002C/0x8000003c/0x80000000/0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}},
		{"m/0x8000002C/0x8000003c/0x80000000/0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0x80000000 + 0}},

		// Hexadecimal relative derivation paths.
		{"0x00", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}},
		{"0x80", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 128}},
		{"0x00'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 0}},
		{"0x80'", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 128}},
		{"0x80000000", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0x80000000 + 0}},

		// Weird inputs just to ensure they work.
		{"	m  /   44			'\n/\n   60	\n\n\t'   /\n0 ' /\t\t	0", DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}},

		// Invalid derivation paths
		{"", nil},              // Empty relative derivation path.
		{"m", nil},             // Empty absolute derivation path.
		{"m/", nil},            // Missing last derivation component.
		{"/44'/60'/0'/0", nil}, // Absolute path without m prefix, might be user error.
		{"m/2147483648'", nil}, // Overflows 32 bit integer.
		{"m/-1'", nil},         // Cannot contain negative number.
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseDerivationPath(tt.input)
			if tt.output == nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, got)
			}
		})
	}
}

func FuzzParseDerivationPath(f *testing.F) {
	for _, input := range []string{
		"m",
		"/",
		"0x",
		"44",
		"2147483692",
		"2147483648",
		"'",
		" ",
		"\t",
		"\r",
		"\n",
	} {
		f.Add([]byte(input))
	}
	f.Fuzz(func(t *testing.T, input []byte) {
		ParseDerivationPath(string(input))
	})
}
