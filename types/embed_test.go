package types

import (
	"encoding/json"
	"math/big"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/crypto/kzg4844"
)

func TestTransactionData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		data     *TransactionData
		want     *jsonTransaction
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &TransactionData{},
			want:     &jsonTransaction{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			data: &TransactionData{
				ChainID: ptr(uint64(1)),
				Nonce:   ptr(uint64(2)),
				Signature: &Signature{
					V: big.NewInt(27),
					R: big.NewInt(123456),
					S: big.NewInt(654321),
				},
			},
			want: &jsonTransaction{
				ChainID: NumberFromUint64Ptr(1),
				Nonce:   NumberFromUint64Ptr(2),
				V:       NumberFromBigIntPtr(big.NewInt(27)),
				R:       NumberFromBigIntPtr(big.NewInt(123456)),
				S:       NumberFromBigIntPtr(big.NewInt(654321)),
			},
			wantJSON: `{
                "chainId": "0x1",
                "nonce": "0x2",
                "v": "0x1b",
                "r": "0x1e240",
                "s": "0x9fbf1"
            }`,
		},
		{
			name: "nil signature",
			data: &TransactionData{
				ChainID:   ptr(uint64(1)),
				Nonce:     ptr(uint64(2)),
				Signature: nil,
			},
			want: &jsonTransaction{
				ChainID: NumberFromUint64Ptr(1),
				Nonce:   NumberFromUint64Ptr(2),
				V:       nil,
				R:       nil,
				S:       nil,
			},
			wantJSON: `{
                "chainId": "0x1",
                "nonce": "0x2"
            }`,
		},
		{
			name: "max uint64 values",
			data: &TransactionData{
				ChainID: ptr(^uint64(0)),
				Nonce:   ptr(^uint64(0)),
			},
			want: &jsonTransaction{
				ChainID: NumberFromUint64Ptr(^uint64(0)),
				Nonce:   NumberFromUint64Ptr(^uint64(0)),
			},
			wantJSON: `{
                "chainId": "0xffffffffffffffff",
                "nonce": "0xffffffffffffffff"
            }`,
		},
		{
			name: "zero values",
			data: &TransactionData{
				ChainID: ptr(uint64(0)),
				Nonce:   ptr(uint64(0)),
				Signature: &Signature{
					V: big.NewInt(0),
					R: big.NewInt(0),
					S: big.NewInt(0),
				},
			},
			want: &jsonTransaction{
				ChainID: NumberFromUint64Ptr(0),
				Nonce:   NumberFromUint64Ptr(0),
				V:       NumberFromBigIntPtr(big.NewInt(0)),
				R:       NumberFromBigIntPtr(big.NewInt(0)),
				S:       NumberFromBigIntPtr(big.NewInt(0)),
			},
			wantJSON: `{
                "chainId": "0x0",
                "nonce": "0x0",
                "v": "0x0",
                "r": "0x0",
                "s": "0x0"
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonTransaction
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.ChainID, j.ChainID)
			assert.Equal(t, tt.want.Nonce, j.Nonce)
			assert.Equal(t, tt.want.V, j.V)
			assert.Equal(t, tt.want.R, j.R)
			assert.Equal(t, tt.want.S, j.S)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Test fromJSON
			var data TransactionData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.ChainID, data.ChainID)
			assert.Equal(t, tt.data.Nonce, data.Nonce)
			if tt.data.Signature == nil {
				assert.Nil(t, data.Signature)
			} else {
				assert.NotNil(t, data.Signature)
				assert.Equal(t, tt.data.Signature.V, data.Signature.V)
				assert.Equal(t, tt.data.Signature.R, data.Signature.R)
				assert.Equal(t, tt.data.Signature.S, data.Signature.S)
			}
		})
	}
}

func TestCallData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		data     *CallData
		want     *jsonCall
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &CallData{},
			want:     &jsonCall{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			data: &CallData{
				From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				GasLimit: ptr(uint64(21000)),
				Value:    big.NewInt(1000000000000000000),
				Input:    []byte{0x01, 0x02, 0x03},
			},
			want: &jsonCall{
				From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				GasLimit: NumberFromUint64Ptr(21000),
				Value:    NumberFromBigIntPtr(big.NewInt(1000000000000000000)),
				Input:    []byte{0x01, 0x02, 0x03},
			},
			wantJSON: `{
                "from": "0x1111111111111111111111111111111111111111",
                "to": "0x2222222222222222222222222222222222222222",
                "gas": "0x5208",
                "value": "0xde0b6b3a7640000",
                "input": "0x010203"
            }`,
		},
		{
			name: "max values",
			data: &CallData{
				From:     MustAddressFromHexPtr("0xffffffffffffffffffffffffffffffffffffffff"),
				To:       MustAddressFromHexPtr("0xffffffffffffffffffffffffffffffffffffffff"),
				GasLimit: ptr(^uint64(0)),
				Value:    big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1)),
				Input:    []byte{0x01, 0x02, 0x03},
			},
			want: &jsonCall{
				From:     MustAddressFromHexPtr("0xffffffffffffffffffffffffffffffffffffffff"),
				To:       MustAddressFromHexPtr("0xffffffffffffffffffffffffffffffffffffffff"),
				GasLimit: NumberFromUint64Ptr(^uint64(0)),
				Value:    NumberFromBigIntPtr(big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1))),
				Input:    []byte{0x01, 0x02, 0x03},
			},
			wantJSON: `{
				"from": "0xffffffffffffffffffffffffffffffffffffffff",
				"to": "0xffffffffffffffffffffffffffffffffffffffff",
				"gas": "0xffffffffffffffff",
				"value": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"input": "0x010203"
			}`,
		},
		{
			name: "zero values",
			data: &CallData{
				From:     MustAddressFromHexPtr("0x0000000000000000000000000000000000000000"),
				To:       MustAddressFromHexPtr("0x0000000000000000000000000000000000000000"),
				GasLimit: ptr(uint64(0)),
				Value:    big.NewInt(0),
				Input:    []byte{},
			},
			want: &jsonCall{
				From:     MustAddressFromHexPtr("0x0000000000000000000000000000000000000000"),
				To:       MustAddressFromHexPtr("0x0000000000000000000000000000000000000000"),
				GasLimit: NumberFromUint64Ptr(0),
				Value:    NumberFromBigIntPtr(big.NewInt(0)),
				Input:    []byte{},
			},
			wantJSON: `{
				"from": "0x0000000000000000000000000000000000000000",
				"to": "0x0000000000000000000000000000000000000000",
				"gas": "0x0",
				"value": "0x0"
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonCall
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.From, j.From)
			assert.Equal(t, tt.want.To, j.To)
			assert.Equal(t, tt.want.GasLimit, j.GasLimit)
			assert.Equal(t, tt.want.Value, j.Value)
			assert.Equal(t, tt.want.Input, j.Input)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Test fromJSON
			var data CallData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.From, data.From)
			assert.Equal(t, tt.data.To, data.To)
			assert.Equal(t, tt.data.GasLimit, data.GasLimit)
			assert.Equal(t, tt.data.Value, data.Value)
			assert.Equal(t, tt.data.Input, data.Input)
		})
	}
}

func TestLegacyPriceData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		data     *LegacyPriceData
		want     *jsonCall
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &LegacyPriceData{},
			want:     &jsonCall{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			data: &LegacyPriceData{
				GasPrice: big.NewInt(2000000000),
			},
			want: &jsonCall{
				GasPrice: NumberFromBigIntPtr(big.NewInt(2000000000)),
			},
			wantJSON: `{
                "gasPrice": "0x77359400"
            }`,
		},
		{
			name: "zero value",
			data: &LegacyPriceData{
				GasPrice: big.NewInt(0),
			},
			want: &jsonCall{
				GasPrice: NumberFromBigIntPtr(big.NewInt(0)),
			},
			wantJSON: `{
                "gasPrice": "0x0"
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonCall
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.GasPrice, j.GasPrice)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Test fromJSON
			var data LegacyPriceData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.GasPrice, data.GasPrice)
		})
	}
}

func TestAccessListData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		data     *AccessListData
		want     *jsonCall
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &AccessListData{},
			want:     &jsonCall{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			data: &AccessListData{
				AccessList: AccessList{
					{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
						},
					},
				},
			},
			want: &jsonCall{
				AccessList: AccessList{
					{
						Address: MustAddressFromHex("0x3333333333333333333333333333333333333333"),
						StorageKeys: []Hash{
							MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
							MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
						},
					},
				},
			},
			wantJSON: `{
                "accessList": [
                    {
                        "address": "0x3333333333333333333333333333333333333333",
                        "storageKeys": [
                            "0x4444444444444444444444444444444444444444444444444444444444444444",
                            "0x5555555555555555555555555555555555555555555555555555555555555555"
                        ]
                    }
                ]
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonCall
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.AccessList, j.AccessList)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Test fromJSON
			var data AccessListData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.AccessList, data.AccessList)
		})
	}
}

func TestDynamicFeeData_JSON(t *testing.T) {
	tests := []struct {
		name     string
		data     *DynamicFeeData
		want     *jsonCall
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &DynamicFeeData{},
			want:     &jsonCall{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			data: &DynamicFeeData{
				MaxFeePerGas:         big.NewInt(2000000000),
				MaxPriorityFeePerGas: big.NewInt(1000000000),
			},
			want: &jsonCall{
				MaxFeePerGas:         NumberFromBigIntPtr(big.NewInt(2000000000)),
				MaxPriorityFeePerGas: NumberFromBigIntPtr(big.NewInt(1000000000)),
			},
			wantJSON: `{
                "maxFeePerGas": "0x77359400",
                "maxPriorityFeePerGas": "0x3b9aca00"
            }`,
		},
		{
			name: "zero values",
			data: &DynamicFeeData{
				MaxFeePerGas:         big.NewInt(0),
				MaxPriorityFeePerGas: big.NewInt(0),
			},
			want: &jsonCall{
				MaxFeePerGas:         NumberFromBigIntPtr(big.NewInt(0)),
				MaxPriorityFeePerGas: NumberFromBigIntPtr(big.NewInt(0)),
			},
			wantJSON: `{
                "maxFeePerGas": "0x0",
                "maxPriorityFeePerGas": "0x0"
            }`,
		},
		{
			name: "max values",
			data: &DynamicFeeData{
				MaxFeePerGas:         big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1)),
				MaxPriorityFeePerGas: big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1)),
			},
			want: &jsonCall{
				MaxFeePerGas:         NumberFromBigIntPtr(big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1))),
				MaxPriorityFeePerGas: NumberFromBigIntPtr(big.NewInt(0).Sub(big.NewInt(0).SetBit(big.NewInt(0), 256, 1), big.NewInt(1))),
			},
			wantJSON: `{
				"maxFeePerGas": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"maxPriorityFeePerGas": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
			}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonCall
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.MaxFeePerGas, j.MaxFeePerGas)
			assert.Equal(t, tt.want.MaxPriorityFeePerGas, j.MaxPriorityFeePerGas)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Test fromJSON
			var data DynamicFeeData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.MaxFeePerGas, data.MaxFeePerGas)
			assert.Equal(t, tt.data.MaxPriorityFeePerGas, data.MaxPriorityFeePerGas)
		})
	}
}

func TestBlobData_JSON(t *testing.T) {
	remZerosRx := regexp.MustCompile(`0{128,}`)
	tests := []struct {
		name     string
		data     *BlobData
		want     *jsonCall
		wantJSON string
	}{
		{
			name:     "all fields nil",
			data:     &BlobData{},
			want:     &jsonCall{},
			wantJSON: `{}`,
		},
		{
			name: "blobs with sidecars",
			data: &BlobData{
				MaxFeePerBlobGas: big.NewInt(3000000000),
				Blobs: []BlobInfo{
					{
						Hash: MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
						Sidecar: &BlobSidecar{
							Blob:       kzg4844.Blob{0x01, 0x02, 0x03},
							Commitment: kzg4844.Commitment{0x04, 0x05, 0x06},
							Proof:      kzg4844.Proof{0x07, 0x08, 0x09},
						},
					},
				},
			},
			want: &jsonCall{
				MaxFeePerBlobGas: NumberFromBigIntPtr(big.NewInt(3000000000)),
				BlobHashes: []Hash{
					MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
				},
				Blobs: []kzgBlob{
					{0x01, 0x02, 0x03},
				},
				Commitments: []kzgCommitment{
					{0x04, 0x05, 0x06},
				},
				Proofs: []kzgProof{
					{0x07, 0x08, 0x09},
				},
			},
			// Note, that "blobs" fields should be 128KiB in length, but to keep the test readable,
			// repeated zeros are removed.
			wantJSON: `{
                "maxFeePerBlobGas": "0xb2d05e00",
                "blobVersionedHashes": ["0x6666666666666666666666666666666666666666666666666666666666666666"],
                "blobs": ["0x010203"],
                "commitments": ["0x040506000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"],
                "proofs": ["0x070809000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"]
            }`,
		},
		{
			name: "blobs without sidecars",
			data: &BlobData{
				MaxFeePerBlobGas: big.NewInt(3000000000),
				Blobs: []BlobInfo{
					{
						Hash:    MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", PadNone),
						Sidecar: nil,
					},
				},
			},
			want: &jsonCall{
				MaxFeePerBlobGas: NumberFromBigIntPtr(big.NewInt(3000000000)),
				BlobHashes: []Hash{
					MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", PadNone),
				},
			},
			wantJSON: `{
                "maxFeePerBlobGas": "0xb2d05e00",
                "blobVersionedHashes": ["0x7777777777777777777777777777777777777777777777777777777777777777"]
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test toJSON
			var j jsonCall
			tt.data.toJSON(&j)
			assert.Equal(t, tt.want.MaxFeePerBlobGas, j.MaxFeePerBlobGas)
			assert.Equal(t, tt.want.BlobHashes, j.BlobHashes)
			assert.Equal(t, tt.want.Blobs, j.Blobs)
			assert.Equal(t, tt.want.Commitments, j.Commitments)
			assert.Equal(t, tt.want.Proofs, j.Proofs)

			// Verify generated JSON string
			jsonBytes, err := json.Marshal(j)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(remZerosRx.ReplaceAll(jsonBytes, []byte(""))))

			// Test fromJSON
			var data BlobData
			err = json.Unmarshal(jsonBytes, &j)
			assert.NoError(t, err)
			data.fromJSON(&j)
			assert.Equal(t, tt.data.MaxFeePerBlobGas, data.MaxFeePerBlobGas)
			assert.Equal(t, len(tt.data.Blobs), len(data.Blobs))
			for i := range tt.data.Blobs {
				assert.Equal(t, tt.data.Blobs[i].Hash, data.Blobs[i].Hash)
				if tt.data.Blobs[i].Sidecar == nil {
					assert.Nil(t, data.Blobs[i].Sidecar)
				} else {
					assert.NotNil(t, data.Blobs[i].Sidecar)
					assert.Equal(t, tt.data.Blobs[i].Sidecar.Blob, data.Blobs[i].Sidecar.Blob)
					assert.Equal(t, tt.data.Blobs[i].Sidecar.Commitment, data.Blobs[i].Sidecar.Commitment)
					assert.Equal(t, tt.data.Blobs[i].Sidecar.Proof, data.Blobs[i].Sidecar.Proof)
				}
			}
		})
	}
}
