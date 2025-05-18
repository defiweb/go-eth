package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionOnChain_JSON(t *testing.T) {
	tests := []struct {
		tx       *TransactionOnChain
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "gasPrice": "0x3b9aca00",
				  "input": "0x01020304",
				  "nonce": "0x1",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
				  "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
				  "blockHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
				  "blockNumber": "0x3",
				  "transactionIndex": "0x4"
				}
			`,
			tx: &TransactionOnChain{
				Transaction: &TransactionLegacy{
					TransactionFields: TransactionFields{
						Nonce:     ptr(uint64(1)),
						Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
					},
					CallLegacy: CallLegacy{
						CallFields: CallFields{
							To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
							Value:    big.NewInt(1000000000000000000),
							GasLimit: ptr(uint64(100000)),
							Input:    []byte{1, 2, 3, 4},
						},
						LegacyPriceField: LegacyPriceField{
							GasPrice: big.NewInt(1000000000),
						},
					},
				},
				Hash:             MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				BlockHash:        MustHashFromHexPtr("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
				BlockNumber:      big.NewInt(3),
				TransactionIndex: ptr(uint64(4)),
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			tx := &TransactionOnChain{}

			err := tx.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.tx, tx)

			j, err := tx.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}
