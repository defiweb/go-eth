package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionDecoder_DecodeRLP(t *testing.T) {
	tests := []struct {
		tx  Transaction
		rlp string
	}{
		{
			tx:  &TransactionLegacy{},
			rlp: "0xc9808080808080808080",
		},
		{
			tx: &TransactionLegacy{
				SigningData: SigningData{
					ChainID:   ptr(uint64(38)),
					Nonce:     ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallLegacy: CallLegacy{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyFeeData: LegacyFeeData{
						GasPrice: big.NewInt(1000000000),
					},
				},
			},
			rlp: "0xf87001843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a764000084010203046fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
		{
			tx: &TransactionLegacy{
				SigningData: SigningData{
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
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(21000)),
					},
					LegacyFeeData: LegacyFeeData{
						GasPrice: big.NewInt(20000000000),
					},
				},
			},
			rlp: "0xf86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83",
		},
		{
			tx:  &TransactionAccessList{},
			rlp: "0x01cb80808080808080c0808080",
		},
		{
			tx: &TransactionAccessList{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallAccessList: CallAccessList{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyFeeData: LegacyFeeData{
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
			rlp: "0x01f8ce0101843b9aca00830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
		{
			tx:  &TransactionDynamicFee{},
			rlp: "0x02cc8080808080808080c0808080",
		},
		{
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
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
			rlp: "0x02f8d30101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a055555555555555555555555555555555555555555555555555555555555555556fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
		{
			tx:  &TransactionBlob{},
			rlp: "0x03ce8080808080808080c080c0808080",
		},
		{
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					AccessListData: AccessListData{
						AccessList: []AccessTuple{{
							Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
							StorageKeys: []Hash{
								MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
								MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
							},
						}},
					},
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							{
								Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
							},
						},
					},
				},
			},
			rlp: "0x03f8fa0101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a0555555555555555555555555555555555555555555555555555555555555555584b2d05e00e1a066666666666666666666666666666666666666666666666666666666666666666fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			tx, err := DefaultTransactionDecoder.DecodeRLP(hexutil.MustHexToBytes(tt.rlp))
			require.NoError(t, err)
			assertEqualTX(t, tt.tx, tx)
		})
	}
}

func TestTransactionDeocder_DecodeJSON(t *testing.T) {
	tests := []struct {
		tx   Transaction
		json string
	}{
		{
			tx:   &TransactionLegacy{},
			json: `{}`,
		},
		{
			tx: &TransactionLegacy{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallLegacy: CallLegacy{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyFeeData: LegacyFeeData{
						GasPrice: big.NewInt(1000000000),
					},
				},
			},
			json: `
				{
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "gasPrice": "0x3b9aca00",
				  "input": "0x01020304",
				  "nonce": "0x1",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"
				}
			`,
		},
		{
			tx: &TransactionLegacy{
				SigningData: SigningData{
					Nonce: ptr(uint64(9)),
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
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x3535353535353535353535353535353535353535"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(21000)),
					},
					LegacyFeeData: LegacyFeeData{
						GasPrice: big.NewInt(20000000000),
					},
				},
			},
			json: `
				{
				  "to": "0x3535353535353535353535353535353535353535",
				  "gas": "0x5208",
				  "gasPrice": "0x4a817c800",
				  "nonce": "0x9",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x25",
				  "r": "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276",
				  "s": "0x67cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"
				}
			`,
		},
		{
			tx: &TransactionAccessList{},
			json: `
				{
				  "type": "0x1"
				}
			`,
		},
		{
			tx: &TransactionAccessList{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallAccessList: CallAccessList{
					ExecutionData: ExecutionData{
						To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
						Value:    big.NewInt(1000000000000000000),
						GasLimit: ptr(uint64(100000)),
						Input:    []byte{1, 2, 3, 4},
					},
					LegacyFeeData: LegacyFeeData{
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
			json: `
				{
				  "chainId": "0x1",
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "gasPrice": "0x3b9aca00",
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
				}
			`,
		},
		{
			tx: &TransactionDynamicFee{},
			json: `
				{
				  "type": "0x2"
				}
			`,
		},
		{
			tx: &TransactionDynamicFee{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
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
			json: `
				{
				  "chainId": "0x1",
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "maxFeePerGas": "0x77359400",
				  "maxPriorityFeePerGas": "0x3b9aca00",
				  "input": "0x01020304",
				  "Nonce": "0x1",
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
				}
			`,
		},
		{
			tx: &TransactionBlob{},
			json: `
				{
				  "type": "0x3"
				}
			`,
		},
		{
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					AccessListData: AccessListData{
						AccessList: []AccessTuple{{
							Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
							StorageKeys: []Hash{
								MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
								MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
							},
						}},
					},
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							{
								Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
							},
						},
					},
				},
			},
			json: `
				{
				  "chainId": "0x1",
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "maxFeePerGas": "0x77359400",
				  "maxFeePerBlobGas": "0xb2d05e00",
				  "maxPriorityFeePerGas": "0x3b9aca00",
				  "input": "0x01020304",
				  "Nonce": "0x1",
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
				  "blobVersionedHashes": [
					"0x6666666666666666666666666666666666666666666666666666666666666666"
				  ],
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"
				}
			`,
		},
		{
			tx: &TransactionLegacy{},
			json: `
				{
				  "accessList": [],
				  "maxFeePerGas": "0x0",
				  "blobVersionedHashes": [
					"0x6666666666666666666666666666666666666666666666666666666666666666"
				  ],
				  "type": "0x0"
				}
			`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			tx, err := DefaultTransactionDecoder.DecodeJSON([]byte(tt.json))
			require.NoError(t, err)
			assertEqualTX(t, tt.tx, tx)
		})
	}
}
