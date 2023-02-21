package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransaction_RLP(t1 *testing.T) {
	tests := []struct {
		tx   *Transaction
		want []byte
	}{
		{
			tx: &Transaction{
				Type:      LegacyTxType,
				From:      HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:        HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:       100000,
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
			},
			want: hexutil.MustHexToBytes("f87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		{
			tx: &Transaction{
				Type:      AccessListTxType,
				From:      HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:        HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:       100000,
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
				AccessList: AccessList{
					AccessTuple{
						Address: MustHexToAddress("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHexToHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: hexutil.MustHexToBytes("01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:                   HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  100000,
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
				AccessList: AccessList{
					AccessTuple{
						Address: MustHexToAddress("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHexToHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: hexutil.MustHexToBytes("02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:                   HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  100000,
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
			},
			want: hexutil.MustHexToBytes("02f8770101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304c06fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
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
			assert.Equal(t1, tt.tx.Type, tx.Type)
			assert.Equal(t1, tt.tx.To, tx.To)
			assert.Equal(t1, tt.tx.Gas, tx.Gas)
			assert.Equal(t1, tt.tx.GasPrice, tx.GasPrice)
			assert.Equal(t1, tt.tx.Input, tx.Input)
			assert.Equal(t1, tt.tx.Nonce, tx.Nonce)
			assert.Equal(t1, tt.tx.Value, tx.Value)
			assert.Equal(t1, tt.tx.Signature, tx.Signature)
			assert.Equal(t1, tt.tx.ChainID, tx.ChainID)
			assert.Equal(t1, tt.tx.MaxPriorityFeePerGas, tx.MaxPriorityFeePerGas)
			assert.Equal(t1, tt.tx.MaxFeePerGas, tx.MaxFeePerGas)
			for i, accessTuple := range tt.tx.AccessList {
				assert.Equal(t1, accessTuple.Address, tx.AccessList[i].Address)
				assert.Equal(t1, accessTuple.StorageKeys, tx.AccessList[i].StorageKeys)
			}
		})
	}
}

func TestTransaction_SingingHash(t1 *testing.T) {
	tests := []struct {
		tx   *Transaction
		want Hash
	}{
		{
			tx: &Transaction{
				Type:      LegacyTxType,
				From:      HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:        HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:       100000,
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
			},
			want: MustHexToHash("1efbe489013ac8c0dad2202f68ac12657471df8d80f70e0683ec07b0564a32ca"),
		},
		{
			tx: &Transaction{
				Type:      AccessListTxType,
				From:      HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:        HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:       100000,
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
				AccessList: AccessList{
					AccessTuple{
						Address: MustHexToAddress("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHexToHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: MustHexToHash("71cba0039a020b7a524d7746b79bf6d1f8a521eb1a76715d00116ef1c0f56107"),
		},
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:                   HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  100000,
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
				AccessList: AccessList{
					AccessTuple{
						Address: MustHexToAddress("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHexToHash("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHexToHash("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: MustHexToHash("a66ab756479bfd56f29658a8a199319094e84711e8a2de073ec136ef5179c4c9"),
		},
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 HexToAddressPtr("0x1111111111111111111111111111111111111111"),
				To:                   HexToAddressPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  100000,
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            HexToSignature("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
			},
			want: MustHexToHash("c3266152306909bfe339f90fad4f73f958066860300b5a22b98ee6a1d629706c"),
		},
	}
	for n, tt := range tests {
		t1.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			sh, err := tt.tx.SigningHash(hashFunc)
			require.NoError(t1, err)
			require.Equal(t1, tt.want, sh)
		})
	}
}

func hashFunc(data ...[]byte) Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range data {
		h.Write(i)
	}
	return MustBytesToHash(h.Sum(nil))
}
