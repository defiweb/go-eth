package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/crypto/primitives"
	"github.com/defiweb/go-eth/hexutil"
)

func mustKZGHashFromHex(h string) (k primitives.KZGHash) {
	b, err := hexutil.HexToBytes(h)
	if err != nil {
		return primitives.KZGHash{}
	}
	copy(k[:], b)
	return k
}

func assertEqualTX(t *testing.T, actual, expected Transaction) {
	assert.Equal(t, deref(reflect.TypeOf(actual)), deref(reflect.TypeOf(expected)))
	if _, ok := expected.(HasSigningData); ok {
		assert.Equal(t, expected.(HasSigningData).GetSigningData(), actual.(HasSigningData).GetSigningData())
	}
	if _, ok := expected.(HasExecutionData); ok {
		assert.Equal(t, actual.(HasExecutionData).GetExecutionData(), actual.(HasExecutionData).GetExecutionData())
	}
	if _, ok := actual.(HasLegacyFeeData); ok {
		assert.Equal(t, expected.(HasLegacyFeeData).GetLegacyFeeData(), actual.(HasLegacyFeeData).GetLegacyFeeData())
	}
	if _, ok := actual.(HasAccessListData); ok {
		assert.Equal(t, expected.(HasAccessListData).GetAccessListData(), actual.(HasAccessListData).GetAccessListData())
	}
	if _, ok := actual.(HasDynamicFeeData); ok {
		assert.Equal(t, expected.(HasDynamicFeeData).GetDynamicFeeData(), actual.(HasDynamicFeeData).GetDynamicFeeData())
	}
	if _, ok := actual.(HasBlobData); ok {
		assert.Equal(t, expected.(HasBlobData).GetBlobData(), actual.(HasBlobData).GetBlobData())
	}
}

func assertEqualCall(t *testing.T, actual, expected Call) {
	assert.Equal(t, deref(reflect.TypeOf(actual)), deref(reflect.TypeOf(expected)))
	if _, ok := expected.(HasExecutionData); ok {
		assert.Equal(t, actual.(HasExecutionData).GetExecutionData(), actual.(HasExecutionData).GetExecutionData())
	}
	if _, ok := actual.(HasLegacyFeeData); ok {
		assert.Equal(t, expected.(HasLegacyFeeData).GetLegacyFeeData(), actual.(HasLegacyFeeData).GetLegacyFeeData())
	}
	if _, ok := actual.(HasAccessListData); ok {
		assert.Equal(t, expected.(HasAccessListData).GetAccessListData(), actual.(HasAccessListData).GetAccessListData())
	}
	if _, ok := actual.(HasDynamicFeeData); ok {
		assert.Equal(t, expected.(HasDynamicFeeData).GetDynamicFeeData(), actual.(HasDynamicFeeData).GetDynamicFeeData())
	}
	if _, ok := actual.(HasBlobData); ok {
		assert.Equal(t, expected.(HasBlobData).GetBlobData(), actual.(HasBlobData).GetBlobData())
	}
}

func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	return t
}

func ptr[T any](x T) *T {
	return &x
}
