package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualTX(t *testing.T, actual, expected Transaction) {
	assert.Equal(t, deref(reflect.TypeOf(actual)), deref(reflect.TypeOf(expected)))
	assert.Equal(t, actual.GetTransactionData(), expected.GetTransactionData())
	if _, ok := expected.(HasCallData); ok {
		assert.Equal(t, actual.(HasCallData).GetCallData(), actual.(HasCallData).GetCallData())
	}
	if _, ok := actual.(HasLegacyPriceData); ok {
		assert.Equal(t, expected.(HasLegacyPriceData).GetLegacyPriceData(), actual.(HasLegacyPriceData).GetLegacyPriceData())
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
	if _, ok := expected.(HasCallData); ok {
		assert.Equal(t, actual.(HasCallData).GetCallData(), actual.(HasCallData).GetCallData())
	}
	if _, ok := actual.(HasLegacyPriceData); ok {
		assert.Equal(t, expected.(HasLegacyPriceData).GetLegacyPriceData(), actual.(HasLegacyPriceData).GetLegacyPriceData())
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
