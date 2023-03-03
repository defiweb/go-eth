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
		// Empty transaction:
		{
			tx: &Transaction{
				Gas:      func() *uint64 { v := uint64(0); return &v }(),
				GasPrice: func() *big.Int { v := big.NewInt(0); return v }(),
				Nonce:    func() *big.Int { v := big.NewInt(0); return v }(),
				Value:    func() *big.Int { v := big.NewInt(0); return v }(),
			},
			want: hexutil.MustHexToBytes("c9808080808080808080"),
		},
		// Legacy transaction:
		{
			tx: &Transaction{
				Type:      LegacyTxType,
				From:      AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:        AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:       func() *uint64 { v := uint64(100000); return &v }(),
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
			},
			want: hexutil.MustHexToBytes("f87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Access list transaction:
		{
			tx: &Transaction{
				Type:      AccessListTxType,
				From:      AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:        AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:       func() *uint64 { v := uint64(100000); return &v }(),
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
				AccessList: AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: hexutil.MustHexToBytes("01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Dynamic fee transaction:
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:                   AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  func() *uint64 { v := uint64(100000); return &v }(),
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
				AccessList: AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: hexutil.MustHexToBytes("02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Dynamic fee transaction with no access list:
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:                   AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  func() *uint64 { v := uint64(100000); return &v }(),
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
			},
			want: hexutil.MustHexToBytes("02f8770101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304c06fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"),
		},
		// Example from EIP-155:
		{
			tx: &Transaction{
				Type:     LegacyTxType,
				ChainID:  big.NewInt(1),
				To:       AddressFromHexPtr("0x3535353535353535353535353535353535353535"),
				Gas:      func() *uint64 { v := uint64(21000); return &v }(),
				GasPrice: func() *big.Int { v, _ := new(big.Int).SetString("20000000000", 10); return v }(),
				Nonce:    func() *big.Int { v := big.NewInt(9); return v }(),
				Value:    func() *big.Int { v, _ := new(big.Int).SetString("1000000000000000000", 10); return v }(),
				Signature: MustSignatureFromBigIntPtr(
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
				),
			},
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

func TestTransaction_SingingHash(t1 *testing.T) {
	tests := []struct {
		tx   *Transaction
		want Hash
	}{
		// Empty transaction:
		{
			tx:   &Transaction{},
			want: MustHashFromHex("5460be86ce1e4ca0564b5761c6e7070d9f054b671f5404268335000806423d75"),
		},
		// Legacy transaction:
		{
			tx: &Transaction{
				Type:      LegacyTxType,
				From:      AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:        AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:       func() *uint64 { v := uint64(100000); return &v }(),
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
			},
			want: MustHashFromHex("1efbe489013ac8c0dad2202f68ac12657471df8d80f70e0683ec07b0564a32ca"),
		},
		// Access list transaction:
		{
			tx: &Transaction{
				Type:      AccessListTxType,
				From:      AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:        AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:       func() *uint64 { v := uint64(100000); return &v }(),
				GasPrice:  new(big.Int).SetUint64(1000000000),
				Input:     []byte{1, 2, 3, 4},
				Nonce:     new(big.Int).SetUint64(1),
				Value:     new(big.Int).SetUint64(1000000000000000000),
				Signature: SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:   new(big.Int).SetUint64(1),
				AccessList: AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: MustHashFromHex("71cba0039a020b7a524d7746b79bf6d1f8a521eb1a76715d00116ef1c0f56107"),
		},
		// Dynamic fee transaction with access list:
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:                   AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  func() *uint64 { v := uint64(100000); return &v }(),
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
				AccessList: AccessList{
					AccessTuple{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				},
			},
			want: MustHashFromHex("a66ab756479bfd56f29658a8a199319094e84711e8a2de073ec136ef5179c4c9"),
		},
		// Dynamic fee transaction with no access list:
		{
			tx: &Transaction{
				Type:                 DynamicFeeTxType,
				From:                 AddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:                   AddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				Gas:                  func() *uint64 { v := uint64(100000); return &v }(),
				Input:                []byte{1, 2, 3, 4},
				Nonce:                new(big.Int).SetUint64(1),
				Value:                new(big.Int).SetUint64(1000000000000000000),
				Signature:            SignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				ChainID:              new(big.Int).SetUint64(1),
				MaxPriorityFeePerGas: new(big.Int).SetUint64(1000000000),
				MaxFeePerGas:         new(big.Int).SetUint64(2000000000),
			},
			want: MustHashFromHex("c3266152306909bfe339f90fad4f73f958066860300b5a22b98ee6a1d629706c"),
		},
		// Example from EIP-155:
		{
			tx: &Transaction{
				Type:     LegacyTxType,
				ChainID:  big.NewInt(1),
				To:       AddressFromHexPtr("0x3535353535353535353535353535353535353535"),
				Gas:      func() *uint64 { v := uint64(21000); return &v }(),
				GasPrice: func() *big.Int { v, _ := new(big.Int).SetString("20000000000", 10); return v }(),
				Nonce:    func() *big.Int { v := big.NewInt(9); return v }(),
				Value:    func() *big.Int { v, _ := new(big.Int).SetString("1000000000000000000", 10); return v }(),
				Signature: MustSignatureFromBigIntPtr(
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
				),
			},
			want: MustHashFromHex("daf5a779ae972f972197303d7b574746c7ef83eadac0f2791ad23db92e4c8e53"),
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

func equalTx(t *testing.T, expected, got *Transaction) {
	assert.Equal(t, expected.Type, got.Type)
	assert.Equal(t, expected.To, got.To)
	assert.Equal(t, expected.Gas, got.Gas)
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

func hashFunc(data ...[]byte) Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range data {
		h.Write(i)
	}
	return MustHashFromBytes(h.Sum(nil))
}
