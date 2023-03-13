package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransaction_RLP(t1 *testing.T) {
	tests := []struct {
		tx   *Transaction
		want []byte
	}{
		// Empty transaction:
		{
			tx: (&Transaction{}).
				SetGasLimit(0).
				SetGasPrice(big.NewInt(0)).
				SetNonce(0).
				SetValue(big.NewInt(0)),
			want: hexutil.MustHexToBytes("c9808080808080808080"),
		},
		// Legacy transaction:
		{
			tx: (&Transaction{}).
				SetType(LegacyTxType).
				SetFrom(MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetGasPrice(big.NewInt(1000000000)).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")),
			want: hexutil.MustHexToBytes("f87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Access list transaction:
		{
			tx: (&Transaction{}).
				SetType(AccessListTxType).
				SetFrom(MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetGasPrice(big.NewInt(1000000000)).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetAccessList(AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
						},
					}}),
			want: hexutil.MustHexToBytes("01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Dynamic fee transaction:
		{
			tx: (&Transaction{}).
				SetType(DynamicFeeTxType).
				SetFrom(MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetMaxPriorityFeePerGas(big.NewInt(1000000000)).
				SetMaxFeePerGas(big.NewInt(2000000000)).
				SetAccessList(AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
						},
					},
				}),
			want: hexutil.MustHexToBytes("02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Dynamic fee transaction with no access list:
		{
			tx: (&Transaction{}).
				SetType(DynamicFeeTxType).
				SetFrom(MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetMaxPriorityFeePerGas(big.NewInt(1000000000)).
				SetMaxFeePerGas(big.NewInt(2000000000)),
			want: hexutil.MustHexToBytes("02f8770101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304c06fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Example from EIP-155:
		{
			tx: (&Transaction{}).
				SetType(LegacyTxType).
				SetChainID(1).
				SetTo(MustAddressFromHex("0x3535353535353535353535353535353535353535")).
				SetGasLimit(21000).
				SetGasPrice(big.NewInt(20000000000)).
				SetNonce(9).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(SignatureFromVRS(
					func() *big.Int {
						v, _ := new(big.Int).SetString("37", 10)
						return v
					}(),
					func() *big.Int {
						v, _ := new(big.Int).SetString("18515461264373351373200002665853028612451056578545711640558177340181847433846", 10)
						return v
					}(),
					func() *big.Int {
						v, _ := new(big.Int).SetString("46948507304638947509940763649030358759909902576025900602547168820602576006531", 10)
						return v
					}(),
				)),
			want: hexutil.MustHexToBytes("f86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"),
		},
	}
	for n, tt := range tests {
		t1.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			// Encode
			rlp, err := tt.tx.Raw()
			require.NoError(t1, err)
			assert.Equal(t1, tt.want, rlp)

			// Decode
			tx := new(Transaction)
			_, err = tx.DecodeRLP(rlp)
			require.NoError(t1, err)
			equalTx(t1, tt.tx, tx)
		})
	}
}

func equalTx(t *testing.T, expected, got *Transaction) {
	assert.Equal(t, expected.Type, got.Type)
	assert.Equal(t, expected.To, got.To)
	assert.Equal(t, expected.GasLimit, got.GasLimit)
	assert.Equal(t, expected.GasPrice, got.GasPrice)
	assert.Equal(t, expected.Input, got.Input)
	assert.Equal(t, expected.Nonce, got.Nonce)
	assert.Equal(t, expected.Value, got.Value)
	assert.Equal(t, expected.Signature, got.Signature)
	if expected.Type != LegacyTxType {
		assert.Equal(t, expected.ChainID, got.ChainID)
	}
	assert.Equal(t, expected.MaxPriorityFeePerGas, got.MaxPriorityFeePerGas)
	assert.Equal(t, expected.MaxFeePerGas, got.MaxFeePerGas)
	for i, accessTuple := range expected.AccessList {
		assert.Equal(t, accessTuple.Address, got.AccessList[i].Address)
		assert.Equal(t, accessTuple.StorageKeys, got.AccessList[i].StorageKeys)
	}
}
