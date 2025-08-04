package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionLegacy_JSON(t *testing.T) {
	tests := []struct {
		name     string
		tx       *TransactionLegacy
		wantJSON string
	}{
		{
			name:     "empty transaction",
			tx:       &TransactionLegacy{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					Nonce:     ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallLegacy: CallLegacy{
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
				},
			},
			wantJSON: `{
				  "from": "0x1111111111111111111111111111111111111111",
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "gasPrice": "0x3b9aca00",
				  "input": "0x01020304",
				  "nonce": "0x1",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"
				}`,
		},
		{
			name: "example from EIP-155",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(9)),
					Signature: SignatureFromVRSPtr(
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
				CallLegacy: CallLegacy{
					CallData: CallData{
						To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(21000)),
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(20000000000),
					},
				},
			},
			wantJSON: `{
				  "chainId": "0x1",
				  "to": "0x3535353535353535353535353535353535353535",
				  "gas": "0x5208",
				  "gasPrice": "0x4a817c800",
				  "nonce": "0x9",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x25",
				  "r": "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276",
				  "s": "0x67cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"
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
			tx := NewTransactionLegacy()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionLegacy_RLP(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionLegacy
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionLegacy{},
			wantHex: "0xc9808080808080808080",
		},
		{
			name: "all fields set",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					Nonce:     ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallLegacy: CallLegacy{
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
				},
			},
			wantHex: "0xf87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
		{
			name: "example from EIP-155",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(9)),
					Signature: SignatureFromVRSPtr(
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
				CallLegacy: CallLegacy{
					CallData: CallData{
						To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(21000)),
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(20000000000),
					},
				},
			},
			wantHex: "0xf86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RLP
			rlpBytes, err := tt.tx.EncodeRLP()
			require.NoError(t, err)
			assert.Equal(t, tt.wantHex, hexutil.BytesToHex(rlpBytes))

			// Decode from RLP
			tx := NewTransactionLegacy()
			_, err = tx.DecodeRLP(rlpBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionLegacy_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionLegacy
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionLegacy{},
			wantHex: "0x5460be86ce1e4ca0564b5761c6e7070d9f054b671f5404268335000806423d75",
		},
		{
			name: "all fields set",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallLegacy: CallLegacy{
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
				},
			},
			wantHex: "0x1efbe489013ac8c0dad2202f68ac12657471df8d80f70e0683ec07b0564a32ca",
		},
		{
			name: "example from EIP-155",
			tx: &TransactionLegacy{
				TransactionData: TransactionData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(9)),
				},
				CallLegacy: CallLegacy{
					CallData: CallData{
						To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(21000)),
					},
					LegacyPriceData: LegacyPriceData{
						GasPrice: big.NewInt(20000000000),
					},
				},
			},
			wantHex: "0xdaf5a779ae972f972197303d7b574746c7ef83eadac0f2791ad23db92e4c8e53",
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
