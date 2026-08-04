package ecdsa

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/crypto/primitives"
	"github.com/defiweb/go-eth/hexutil"
)

// testPrivateKey returns a deterministic private key (32 bytes of 0x01) used
// across the tests.
func testPrivateKey() primitives.PrivateKey {
	return primitives.PrivateKey([32]byte(bytes.Repeat([]byte{0x01}, 32)))
}

func TestSignHash(t *testing.T) {
	hash := primitives.Hash{}
	copy(hash[:], bytes.Repeat([]byte{0x02}, 32))
	signature, err := SignHash(testPrivateKey(), hash)

	require.NoError(t, err)
	require.NotNil(t, signature)
	assert.Equal(t, "0", signature.V.Text(16))
	assert.Equal(t, "97ef30233ead25d10f7bb2bf9eaf571a16f2deb33a75f20819284f0cb8ff3cc1", signature.R.Text(16))
	assert.Equal(t, "4870ca05940199c113b4dc77866f001702691cde269f6835581e7aea1ead2660", signature.S.Text(16))
}

func TestSignMessage(t *testing.T) {
	signature, err := SignMessage(testPrivateKey(), []byte("hello world"))

	require.NoError(t, err)
	require.NotNil(t, signature)
	assert.Equal(t, "1b", signature.V.Text(16))
	assert.Equal(t, "f2b67e452d18ce781203f10380ea5a2726494162c49c495069cf99118bcf199", signature.R.Text(16))
	assert.Equal(t, "51601fe3219055482c45a14bf616c3e2bc7914c953f438627de2aa541eef61b5", signature.S.Text(16))
}

func TestRecoverHash(t *testing.T) {
	hash := primitives.Hash{}
	copy(hash[:], bytes.Repeat([]byte{0x02}, 32))
	signature := primitives.Signature{
		V: hexutil.MustHexToBigInt("1b"),
		R: hexutil.MustHexToBigInt("97ef30233ead25d10f7bb2bf9eaf571a16f2deb33a75f20819284f0cb8ff3cc1"),
		S: hexutil.MustHexToBigInt("4870ca05940199c113b4dc77866f001702691cde269f6835581e7aea1ead2660"),
	}
	addr, err := RecoverHash(hash, signature)

	require.NoError(t, err)
	require.NotNil(t, addr)
	assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", hexutil.BytesToHex(addr[:]))
}

func TestRecoverMessage(t *testing.T) {
	signature := primitives.Signature{
		V: hexutil.MustHexToBigInt("1b"),
		R: hexutil.MustHexToBigInt("f2b67e452d18ce781203f10380ea5a2726494162c49c495069cf99118bcf199"),
		S: hexutil.MustHexToBigInt("51601fe3219055482c45a14bf616c3e2bc7914c953f438627de2aa541eef61b5"),
	}
	addr, err := RecoverMessage([]byte("hello world"), signature)

	require.NoError(t, err)
	require.NotNil(t, addr)
	assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", hexutil.BytesToHex(addr[:]))
}

func TestPublicKeyToAddress(t *testing.T) {
	addr := PublicKeyToAddress(PrivateKeyToPublicKey(testPrivateKey()))

	assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", hexutil.BytesToHex(addr[:]))
}
