package types

import (
	"math/big"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/kzg4844"
	"github.com/defiweb/go-eth/hexutil"
)

func TestTransactionBlob_JSON(t *testing.T) {
	remZerosRx := regexp.MustCompile(`0{128,}`)
	tests := []struct {
		name     string
		tx       *TransactionBlob
		wantJSON string
	}{
		{
			name:     "empty transaction",
			tx:       &TransactionBlob{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							{Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone)},
							{Hash: MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", PadNone)},
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
				  "maxFeePerBlobGas": "0xb2d05e00",
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
				  "blobVersionedHashes": [
					"0x6666666666666666666666666666666666666666666666666666666666666666",
					"0x7777777777777777777777777777777777777777777777777777777777777777"
				  ],
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84"
				}`,
		},
		{
			name: "blobs with shortened zero fields",
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							newBlob("blob1"),
							newBlob("blob2"),
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
				  "maxFeePerBlobGas": "0xb2d05e00",
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
				  "blobVersionedHashes": [
					"0x01e951827dab35ecb4ce6e29ca6779ad0ac958f06ac1d54eb6e7523f3e3febeb",
					"0x01c6f423eec4f5e9ebd79e6fb5eb9bce57dd102f2f0a6fdee0bf58c4e109e27a"
				  ],
				  "blobs": [
					"0x626c6f6231",
					 "0x626c6f6232"
				  ],
				  "commitments": [
					"0x832726ece34fb93100194291b75b7f5fa920d5b896e5edafeed46f3636fadc485f493490bd596c76171516024bcf7a00",
					"0x896a515deb6c1ac23436f5de75186316b53851be3cb4225437e58c34e376bbd6dea13d7a39b8d5dd9f6cba1c052ab0cb"
				  ],
				  "proofs": [
					"0xacbe7bde870d1e7c239063368dbf62ce3bfaef1c08006e8e724da2837294e19d991ba5677a628e64b6c6906ff1786b3c",
					"0x8fd7fb4172048b9d8deccb939eeb787c45d0bb00d5f238e13084bd3fbe8050215ce488999f86fccddc36d1a257a19513"
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
			assert.JSONEq(t, tt.wantJSON, string(remZerosRx.ReplaceAll(jsonBytes, []byte(""))))

			// Decode from JSON
			tx := NewTransactionBlob()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			tx.ChainID = tt.tx.ChainID
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionBlob_RLP(t *testing.T) {
	tests := []struct {
		name     string
		tx       *TransactionBlob
		wantHex  string
		wantHash bool
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionBlob{},
			wantHex: "0x03ce8080808080808080c080c0808080",
		},
		{
			name: "all fields set",
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							{Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone)},
							{Hash: MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", PadNone)},
						},
					},
				},
			},
			wantHex: "0x03f9011c0101843b9aca008477359400830186a0942222222222222222222222222222222222222222880de0b6b3a76400008401020304f85bf859943333333333333333333333333333333333333333f842a04444444444444444444444444444444444444444444444444444444444444444a0555555555555555555555555555555555555555555555555555555555555555584b2d05e00f842a06666666666666666666666666666666666666666666666666666666666666666a077777777777777777777777777777777777777777777777777777777777777776fa0a3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490a08051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		},
		{
			name: "hash output",
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							newBlob("blob1"),
							newBlob("blob2"),
						},
					},
				},
			},
			wantHex:  "0x848eb4e644a60e42df3b639eb40c0f3763d13ebc2a33aa06e9b2acc22c51f59e",
			wantHash: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to RLP
			rlpBytes, err := tt.tx.EncodeRLP()
			require.NoError(t, err)

			if tt.wantHash {
				hash := crypto.Keccak256(rlpBytes)
				assert.Equal(t, tt.wantHex, hexutil.BytesToHex(hash[:]))
			} else {
				assert.Equal(t, tt.wantHex, hexutil.BytesToHex(rlpBytes))
			}

			// Decode from RLP
			tx := NewTransactionBlob()
			_, err = tx.DecodeRLP(rlpBytes)
			require.NoError(t, err)

			// Compare the original and decoded transactions
			tx.From = tt.tx.From
			assertEqualTX(t, tx, tt.tx)
		})
	}
}

func TestTransactionBlob_CalculateSigningHash(t *testing.T) {
	tests := []struct {
		name    string
		tx      *TransactionBlob
		wantHex string
	}{
		{
			name:    "empty transaction",
			tx:      &TransactionBlob{},
			wantHex: "0x846c9b47f161837f5068b0ffb0c1a98785302f89d613338ccfa9a1c72c9f951d",
		},
		{
			name: "all fields set",
			tx: &TransactionBlob{
				SigningData: SigningData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
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
				},
			},
			wantHex: "0x0604b49731147cf745c666f1a67bf1b5e9fbee127085b3d4c4958191590e8bce",
		},
		{
			name: "all fields set with access list",
			tx: &TransactionBlob{
				SigningData: SigningData{
					ChainID: ptr(uint64(1)),
					Nonce:   ptr(uint64(1)),
				},
				CallBlob: CallBlob{
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
			wantHex: "0x3faa63efab3e460606c31cd9a2e8791d87e91954137571eddb3b4b0abc69e2cd",
		},
		{
			name: "with blobs and access list",
			tx: &TransactionBlob{
				SigningData: SigningData{
					Nonce:     ptr(uint64(1)),
					ChainID:   ptr(uint64(1)),
					Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
				},
				CallBlob: CallBlob{
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
					BlobData: BlobData{
						MaxFeePerBlobGas: big.NewInt(3000000000),
						Blobs: []BlobInfo{
							newBlob("blob1"),
							newBlob("blob2"),
						},
					},
				},
			},
			wantHex: "0x09f9204d83af238e1c0044bf22b4dd52ea5c25390b27bdbd38024212e238934d",
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

func newBlob(data string) BlobInfo {
	d := new(kzg4844.Blob)
	copy(d[:], data)
	b, err := NewBlobInfo(d)
	if err != nil {
		panic(err)
	}
	return b
}
