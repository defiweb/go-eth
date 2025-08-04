package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionAccessList_JSON(t *testing.T) {
	tests := []struct {
		name     string
		tx       *TransactionAccessList
		wantJSON string
	}{
		{
			name:     "all fields nil",
			tx:       &TransactionAccessList{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			tx: &TransactionAccessList{
				TransactionData: TransactionData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallAccessList: CallAccessList{
					CallData: CallData{
						From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(1000000000),
					},
					AccessListData: AccessListData{
						AccessList: []AccessTuple{{
							Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
							StorageKeys: []Hash{
								MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
								MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
							},
						}},
					},
				},
			},
			wantJSON: `{
                "chainId": "0x1",
                "nonce": "0x1",
                "gasPrice": "0x3b9aca00",
                "gas": "0x186a0",
                "to": "0x2222222222222222222222222222222222222222",
                "value": "0xde0b6b3a7640000",
                "input": "0x01020304",
                "accessList": [
                    {
                        "address": "0x3333333333333333333333333333333333333333",
                        "storageKeys": [
                            "0x4444444444444444444444444444444444444444444444444444444444444444",
                            "0x5555555555555555555555555555555555555555555555555555555555555555"
                        ]
                    }
                ],
                "v": "0x6f",
                "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
                "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
                "from": "0x1111111111111111111111111111111111111111"
            }`,
		},
		{
			name: "invalid negative nonce",
			tx: &TransactionAccessList{
				TransactionData: TransactionData{
					Nonce: ptr(uint64(18446744073709551615)), // Max uint64 value to simulate negative when interpreted incorrectly
				},
			},
			wantJSON: `{
                "nonce": "0xffffffffffffffff"
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to JSON
			jsonBytes, err := tt.tx.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Decode from JSON
			tx := NewTransactionAccessList()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionAccessList_RLP(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionAccessList
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionAccessList{},
			wantHex: "0x01cb80808080808080c0808080",
		},
		{
			name: "all fields set",
			tx: &TransactionAccessList{
				TransactionData: TransactionData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallAccessList: CallAccessList{
					CallData: CallData{
						From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(1000000000),
					},
					AccessListData: AccessListData{
						AccessList: []AccessTuple{{
							Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
							StorageKeys: []Hash{
								MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
								MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
							},
						}},
					},
				},
			},
			wantHex: "0x01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RLP
			rlpBytes, err := tt.tx.EncodeRLP()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, hexutil.BytesToHex(rlpBytes))

			// Decode from RLP
			tx := NewTransactionAccessList()
			_, err = tx.DecodeRLP(rlpBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionAccessList_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionAccessList
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionAccessList{},
			wantHex: "0xc0157440e7609b2ddee74686831421f05b238ed4c981363e64df8eb1c1ea6afc",
		},
		{
			name: "all fields set",
			tx: &TransactionAccessList{
				TransactionData: TransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallAccessList: CallAccessList{
					CallData: CallData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(1000000000),
					},
				},
			},
			wantHex: "0x46ba790cdf341de06f08944eecd84721e9ae3c4324098f882597d9817eeba63b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh, err := tt.tx.CalculateSigningHash()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, sh.String())
		})
	}
}
