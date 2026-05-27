package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionDynamicFee_JSON(t *testing.T) {
	tests := []struct {
		name     string
		tx       *TransactionDynamicFee
		wantJSON string
	}{
		{
			name:     "empty transaction",
			tx:       &TransactionDynamicFee{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallDynamicFee: CallDynamicFee{
					ExecutionData: ExecutionData{
						From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					DynamicFeeData: DynamicFeeData{
						MaxPriorityFeePerGas: big.NewInt(1000000000),
						MaxFeePerGas:         big.NewInt(2000000000),
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
				  "from": "0x1111111111111111111111111111111111111111",
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "maxFeePerGas": "0x77359400",
				  "maxPriorityFeePerGas": "0x3b9aca00",
				  "input": "0x01020304",
				  "nonce": "0x1",
				  "value": "0xde0b6b3a7640000",
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
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"
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
			tx := NewTransactionDynamicFee()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionDynamicFee_RLP(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionDynamicFee
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionDynamicFee{},
			wantHex: "0x02cc8080808080808080c0808080",
		},
		{
			name: "all fields set",
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallDynamicFee: CallDynamicFee{
					ExecutionData: ExecutionData{
						From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					DynamicFeeData: DynamicFeeData{
						MaxPriorityFeePerGas: big.NewInt(1000000000),
						MaxFeePerGas:         big.NewInt(2000000000),
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
			wantHex: "0x02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RLP
			rlpBytes, err := tt.tx.EncodeRLP()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, hexutil.BytesToHex(rlpBytes))

			// Decode from RLP
			tx := NewTransactionDynamicFee()
			_, err = tx.DecodeRLP(rlpBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionDynamicFee_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionDynamicFee
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionDynamicFee{},
			wantHex: "0x292edeba1be7c90f4dbaed50c44b7f6378633f933202ffe4f547e5a5c2ca3304",
		},
		{
			name: "all fields set",
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallDynamicFee: CallDynamicFee{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					DynamicFeeData: DynamicFeeData{
						MaxPriorityFeePerGas: big.NewInt(1000000000),
						MaxFeePerGas:         big.NewInt(2000000000),
					},
				},
			},
			wantHex: "0xc3266152306909bfe339f90fad4f73f958066860300b5a22b98ee6a1d629706c",
		},
		{
			name: "all fields set with access list",
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallDynamicFee: CallDynamicFee{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					AccessListData: AccessListData{
						AccessList: AccessList{
							AccessTuple{
								Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
								StorageKeys: []Hash{
									MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
									MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
								},
							},
						},
					},
					DynamicFeeData: DynamicFeeData{
						MaxPriorityFeePerGas: big.NewInt(1000000000),
						MaxFeePerGas:         big.NewInt(2000000000),
					},
				},
			},
			wantHex: "0xa66ab756479bfd56f29658a8a199319094e84711e8a2de073ec136ef5179c4c9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sh, err := tt.tx.SigningHash()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, sh.String())
		})
	}
}
