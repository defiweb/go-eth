package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "foo((uint256,bytes32)[])(uint256)", expected: "function foo((uint256, bytes32)[]) returns (uint256)"},
		{signature: "foo((uint256 a, bytes32 b)[] c)(uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) returns (uint256 d)"},
		{signature: "function foo(tuple(uint256 a, bytes32 b)[] memory c) pure returns (uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) returns (uint256 d)"},
		{signature: "event foo(uint256)", wantErr: true},
		{signature: "error foo(uint256)", wantErr: true},
		{signature: "constructor(uint256)", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			m, err := ParseMethod(tt.signature)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, m.String())
			}
		})
	}
}

func TestMethod_EncodeArgs(t *testing.T) {
	const ZeroAddr = "0x0000000000000000000000000000000000000000"
	const FuncERC20Approve = "function approve(address spender,uint256 value) external returns (bool)"
	const FuncTryAggregate = "function tryAggregate(bool requireSuccess, Call[] memory calls) public returns (Result[] memory returnData)"

	innerFunc := MustParseMethod(FuncERC20Approve)

	const TypeCall = "(address target,bytes callData)"
	type Call struct {
		Address types.Address
		Data    []byte
	}
	require.NoError(t, addType("Call", TypeCall))
	const TypeResult = "(bool success,bytes returnData)"
	type Result struct {
		Success bool
		Data    []byte
	}
	require.NoError(t, addType("Result", TypeResult))
	outerFunc := MustParseMethod(FuncTryAggregate)

	data := make([]Call, 2)
	data[0] = Call{
		Address: types.MustHexToAddress(ZeroAddr),
		Data:    mustFn[[]byte](t)(innerFunc.EncodeArgs(types.MustHexToAddress(ZeroAddr), 0)),
	}
	data[1] = Call{
		Address: types.MustHexToAddress(ZeroAddr),
		Data:    mustFn[[]byte](t)(innerFunc.EncodeArgs(types.MustHexToAddress(ZeroAddr), 1)),
	}

	mustFn[[]byte](t)(outerFunc.EncodeArgs(false, data))
}

func addType(name, signature string) error {
	var err error
	Default.Types[name], err = ParseType(signature)
	return err
}
func mustFn[C any](t *testing.T) func(c C, err error) C {
	return func(c C, err error) C {
		require.NoError(t, err)
		return c
	}
}
