package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallLegacy_JSON(t *testing.T) {
	tests := []struct {
		name     string
		call     *CallLegacy
		wantJSON string
	}{
		{
			name:     "all fields nil",
			call:     &CallLegacy{},
			wantJSON: `{}`,
		},
		{
			name: "all fields set",
			call: &CallLegacy{
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
			wantJSON: `{
                "gasPrice": "0x3b9aca00",
                "gas": "0x186a0",
                "to": "0x2222222222222222222222222222222222222222",
                "value": "0xde0b6b3a7640000",
                "input": "0x01020304",
                "from": "0x1111111111111111111111111111111111111111"
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
			call := NewCallLegacy()
			err = call.UnmarshalJSON(jsonBytes)
			require.NoError(t, err)

			// Compare the original and decoded call
			call.From = tt.call.From
			assertEqualCall(t, call, tt.call)
		})
	}
}
