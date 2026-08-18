package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionSetCode_JSON(t *testing.T) {
	tests := []struct {
		name     string
		tx       *TransactionSetCode
		wantJSON string
	}{
		{
			name:     "empty transaction",
			tx:       &TransactionSetCode{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			tx: &TransactionSetCode{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallSetCode: CallSetCode{
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
					AuthorizationData: AuthorizationData{
						AuthorizationList: AuthorizationList{
							{
								ChainID: 1,
								Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
								Nonce:   0,
							},
						},
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
				"authorizationList": [
					{
						"chainId": "0x1",
						"address": "0x3333333333333333333333333333333333333333",
						"nonce": "0x0"
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
			tx := NewTransactionSetCode()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionSetCode_RLP(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionSetCode
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionSetCode{},
			wantHex: "0x04cd8080808080808080c0c0808080",
		},
		{
			name: "all fields set",
			tx: &TransactionSetCode{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallSetCode: CallSetCode{
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
					AuthorizationData: AuthorizationData{
						AuthorizationList: AuthorizationList{
							{
								ChainID: 1,
								Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
								Nonce:   0,
							},
						},
					},
				},
			},
			wantHex: "0x04f8ef0101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a05555555555555555555555555555555555555555555555555555555555555555dbda01943333333333333333333333333333333333333333808080806fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RLP
			rlpBytes, err := tt.tx.EncodeRLP()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, hexutil.BytesToHex(rlpBytes))

			// Decode from RLP
			tx := NewTransactionSetCode()
			_, err = tx.DecodeRLP(rlpBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionSetCode_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionSetCode
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionSetCode{},
			wantHex: "0x1aebc3534fcd957f1b31c053815f3f56d7d7dc60b201027942f8d2cd1380cbfb",
		},
		{
			name: "all fields set",
			tx: &TransactionSetCode{
				SigningData: SigningData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallSetCode: CallSetCode{
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
			wantHex: "0xb826798c98525b8c19b2abe284842f53cc5a70e685a7ea8a890165fc61889026",
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

func TestAuthorization_SigningHash(t *testing.T) {
	tests := []struct {
		name    string
		auth    Authorization
		wantHex string
	}{
		{
			name: "chain=1 addr=0x3333 nonce=1",
			auth: Authorization{
				ChainID: 1,
				Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
				Nonce:   1,
			},
			wantHex: "0x689a4ce94c04490ed6d26c33123f595f4203153a93f0ed9e72186c69194e8425",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := tt.auth.SigningHash()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, h.String())
		})
	}
}
