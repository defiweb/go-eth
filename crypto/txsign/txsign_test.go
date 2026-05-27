package txsign

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/crypto/ecdsa"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func TestSign(t *testing.T) {
	t.Run("legacy", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(bytes.Repeat([]byte{0x01}, 32))
		tx := types.NewTransactionLegacy()
		tx.SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535"))
		tx.SetGasLimit(21000)
		tx.SetGasPrice(big.NewInt(20000000000))
		tx.SetNonce(9)
		tx.SetValue(big.NewInt(1000000000000000000))

		err := Sign(&ecdsa.PrivateKey{D: key.ToECDSA().D}, tx)
		require.NoError(t, err)
		txData := tx.GetSigningData()
		require.NotNil(t, txData.Signature)
		assert.Equal(t, "1b", txData.Signature.V.Text(16))
		assert.Equal(t, "2bfad43ba1b40e7f3ffb6342b1a6eecc700dd344fb0aba543aed5c10fd1a9470", txData.Signature.R.Text(16))
		assert.Equal(t, "615bff48c483d368ed4f6e327a6ddd8831e544d0ca08f1345433e4ed204f8537", txData.Signature.S.Text(16))
	})
	t.Run("dynamic-fee", func(t *testing.T) {
		key, _ := btcec.PrivKeyFromBytes(bytes.Repeat([]byte{0x01}, 32))
		tx := types.NewTransactionDynamicFee()
		tx.SetChainID(1)
		tx.SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535"))
		tx.SetGasLimit(21000)
		tx.SetMaxFeePerGas(big.NewInt(20000000000))
		tx.SetMaxPriorityFeePerGas(big.NewInt(20000000000))
		tx.SetNonce(9)
		tx.SetValue(big.NewInt(1000000000000000000))

		err := Sign(&ecdsa.PrivateKey{D: key.ToECDSA().D}, tx)
		require.NoError(t, err)
		txData := tx.GetSigningData()
		require.NotNil(t, txData.Signature)
		assert.Equal(t, "0", txData.Signature.V.Text(16))
		assert.Equal(t, "62072d055f9ceb871a47f2d81aeb5aa34df50c625da16c6d0d57d232fa3cd152", txData.Signature.R.Text(16))
		assert.Equal(t, "57fd88df7c85076f5729493be7e87f51b618a78bc89441ed741bdfdb9d1d5572", txData.Signature.S.Text(16))
	})
}

func TestRecover(t *testing.T) {
	t.Run("legacy", func(t *testing.T) {
		tx := types.NewTransactionLegacy()
		tx.SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535"))
		tx.SetGasLimit(21000)
		tx.SetGasPrice(big.NewInt(20000000000))
		tx.SetNonce(9)
		tx.SetValue(big.NewInt(1000000000000000000))
		txData := tx.GetSigningData()
		txData.SetSignature(types.SignatureFromVRS(
			hexutil.MustHexToBigInt("1b"),
			hexutil.MustHexToBigInt("2bfad43ba1b40e7f3ffb6342b1a6eecc700dd344fb0aba543aed5c10fd1a9470"),
			hexutil.MustHexToBigInt("615bff48c483d368ed4f6e327a6ddd8831e544d0ca08f1345433e4ed204f8537"),
		))

		addr, err := Recover(tx)
		require.NoError(t, err)
		require.NotNil(t, addr)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
	t.Run("dynamic-fee", func(t *testing.T) {
		tx := types.NewTransactionDynamicFee()
		tx.SetChainID(1)
		tx.SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535"))
		tx.SetGasLimit(21000)
		tx.SetMaxFeePerGas(big.NewInt(20000000000))
		tx.SetMaxPriorityFeePerGas(big.NewInt(20000000000))
		tx.SetNonce(9)
		tx.SetValue(big.NewInt(1000000000000000000))
		txData := tx.GetSigningData()
		txData.SetSignature(types.SignatureFromVRS(
			hexutil.MustHexToBigInt("0"),
			hexutil.MustHexToBigInt("62072d055f9ceb871a47f2d81aeb5aa34df50c625da16c6d0d57d232fa3cd152"),
			hexutil.MustHexToBigInt("57fd88df7c85076f5729493be7e87f51b618a78bc89441ed741bdfdb9d1d5572"),
		))

		addr, err := Recover(tx)
		require.NoError(t, err)
		require.NotNil(t, addr)
		assert.Equal(t, "0x1a642f0e3c3af545e7acbd38b07251b3990914f1", addr.String())
	})
}
