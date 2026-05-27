package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallAccessList_JSON(t *testing.T) {
	tests := []struct {
		name     string
		call     *CallAccessList
		wantJSON string
	}{
		{
			name:     "all fields nil",
			call:     &CallAccessList{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			call: &CallAccessList{
				ExecutionData: ExecutionData{
					From:     MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
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
			wantJSON: `{
                "from": "0x1111111111111111111111111111111111111111",
                "to": "0x2222222222222222222222222222222222222222",
                "value": "0xde0b6b3a7640000",
                "gas": "0x186a0",
                "input": "0x01020304",
                "gasPrice": "0x3b9aca00",
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
			// Encode to JSON
			jsonBytes, err := tt.call.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))

			// Decode from JSON
			call := NewCallAccessList()
			err = call.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded call
			call.From = tt.call.From
			assertEqualCall(t, call, tt.call)
		})
	}
}
