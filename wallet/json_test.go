package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKeyFromJSON(t *testing.T) {
	key, err := NewKeyFromJSON("./testdata/1.json", "test123")
	require.NoError(t, err)
	assert.Equal(t, "0x2d800d93b065ce011af83f316cef9f0d005b0aa4", key.Address().String())
}
