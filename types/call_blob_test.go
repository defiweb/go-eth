package types

import (
	"math/big"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallBlob_JSON(t *testing.T) {
	remZerosRx := regexp.MustCompile(`0{128,}`)
	tests := []struct {
		name     string
		call     *CallBlob
		wantJSON string
	}{
		{
			name:     "empty call",
			call:     &CallBlob{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			call: &CallBlob{
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
			wantJSON: `{
			  "from": "0x1111111111111111111111111111111111111111",
			  "to": "0x2222222222222222222222222222222222222222",
			  "gas": "0x186a0",
			  "maxFeePerGas": "0x77359400",
			  "maxFeePerBlobGas": "0xb2d05e00",
			  "maxPriorityFeePerGas": "0x3b9aca00",
			  "input": "0x01020304",
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
			  ]
			}`,
		},
		{
			name: "blobs with shortened zero fields",
			call: &CallBlob{
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
			wantJSON: `{
			  "from": "0x1111111111111111111111111111111111111111",
			  "to": "0x2222222222222222222222222222222222222222",
			  "gas": "0x186a0",
			  "maxFeePerGas": "0x77359400",
			  "maxFeePerBlobGas": "0xb2d05e00",
			  "maxPriorityFeePerGas": "0x3b9aca00",
			  "input": "0x01020304",
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
			  ]
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode to JSON
			jsonBytes, err := tt.call.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(remZerosRx.ReplaceAll(jsonBytes, []byte(""))))

			// Decode from JSON
			tx := NewCallBlob()
			err = tx.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded calls
			tx.From = tt.call.From
			assertEqualCall(t, tx, tt.call)
		})
	}
}
