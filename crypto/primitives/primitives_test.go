package primitives

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateKey_IsZero(t *testing.T) {
	assert.True(t, PrivateKey{}.IsZero())
	assert.False(t, PrivateKey{31: 1}.IsZero())
}

func TestPrivateKey_Zero(t *testing.T) {
	k := PrivateKey{0: 1, 31: 2}
	k.Zero()
	assert.True(t, k.IsZero())
}

func TestPrivateKey_Bytes(t *testing.T) {
	k := PrivateKey{0: 1, 31: 2}
	b := k.Bytes()
	require.Len(t, b, PrivateKeySize)
	assert.Equal(t, k[:], b)

	// The returned slice must not alias the key.
	b[0] = 0xff
	assert.Equal(t, byte(1), k[0])
}

func TestPrivateKey_Scalar(t *testing.T) {
	k := PrivateKey{31: 5}
	assert.Equal(t, int64(5), k.Scalar().Int64())
}

// TestPrivateKey_NoDisclose verifies that the key material is not disclosed by
// the formatting and marshaling methods.
func TestPrivateKey_NoDisclose(t *testing.T) {
	k := PrivateKey{0: 0xde, 1: 0xad, 30: 0xbe, 31: 0xef}

	assert.Equal(t, "PrivateKey(redacted)", k.String())
	assert.NotContains(t, fmt.Sprintf("%v", k), "de")
	assert.NotContains(t, fmt.Sprintf("%s", k), "de")
	assert.NotContains(t, fmt.Sprintf("%#v", k), "de")

	_, err := k.MarshalText()
	assert.Error(t, err)

	_, err = k.MarshalJSON()
	assert.Error(t, err)

	_, err = json.Marshal(k)
	assert.Error(t, err)

	_, err = json.Marshal(struct{ Key PrivateKey }{k})
	assert.Error(t, err)
}

func TestPublicKey_IsZero(t *testing.T) {
	assert.True(t, PublicKey{}.IsZero())
	assert.False(t, PublicKey{63: 1}.IsZero())
}

func TestPublicKey_Bytes(t *testing.T) {
	k := PublicKey{0: 1, 63: 2}
	b := k.Bytes()
	require.Len(t, b, PublicKeySize)
	assert.Equal(t, k[:], b)

	// The returned slice must not alias the key.
	b[0] = 0xff
	assert.Equal(t, byte(1), k[0])
}

func TestPublicKey_XY(t *testing.T) {
	k := PublicKey{31: 7, 63: 9}
	assert.Equal(t, int64(7), k.X().Int64())
	assert.Equal(t, int64(9), k.Y().Int64())
}
