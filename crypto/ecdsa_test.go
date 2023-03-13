package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func TestEthereumSigner_SignHash(t *testing.T) {
	key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
	signature, err := ecSignHash((*ecdsa.PrivateKey)(key), types.MustHashFromBytes(bytes.Repeat([]byte{0x02}, 32)))

	require.NoError(t, err)
	require.NotNil(t, signature)
	assert.Equal(t, "1b", signature.V.Text(16))
	assert.Equal(t, "97ef30233ead25d10f7bb2bf9eaf571a16f2deb33a75f20819284f0cb8ff3cc1", signature.R.Text(16))
	assert.Equal(t, "4870ca05940199c113b4dc77866f001702691cde269f6835581e7aea1ead2660", signature.S.Text(16))
}

func TestEthereumSigner_SignMessage(t *testing.T) {
	key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
	signature, err := ecSignMessage((*ecdsa.PrivateKey)(key), []byte("hello world"))

	require.NoError(t, err)
	require.NotNil(t, signature)
	assert.Equal(t, "1b", signature.V.Text(16))
	assert.Equal(t, "f2b67e452d18ce781203f10380ea5a2726494162c49c495069cf99118bcf199", signature.R.Text(16))
	assert.Equal(t, "51601fe3219055482c45a14bf616c3e2bc7914c953f438627de2aa541eef61b5", signature.S.Text(16))
}

func TestEthereumSigner_SignTransaction(t *testing.T) {
	t.Run("legacy", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
		tx := (&types.Transaction{}).
			SetType(types.LegacyTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000))
		err := ecSignTransaction((*ecdsa.PrivateKey)(key), tx)

		require.NoError(t, err)
		assert.Equal(t, "1b", tx.Signature.V.Text(16))
		assert.Equal(t, "2bfad43ba1b40e7f3ffb6342b1a6eecc700dd344fb0aba543aed5c10fd1a9470", tx.Signature.R.Text(16))
		assert.Equal(t, "615bff48c483d368ed4f6e327a6ddd8831e544d0ca08f1345433e4ed204f8537", tx.Signature.S.Text(16))
	})
	t.Run("legacy-eip155", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
		tx := (&types.Transaction{}).
			SetType(types.LegacyTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000)).
			SetChainID(1337)
		err := ecSignTransaction((*ecdsa.PrivateKey)(key), tx)

		require.NoError(t, err)
		assert.Equal(t, "a95", tx.Signature.V.Text(16))
		assert.Equal(t, "14702a15dd7739397f25e3902a0c2bf6989e93888201139aac2c67a8f33a2f3f", tx.Signature.R.Text(16))
		assert.Equal(t, "4a10ba6cf47ace7e3c847e38583f5b1e1c7d8a862f4b43cd74480a03007363f7", tx.Signature.S.Text(16))
	})
	t.Run("access-list", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
		tx := (&types.Transaction{}).
			SetType(types.AccessListTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000))
		err := ecSignTransaction((*ecdsa.PrivateKey)(key), tx)

		require.NoError(t, err)
		assert.Equal(t, "1", tx.Signature.V.Text(16))
		assert.Equal(t, "dc1fcd0c6f56eddc8dbe70635690cce521276b8a6e167f8e57e4064db8a5738e", tx.Signature.R.Text(16))
		assert.Equal(t, "2743f261c001ee472c9664258708eaf849fc85623ee337d2018d37fc6f397d8c", tx.Signature.S.Text(16))
	})
	t.Run("dynamic-fee", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(s256, bytes.Repeat([]byte{0x01}, 32))
		tx := (&types.Transaction{}).
			SetType(types.DynamicFeeTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetMaxFeePerGas(big.NewInt(20000000000)).
			SetMaxPriorityFeePerGas(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000))
		err := ecSignTransaction((*ecdsa.PrivateKey)(key), tx)

		require.NoError(t, err)
		assert.Equal(t, "0", tx.Signature.V.Text(16))
		assert.Equal(t, "62072d055f9ceb871a47f2d81aeb5aa34df50c625da16c6d0d57d232fa3cd152", tx.Signature.R.Text(16))
		assert.Equal(t, "57fd88df7c85076f5729493be7e87f51b618a78bc89441ed741bdfdb9d1d5572", tx.Signature.S.Text(16))
	})
}

func TestEthereumSigner_RecoverHash(t *testing.T) {
	addr, err := ecRecoverHash(
		types.MustHashFromBytes(bytes.Repeat([]byte{0x02}, 32)),
		types.SignatureFromVRS(
			hexutil.MustHexToBigInt("1b"),
			hexutil.MustHexToBigInt("97ef30233ead25d10f7bb2bf9eaf571a16f2deb33a75f20819284f0cb8ff3cc1"),
			hexutil.MustHexToBigInt("4870ca05940199c113b4dc77866f001702691cde269f6835581e7aea1ead2660"),
		),
	)

	require.NoError(t, err)
	assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
}

func TestEthereumSigner_RecoverMessage(t *testing.T) {
	addr, err := ecRecoverMessage(
		[]byte("hello world"),
		types.SignatureFromVRS(
			hexutil.MustHexToBigInt("1b"),
			hexutil.MustHexToBigInt("f2b67e452d18ce781203f10380ea5a2726494162c49c495069cf99118bcf199"),
			hexutil.MustHexToBigInt("51601fe3219055482c45a14bf616c3e2bc7914c953f438627de2aa541eef61b5"),
		),
	)

	require.NoError(t, err)
	assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
}

func TestEthereumSigner_RecoverTransaction(t *testing.T) {
	t.Run("legacy", func(t *testing.T) {
		tx := (&types.Transaction{}).
			SetType(types.LegacyTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000)).
			SetSignature(types.SignatureFromVRS(
				hexutil.MustHexToBigInt("1b"),
				hexutil.MustHexToBigInt("2bfad43ba1b40e7f3ffb6342b1a6eecc700dd344fb0aba543aed5c10fd1a9470"),
				hexutil.MustHexToBigInt("615bff48c483d368ed4f6e327a6ddd8831e544d0ca08f1345433e4ed204f8537"),
			))
		addr, err := ecRecoverTransaction(tx)

		require.NoError(t, err)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
	t.Run("legacy-eip155", func(t *testing.T) {
		tx := (&types.Transaction{}).
			SetType(types.LegacyTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000)).
			SetChainID(1337).
			SetSignature(types.SignatureFromVRS(
				hexutil.MustHexToBigInt("a95"),
				hexutil.MustHexToBigInt("14702a15dd7739397f25e3902a0c2bf6989e93888201139aac2c67a8f33a2f3f"),
				hexutil.MustHexToBigInt("4a10ba6cf47ace7e3c847e38583f5b1e1c7d8a862f4b43cd74480a03007363f7"),
			))
		addr, err := ecRecoverTransaction(tx)

		require.NoError(t, err)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
	t.Run("access-list", func(t *testing.T) {
		tx := (&types.Transaction{}).
			SetType(types.AccessListTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetGasPrice(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000)).
			SetSignature(types.SignatureFromVRS(
				hexutil.MustHexToBigInt("1"),
				hexutil.MustHexToBigInt("dc1fcd0c6f56eddc8dbe70635690cce521276b8a6e167f8e57e4064db8a5738e"),
				hexutil.MustHexToBigInt("2743f261c001ee472c9664258708eaf849fc85623ee337d2018d37fc6f397d8c"),
			))
		addr, err := ecRecoverTransaction(tx)

		require.NoError(t, err)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
	t.Run("dynamic-fee", func(t *testing.T) {
		tx := (&types.Transaction{}).
			SetType(types.DynamicFeeTxType).
			SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
			SetGasLimit(21000).
			SetMaxFeePerGas(big.NewInt(20000000000)).
			SetMaxPriorityFeePerGas(big.NewInt(20000000000)).
			SetNonce(9).
			SetValue(big.NewInt(1000000000000000000)).
			SetSignature(types.SignatureFromVRS(
				hexutil.MustHexToBigInt("0"),
				hexutil.MustHexToBigInt("62072d055f9ceb871a47f2d81aeb5aa34df50c625da16c6d0d57d232fa3cd152"),
				hexutil.MustHexToBigInt("57fd88df7c85076f5729493be7e87f51b618a78bc89441ed741bdfdb9d1d5572"),
			))
		addr, err := ecRecoverTransaction(tx)

		require.NoError(t, err)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
}
