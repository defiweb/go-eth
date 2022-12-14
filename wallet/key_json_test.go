package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeyFromJSON(t *testing.T) {
	t.Run("scrypt", func(t *testing.T) {
		key, err := NewKeyFromJSON("./testdata/scrypt.json", "test123")
		require.NoError(t, err)
		assert.Equal(t, "0x2d800d93b065ce011af83f316cef9f0d005b0aa4", key.Address().String())
	})
	t.Run("pbkdf2", func(t *testing.T) {
		key, err := NewKeyFromJSON("./testdata/pbkdf2.json", "testpassword")
		require.NoError(t, err)
		assert.Equal(t, "0x008aeeda4d805471df9b2a5b0f38a0c3bcba786b", key.Address().String())
	})
}
