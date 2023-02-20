package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransaction_Raw(t1 *testing.T) {
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
				ChainID:   new(big.Int).SetUint64(1),
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
			want: hexutil.MustHexToBytes("01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a05555555555555555555555555555555555555555555555555555555555555555a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
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
			want: hexutil.MustHexToBytes("02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a05555555555555555555555555555555555555555555555555555555555555555a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
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
			rlp, err := tt.tx.Raw()
			require.NoError(t1, err)
			require.Equal(t1, tt.want, rlp)
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
			want: MustHexToHash("0f0ef55dca9a5f856088348ced1f393078ccbf2ddcda78b7146dd8280a824b4a"),
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
			want: MustHexToHash("5fb17026de309d639bdc2cd78050bc1629aeda743ad2b9b51e131e19f73df6b9"),
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
