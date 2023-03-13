package crypto

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

func Test_singingHash(t1 *testing.T) {
	tests := []struct {
		tx   *types.Transaction
		want types.Hash
	}{
		// Empty transaction:
		{
			tx:   &types.Transaction{},
			want: types.MustHashFromHex("5460be86ce1e4ca0564b5761c6e7070d9f054b671f5404268335000806423d75"),
		},
		// Legacy transaction:
		{
			tx: (&types.Transaction{}).
				SetType(types.LegacyTxType).
				SetFrom(types.MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(types.MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetGasPrice(big.NewInt(1000000000)).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1),
			want: types.MustHashFromHex("1efbe489013ac8c0dad2202f68ac12657471df8d80f70e0683ec07b0564a32ca"),
		},
		// Access list transaction:
		{
			tx: (&types.Transaction{}).
				SetType(types.AccessListTxType).
				SetFrom(types.MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(types.MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetGasPrice(big.NewInt(1000000000)).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetAccessList(types.AccessList{
					types.AccessTuple{
						Address: types.MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []types.Hash{
							types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							types.MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				}),
			want: types.MustHashFromHex("71cba0039a020b7a524d7746b79bf6d1f8a521eb1a76715d00116ef1c0f56107"),
		},
		// Dynamic fee transaction with access list:
		{
			tx: (&types.Transaction{}).
				SetType(types.DynamicFeeTxType).
				SetFrom(types.MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(types.MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetMaxPriorityFeePerGas(big.NewInt(1000000000)).
				SetMaxFeePerGas(big.NewInt(2000000000)).
				SetAccessList(types.AccessList{
					types.AccessTuple{
						Address: types.MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []types.Hash{
							types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444"),
							types.MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555"),
						},
					},
				}),
			want: types.MustHashFromHex("a66ab756479bfd56f29658a8a199319094e84711e8a2de073ec136ef5179c4c9"),
		},
		// Dynamic fee transaction with no access list:
		{
			tx: (&types.Transaction{}).
				SetType(types.DynamicFeeTxType).
				SetFrom(types.MustAddressFromHex("0x1111111111111111111111111111111111111111")).
				SetTo(types.MustAddressFromHex("0x2222222222222222222222222222222222222222")).
				SetGasLimit(100000).
				SetInput([]byte{1, 2, 3, 4}).
				SetNonce(1).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f")).
				SetChainID(1).
				SetMaxPriorityFeePerGas(big.NewInt(1000000000)).
				SetMaxFeePerGas(big.NewInt(2000000000)),
			want: types.MustHashFromHex("c3266152306909bfe339f90fad4f73f958066860300b5a22b98ee6a1d629706c"),
		},
		// Example from EIP-155:
		{
			tx: (&types.Transaction{}).
				SetType(types.LegacyTxType).
				SetChainID(1).
				SetTo(types.MustAddressFromHex("0x3535353535353535353535353535353535353535")).
				SetGasLimit(21000).
				SetGasPrice(big.NewInt(20000000000)).
				SetNonce(9).
				SetValue(big.NewInt(1000000000000000000)).
				SetSignature(types.SignatureFromVRS(
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
			want: types.MustHashFromHex("daf5a779ae972f972197303d7b574746c7ef83eadac0f2791ad23db92e4c8e53"),
		},
	}
	for n, tt := range tests {
		t1.Run(fmt.Sprintf("case-%d", n+1), func(t1 *testing.T) {
			sh, err := signingHash(tt.tx)
			require.NoError(t1, err)
			require.Equal(t1, tt.want, sh)
		})
	}
}
