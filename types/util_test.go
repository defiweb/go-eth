package types

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualTX(t *testing.T, actual, expected Transaction) {
	assert.Equal(t, deref(reflect.TypeOf(actual)), deref(reflect.TypeOf(expected)))
	assert.Equal(t, actual.TransactionData(), expected.TransactionData())
	if _, ok := expected.(CallData); ok {
		assert.Equal(t, actual.(CallData).CallData(), actual.(CallData).CallData())
	}
	if _, ok := actual.(LegacyPriceData); ok {
		assert.Equal(t, expected.(LegacyPriceData).LegacyPriceData(), actual.(LegacyPriceData).LegacyPriceData())
	}
	if _, ok := actual.(AccessListData); ok {
		assert.Equal(t, expected.(AccessListData).AccessListData(), actual.(AccessListData).AccessListData())
	}
	if _, ok := actual.(DynamicFeeData); ok {
		assert.Equal(t, expected.(DynamicFeeData).DynamicFeeData(), actual.(DynamicFeeData).DynamicFeeData())
	}
	if _, ok := actual.(BlobData); ok {
		assert.Equal(t, expected.(BlobData).BlobData(), actual.(BlobData).BlobData())
	}
}

func assertEqualCall(t *testing.T, actual, expected Call) {
	assert.Equal(t, deref(reflect.TypeOf(actual)), deref(reflect.TypeOf(expected)))
	if _, ok := expected.(CallData); ok {
		assert.Equal(t, actual.(CallData).CallData(), actual.(CallData).CallData())
	}
	if _, ok := actual.(LegacyPriceData); ok {
		assert.Equal(t, expected.(LegacyPriceData).LegacyPriceData(), actual.(LegacyPriceData).LegacyPriceData())
	}
	if _, ok := actual.(AccessListData); ok {
		assert.Equal(t, expected.(AccessListData).AccessListData(), actual.(AccessListData).AccessListData())
	}
	if _, ok := actual.(DynamicFeeData); ok {
		assert.Equal(t, expected.(DynamicFeeData).DynamicFeeData(), actual.(DynamicFeeData).DynamicFeeData())
	}
	if _, ok := actual.(BlobData); ok {
		assert.Equal(t, expected.(BlobData).BlobData(), actual.(BlobData).BlobData())
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
