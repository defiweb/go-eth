package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
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

func TestNewKeyFromDirectory(t *testing.T) {
	t.Run("key-1", func(t *testing.T) {
		key, err := NewKeyFromDirectory("./testdata", "test123", types.MustAddressFromHex("0x2d800d93b065ce011af83f316cef9f0d005b0aa4"))
		require.NoError(t, err)
		assert.Equal(t, "0x2d800d93b065ce011af83f316cef9f0d005b0aa4", key.Address().String())
	})
	t.Run("key-2", func(t *testing.T) {
		key, err := NewKeyFromDirectory("./testdata", "testpassword", types.MustAddressFromHex("0x008aeeda4d805471df9b2a5b0f38a0c3bcba786b"))
		require.NoError(t, err)
		assert.Equal(t, "0x008aeeda4d805471df9b2a5b0f38a0c3bcba786b", key.Address().String())
	})
	t.Run("invalid-password", func(t *testing.T) {
		_, err := NewKeyFromDirectory("./testdata", "", types.MustAddressFromHex("0x2d800d93b065ce011af83f316cef9f0d005b0aa4"))
		require.Error(t, err)
	})
	t.Run("missing-key", func(t *testing.T) {
		_, err := NewKeyFromDirectory("./testdata", "", types.MustAddressFromHex("0x0000000000000000000000000000000000000000"))
		require.Error(t, err)
	})
}

func TestPrivateKey_JSON(t *testing.T) {
	t.Run("random", func(t *testing.T) {
		key1 := NewRandomKey()
		j, err := key1.JSON("test123", LightScryptN, LightScryptP)
		require.NoError(t, err)

		key2, err := NewKeyFromJSONContent(j, "test123")
		require.NoError(t, err)

		assert.Equal(t, key1.Address(), key2.Address())
	})
	t.Run("existing", func(t *testing.T) {
		key1, err := NewKeyFromJSON("./testdata/scrypt.json", "test123")
		require.NoError(t, err)

		j, err := key1.JSON("test123", LightScryptN, LightScryptP)
		require.NoError(t, err)

		key2, err := NewKeyFromJSONContent(j, "test123")
		require.NoError(t, err)

		assert.Equal(t, key1.Address(), key2.Address())
	})
}
