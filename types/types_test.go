package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func Test_AddressFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		want    Address
		wantErr bool
	}{
		{
			arg:  "0x00112233445566778899aabbccddeeff00112233",
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:  "00112233445566778899aabbccddeeff00112233",
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:     "0x00",
			wantErr: true,
		},
		{
			arg:     "0x00112233445566778899aabbccddeeff0011223344",
			wantErr: true,
		},
		{
			arg:     "invalid",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v, err := AddressFromHex(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, v)
			}
		})
	}
}

func Test_AddressFromBytes(t *testing.T) {
	tests := []struct {
		arg     []byte
		want    Address
		wantErr bool
	}{
		{
			arg:  []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33},
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:     []byte{},
			wantErr: true,
		},
		{
			arg:     []byte{0x00},
			wantErr: true,
		},
		{
			arg:     []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44},
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := AddressFromBytes(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_AddressType_String(t *testing.T) {
	addr := (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33})
	assert.Equal(t, "0x00112233445566778899aabbccddeeff00112233", addr.String())
}

func Test_AddressType_Bytes(t *testing.T) {
	addr := (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33})
	expected := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}
	assert.Equal(t, expected, addr.Bytes())
}

func Test_AddressType_IsZero(t *testing.T) {
	tests := []struct {
		addr   Address
		isZero bool
	}{
		{addr: Address{}, isZero: true},
		{addr: (Address)([AddressLength]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), isZero: false},
		{addr: (Address)([AddressLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}), isZero: false},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.isZero, tt.addr.IsZero())
		})
	}
}

func Test_AddressType_Checksum(t *testing.T) {
	tests := []struct {
		addr string
	}{
		{addr: "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359"},
		{addr: "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"},
		{addr: "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359"},
		{addr: "0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB"},
		{addr: "0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.addr, MustAddressFromHex(tt.addr).Checksum())
		})
	}
}

func Test_AddressType_MarshalJSON(t *testing.T) {
	tests := []struct {
		arg  Address
		want string
	}{
		{
			arg:  (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
			want: `"0x00112233445566778899aabbccddeeff00112233"`,
		},
		{
			arg:  Address{},
			want: `"0x0000000000000000000000000000000000000000"`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			jsonBytes, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonBytes))
		})
	}
}

func Test_AddressType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg     string
		want    Address
		wantErr bool
	}{
		{
			arg:  `"0x00112233445566778899aabbccddeeff00112233"`,
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:  `"00112233445566778899aabbccddeeff00112233"`,
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:     `0x00112233445566778899aabbccddeeff00112233`,
			wantErr: true,
		},
		{
			arg:     `"0x00"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff0011223344"`,
			wantErr: true,
		},
		{
			arg:     `"invalid"`,
			wantErr: true,
		},
		{
			arg:     `"""`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Address{}
			err := v.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *v)
			}
		})
	}
}

func Test_AddressType_MarshalText(t *testing.T) {
	tests := []struct {
		arg  Address
		want string
	}{
		{
			arg:  (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
			want: `0x00112233445566778899aabbccddeeff00112233`,
		},
		{
			arg:  Address{},
			want: `0x0000000000000000000000000000000000000000`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			textBytes, err := tt.arg.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(textBytes))
		})
	}
}

func Test_AddressType_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg     string
		want    Address
		wantErr bool
	}{
		{
			arg:  `0x00112233445566778899aabbccddeeff00112233`,
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:  `00112233445566778899aabbccddeeff00112233`,
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff00112233"`,
			wantErr: true,
		},
		{
			arg:     `0x00`,
			wantErr: true,
		},
		{
			arg:     `0x00112233445566778899aabbccddeeff0011223344`,
			wantErr: true,
		},
		{
			arg:     `invalid`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Address{}
			err := v.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *v)
			}
		})
	}
}

func Test_AddressType_EncodeRLP(t *testing.T) {
	tests := []struct {
		addr Address
		want []byte
	}{
		{
			addr: Address{},
			want: []byte{0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			addr: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
			want: []byte{0x94, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			rlpBytes, err := tt.addr.EncodeRLP()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, rlpBytes)
		})
	}
}

func Test_AddressType_DecodeRLP(t *testing.T) {
	tests := []struct {
		data    []byte
		want    Address
		wantErr bool
	}{
		{
			data: []byte{0x80},
			want: Address{},
		},
		{
			data: []byte{0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: Address{},
		},
		{
			data: []byte{0x94, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33},
			want: (Address)([AddressLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}),
		},
		{
			data:    []byte{0x81},
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var addr Address
			_, err := addr.DecodeRLP(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, addr)
			}
		})
	}
}

func Test_HashFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		pad     Pad
		want    Hash
		wantErr bool
	}{
		{
			arg:  "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff",
			pad:  PadNone,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:  "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff",
			pad:  PadNone,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:  "80ff",
			pad:  PadLeft,
			want: (Hash)([HashLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xff}),
		},
		{
			arg:  "80ff",
			pad:  PadRight,
			want: (Hash)([HashLength]byte{0x80, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
		},
		{
			arg:     "0x00",
			pad:     PadNone,
			wantErr: true,
		},
		{
			arg:     "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00",
			pad:     PadNone,
			wantErr: true,
		},
		{
			arg:     "invalid",
			pad:     PadLeft,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v, err := HashFromHex(tt.arg, tt.pad)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, v)
			}
		})
	}
}

func Test_HashFromBytes(t *testing.T) {
	tests := []struct {
		arg     []byte
		pad     Pad
		want    Hash
		wantErr bool
	}{
		{
			arg:     []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			pad:     PadNone,
			want:    (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
			wantErr: false,
		},
		{
			arg:  []byte{0x80, 0xff},
			pad:  PadLeft,
			want: (Hash)([HashLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xff}),
		},
		{
			arg:  []byte{0x80, 0xff},
			pad:  PadRight,
			want: (Hash)([HashLength]byte{0x80, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
		},
		{
			arg:     []byte{0x00, 0x11, 0x22},
			pad:     PadNone,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := HashFromBytes(tt.arg, tt.pad)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_HashFromBigInt(t *testing.T) {
	tests := []struct {
		i       *big.Int
		want    Hash
		wantErr bool
	}{
		{
			i:    big.NewInt(0),
			want: Hash{},
		},
		{
			i:    big.NewInt(1),
			want: Hash{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			i:    big.NewInt(-1),
			want: Hash{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		// max uint256
		{
			i:    new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(256)), big.NewInt(1)),
			want: Hash{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		// min int256
		{
			i:    new(big.Int).Lsh(big.NewInt(-1), uint(255)),
			want: Hash{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		// max uint256 + 1
		{
			i:       new(big.Int).Add(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(256)), big.NewInt(1)), big.NewInt(1)),
			wantErr: true,
		},
		// min int256 - 1
		{
			i:       new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(-1), uint(255)), big.NewInt(1)),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := HashFromBigInt(tt.i)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_HashType_String(t *testing.T) {
	hash := (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	assert.Equal(t, "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff", hash.String())
}

func Test_HashType_Bytes(t *testing.T) {
	hash := (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff})
	expected := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	assert.Equal(t, expected, hash.Bytes())
}

func Test_HashType_IsZero(t *testing.T) {
	tests := []struct {
		hash   Hash
		isZero bool
	}{
		{hash: Hash{}, isZero: true},
		{hash: (Hash)([HashLength]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}), isZero: false},
		{hash: (Hash)([HashLength]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}), isZero: false},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.isZero, tt.hash.IsZero())
		})
	}
}

func Test_HashType_MarshalJSON(t *testing.T) {
	tests := []struct {
		arg  Hash
		want string
	}{
		{
			arg:  (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
			want: `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
		},
		{
			arg:  Hash{},
			want: `"0x0000000000000000000000000000000000000000000000000000000000000000"`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			jsonBytes, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonBytes))
		})
	}
}

func Test_HashType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg     string
		want    Hash
		wantErr bool
	}{
		{
			arg:  `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:  `"00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:     `0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`,
			wantErr: true,
		},
		{
			arg:     `"0x00"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"`,
			wantErr: true,
		},
		{
			arg:     `"invalid"`,
			wantErr: true,
		},
		{
			arg:     `"""`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Hash{}
			err := v.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *v)
			}
		})
	}
}

func Test_HashType_MarshalText(t *testing.T) {
	tests := []struct {
		arg  Hash
		want string
	}{
		{
			arg:  (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
			want: `0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`,
		},
		{
			arg:  Hash{},
			want: `0x0000000000000000000000000000000000000000000000000000000000000000`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			textBytes, err := tt.arg.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(textBytes))
		})
	}
}

func Test_HashType_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg     string
		want    Hash
		wantErr bool
	}{
		{
			arg:  `0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:  `00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff`,
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
			wantErr: true,
		},
		{
			arg:     `0x00`,
			wantErr: true,
		},
		{
			arg:     `0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00`,
			wantErr: true,
		},
		{
			arg:     `invalid`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Hash{}
			err := v.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *v)
			}
		})
	}
}

func Test_HashType_EncodeRLP(t *testing.T) {
	tests := []struct {
		hash Hash
		want []byte
	}{
		{
			hash: Hash{},
			want: []byte{0xa0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			hash: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
			want: []byte{0xa0, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			rlpBytes, err := tt.hash.EncodeRLP()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, rlpBytes)
		})
	}
}

func Test_HashType_DecodeRLP(t *testing.T) {
	tests := []struct {
		data    []byte
		want    Hash
		wantErr bool
	}{
		{
			data: []byte{0x80},
			want: Hash{},
		},
		{
			data: []byte{0xa0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: Hash{},
		},
		{
			data: []byte{0xa0, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			want: (Hash)([HashLength]byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
		},
		{
			data:    []byte{0x81},
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var hash Hash
			_, err := hash.DecodeRLP(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, hash)
			}
		})
	}
}

func Test_BlockNumberFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		want    BlockNumber
		wantErr bool
	}{
		{arg: "0x0", want: BlockNumberFromUint64(0)},
		{arg: "0xf", want: BlockNumberFromUint64(15)},
		{arg: "earliest", want: EarliestBlockNumber},
		{arg: "latest", want: LatestBlockNumber},
		{arg: "pending", want: PendingBlockNumber},
		{arg: "safe", want: SafeBlockNumber},
		{arg: "finalized", want: FinalizedBlockNumber},
		{arg: "foo", wantErr: true},
		{arg: "0xZZ", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			bn, err := BlockNumberFromHex(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), bn.String())
			}
		})
	}
}

func Test_BlockNumberType_Big(t *testing.T) {
	tests := []struct {
		arg  BlockNumber
		want *big.Int
	}{
		{arg: BlockNumberFromUint64(0), want: big.NewInt(0)},
		{arg: BlockNumberFromUint64(15), want: big.NewInt(15)},
		{arg: EarliestBlockNumber, want: big.NewInt(-1)},
		{arg: LatestBlockNumber, want: big.NewInt(-2)},
		{arg: PendingBlockNumber, want: big.NewInt(-3)},
		{arg: SafeBlockNumber, want: big.NewInt(-4)},
		{arg: FinalizedBlockNumber, want: big.NewInt(-5)},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.True(t, tt.arg.Big().Cmp(tt.want) == 0)
		})
	}
}

func Test_BlockNumberType_String(t *testing.T) {
	tests := []struct {
		arg  BlockNumber
		want string
	}{
		{arg: BlockNumberFromUint64(0), want: `0x0`},
		{arg: BlockNumberFromUint64(15), want: `0xf`},
		{arg: EarliestBlockNumber, want: `earliest`},
		{arg: LatestBlockNumber, want: `latest`},
		{arg: PendingBlockNumber, want: `pending`},
		{arg: SafeBlockNumber, want: `safe`},
		{arg: FinalizedBlockNumber, want: `finalized`},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.String())
		})
	}
}

func Test_BlockNumberType_MarshalJSON(t *testing.T) {
	tests := []struct {
		arg  BlockNumber
		want string
	}{
		{arg: BlockNumberFromUint64(0), want: `"0x0"`},
		{arg: BlockNumberFromUint64(15), want: `"0xf"`},
		{arg: EarliestBlockNumber, want: `"earliest"`},
		{arg: LatestBlockNumber, want: `"latest"`},
		{arg: PendingBlockNumber, want: `"pending"`},
		{arg: SafeBlockNumber, want: `"safe"`},
		{arg: FinalizedBlockNumber, want: `"finalized"`},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			jsonBytes, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonBytes))
		})
	}
}

func Test_BlockNumberType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg         string
		want        BlockNumber
		wantErr     bool
		isTag       bool
		isEarliest  bool
		isLatest    bool
		isPending   bool
		isSafe      bool
		isFinalized bool
	}{
		{arg: `"0x0"`, want: BlockNumberFromUint64(0)},
		{arg: `"0xF"`, want: BlockNumberFromUint64(15)},
		{arg: `"0"`, want: BlockNumberFromUint64(0)},
		{arg: `"F"`, want: BlockNumberFromUint64(15)},
		{arg: `"earliest"`, want: EarliestBlockNumber, isTag: true, isEarliest: true},
		{arg: `"latest"`, want: LatestBlockNumber, isTag: true, isLatest: true},
		{arg: `"pending"`, want: PendingBlockNumber, isTag: true, isPending: true},
		{arg: `"safe"`, want: SafeBlockNumber, isTag: true, isSafe: true},
		{arg: `"finalized"`, want: FinalizedBlockNumber, isTag: true, isFinalized: true},
		{arg: `"foo"`, wantErr: true},
		{arg: `"0xZZ"`, wantErr: true},
		{arg: `0x0`, wantErr: true},
		{arg: `latest`, wantErr: true},
		{arg: `"""`, wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &BlockNumber{}
			err := v.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), (*v).String())
				assert.Equal(t, tt.isTag, v.IsTag())
				assert.Equal(t, tt.isEarliest, v.IsEarliest())
				assert.Equal(t, tt.isLatest, v.IsLatest())
				assert.Equal(t, tt.isPending, v.IsPending())
				assert.Equal(t, tt.isSafe, v.IsSafe())
				assert.Equal(t, tt.isFinalized, v.IsFinalized())
			}
		})
	}
}

func Test_BlockNumberType_MarshalText(t *testing.T) {
	tests := []struct {
		arg  BlockNumber
		want string
	}{
		{arg: BlockNumberFromUint64(0), want: `0x0`},
		{arg: BlockNumberFromUint64(15), want: `0xf`},
		{arg: EarliestBlockNumber, want: `earliest`},
		{arg: LatestBlockNumber, want: `latest`},
		{arg: PendingBlockNumber, want: `pending`},
		{arg: SafeBlockNumber, want: `safe`},
		{arg: FinalizedBlockNumber, want: `finalized`},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			textBytes, err := tt.arg.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(textBytes))
		})
	}
}

func Test_BlockNumberType_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg         string
		want        BlockNumber
		wantErr     bool
		isTag       bool
		isEarliest  bool
		isLatest    bool
		isPending   bool
		isSafe      bool
		isFinalized bool
	}{
		{arg: `0x0`, want: BlockNumberFromUint64(0)},
		{arg: `0xF`, want: BlockNumberFromUint64(15)},
		{arg: `0`, want: BlockNumberFromUint64(0)},
		{arg: `F`, want: BlockNumberFromUint64(15)},
		{arg: `earliest`, want: EarliestBlockNumber, isTag: true, isEarliest: true},
		{arg: `latest`, want: LatestBlockNumber, isTag: true, isLatest: true},
		{arg: `pending`, want: PendingBlockNumber, isTag: true, isPending: true},
		{arg: `safe`, want: SafeBlockNumber, isTag: true, isSafe: true},
		{arg: `finalized`, want: FinalizedBlockNumber, isTag: true, isFinalized: true},
		{arg: `foo`, wantErr: true},
		{arg: `0xZZ`, wantErr: true},
		{arg: `"0x0"`, wantErr: true},
		{arg: `"latest"`, wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &BlockNumber{}
			err := v.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), (*v).String())
				assert.Equal(t, tt.isTag, v.IsTag())
				assert.Equal(t, tt.isEarliest, v.IsEarliest())
				assert.Equal(t, tt.isLatest, v.IsLatest())
				assert.Equal(t, tt.isPending, v.IsPending())
				assert.Equal(t, tt.isSafe, v.IsSafe())
				assert.Equal(t, tt.isFinalized, v.IsFinalized())
			}
		})
	}
}

var testSignature = Signature{
	V: func() *big.Int {
		v, _ := new(big.Int).SetString("37", 10)
		return v
	}(),
	R: func() *big.Int {
		v, _ := new(big.Int).SetString("18515461264373351373200002665853028612451056578545711640558177340181847433846", 10)
		return v
	}(),
	S: func() *big.Int {
		v, _ := new(big.Int).SetString("46948507304638947509940763649030358759909902576025900602547168820602576006531", 10)
		return v
	}(),
}

func Test_SignatureFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		want    Signature
		wantErr bool
	}{
		{
			arg:  "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325",
			want: testSignature,
		},
		{
			arg:     "0x00",
			wantErr: true,
		},
		{
			arg:     "invalid",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := SignatureFromHex(tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_SignatureFromBytes(t *testing.T) {
	tests := []struct {
		arg     string
		want    Signature
		wantErr bool
	}{
		{
			arg:  "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325",
			want: testSignature,
		},
		{
			arg:     "0x00",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := SignatureFromBytes(hexutil.MustHexToBytes(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_SignatureType_String(t *testing.T) {
	assert.Equal(t, "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325", testSignature.String())
}

func Test_SignatureType_Bytes(t *testing.T) {
	assert.Equal(t, "0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325", hexutil.BytesToHex(testSignature.Bytes()))
}

func Test_SignatureType_IsZero(t *testing.T) {
	tests := []struct {
		arg  Signature
		want bool
	}{
		{arg: Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)}, want: true},
		{arg: Signature{V: nil, R: nil, S: nil}, want: true},
		{arg: Signature{V: nil, R: big.NewInt(0), S: big.NewInt(0)}, want: true},
		{arg: Signature{V: big.NewInt(0), R: nil, S: nil}, want: true},
		{arg: Signature{V: big.NewInt(0), R: big.NewInt(0), S: nil}, want: true},
		{arg: Signature{V: big.NewInt(1), R: big.NewInt(0), S: big.NewInt(0)}, want: false},
		{arg: Signature{V: big.NewInt(0), R: big.NewInt(1), S: big.NewInt(0)}, want: false},
		{arg: Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(1)}, want: false},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.IsZero())
		})
	}
}

func Test_SignatureType_Equal(t *testing.T) {
	tests := []struct {
		arg1 Signature
		arg2 Signature
		want bool
	}{
		{
			arg1: Signature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			arg2: Signature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			want: true,
		},
		{
			arg1: Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)},
			arg2: Signature{V: nil, R: nil, S: nil},
			want: true,
		},
		{
			arg1: Signature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			arg2: Signature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(4)},
			want: false,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg1.Equal(tt.arg2))
		})
	}
}

func Test_SignatureType_MarshalJSON(t *testing.T) {
	tests := []struct {
		arg  Signature
		want string
	}{
		{
			arg:  testSignature,
			want: `"0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325"`,
		},
		{
			arg:  Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)},
			want: `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
		},
		{
			arg:  Signature{V: nil, R: nil, S: nil},
			want: `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			jsonBytes, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(jsonBytes))
		})
	}
}

func Test_SignatureType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg     string
		want    Signature
		wantErr bool
	}{
		{
			arg:  `"0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325"`,
			want: testSignature,
		},
		{
			arg:     `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			want:    Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)},
			wantErr: false,
		},
		{
			arg:     `0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
			wantErr: true,
		},
		{
			arg:     `"0x00"`,
			wantErr: true,
		},
		{
			arg:     `"invalid"`,
			wantErr: true,
		},
		{
			arg:     `"""`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var sig Signature
			err := sig.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.want.Equal(sig))
			}
		})
	}
}

func Test_SignatureType_MarshalText(t *testing.T) {
	tests := []struct {
		arg  Signature
		want string
	}{
		{
			arg:  testSignature,
			want: `0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325`,
		},
		{
			arg:  Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)},
			want: `0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			textBytes, err := tt.arg.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(textBytes))
		})
	}
}

func Test_SignatureType_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg     string
		want    Signature
		wantErr bool
	}{
		{
			arg:  `0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa63627667cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d8325`,
			want: testSignature,
		},
		{
			arg:     `0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000`,
			want:    Signature{V: big.NewInt(0), R: big.NewInt(0), S: big.NewInt(0)},
			wantErr: false,
		},
		{
			arg:     `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			wantErr: true,
		},
		{
			arg:     `0x00`,
			wantErr: true,
		},
		{
			arg:     `invalid`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var sig Signature
			err := sig.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.want.Equal(sig))
			}
		})
	}
}

func Test_NumberFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		want    Number
		wantErr bool
	}{
		{
			arg:  "0x0",
			want: Number{x: *big.NewInt(0)},
		},
		{
			arg:  "0x10",
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:  "0",
			want: Number{x: *big.NewInt(0)},
		},
		{
			arg:  "10",
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:  "-0x10",
			want: Number{x: *big.NewInt(-0x10)},
		},
		{
			arg:  "-10",
			want: Number{x: *big.NewInt(-0x10)},
		},
		{
			arg:     "invalid",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := NumberFromHex(tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.String(), got.String())
			}
		})
	}
}

func Test_NumberType_Big(t *testing.T) {
	tests := []struct {
		arg  Number
		want *big.Int
	}{
		{
			arg:  Number{x: *big.NewInt(0)},
			want: big.NewInt(0),
		},
		{
			arg:  Number{x: *big.NewInt(0x10)},
			want: big.NewInt(0x10),
		},
		{
			arg:  Number{x: *big.NewInt(-0x10)},
			want: big.NewInt(-0x10),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want.Cmp(tt.arg.Big()), 0)
		})
	}
}

func Test_NumberType_String(t *testing.T) {
	tests := []struct {
		arg  Number
		want string
	}{
		{
			arg:  Number{x: *big.NewInt(0)},
			want: "0x0",
		},
		{
			arg:  Number{x: *big.NewInt(0x10)},
			want: "0x10",
		},
		{
			arg:  Number{x: *big.NewInt(-0x10)},
			want: "-0x10",
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.String())
		})
	}
}

func Test_NumberType_Bytes(t *testing.T) {
	tests := []struct {
		arg  Number
		want []byte
	}{
		{
			arg:  Number{x: *big.NewInt(0)},
			want: []byte{},
		},
		{
			arg:  Number{x: *big.NewInt(0x10)},
			want: []byte{0x10},
		},
		{
			arg:  Number{x: *big.NewInt(-0x10)},
			want: []byte{0x10},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.Bytes())
		})
	}
}

func Test_NumberType_MarshalJSON(t *testing.T) {
	tests := []struct {
		arg  Number
		want string
	}{
		{
			arg:  Number{x: *big.NewInt(0)},
			want: `"0x0"`,
		},
		{
			arg:  Number{x: *big.NewInt(0x10)},
			want: `"0x10"`,
		},
		{
			arg:  Number{x: *big.NewInt(-0x10)},
			want: `"-0x10"`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := tt.arg.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func Test_NumberType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg     string
		want    Number
		wantErr bool
	}{
		{
			arg:  `"0x0"`,
			want: Number{x: *big.NewInt(0)},
		},
		{
			arg:  `"0x10"`,
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:  `"-10"`,
			want: Number{x: *big.NewInt(-0x10)},
		},
		{
			arg:  `"10"`,
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:     `0x10"`,
			wantErr: true,
		},
		{
			arg:     `10"`,
			wantErr: true,
		},
		{
			arg:     `"invalid"`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var got Number
			err := got.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_NumberType_MarshalText(t *testing.T) {
	tests := []struct {
		arg  Number
		want string
	}{
		{
			arg:  Number{x: *big.NewInt(0)},
			want: `0x0`,
		},
		{
			arg:  Number{x: *big.NewInt(0x10)},
			want: `0x10`,
		},
		{
			arg:  Number{x: *big.NewInt(-0x10)},
			want: `-0x10`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			textBytes, err := tt.arg.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(textBytes))
		})
	}
}

func Test_NumberType_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg     string
		want    Number
		wantErr bool
	}{
		{
			arg:  `0x0`,
			want: Number{x: *big.NewInt(0)},
		},
		{
			arg:  `0x10`,
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:  `-10`,
			want: Number{x: *big.NewInt(-0x10)},
		},
		{
			arg:  `10`,
			want: Number{x: *big.NewInt(0x10)},
		},
		{
			arg:     `"0x10"`,
			wantErr: true,
		},
		{
			arg:     `"10"`,
			wantErr: true,
		},
		{
			arg:     `invalid`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var got Number
			err := got.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_BytesFromHex(t *testing.T) {
	tests := []struct {
		arg     string
		want    Bytes
		wantErr bool
	}{
		{
			arg:  "0x00112233",
			want: Bytes{0x00, 0x11, 0x22, 0x33},
		},
		{
			arg:  "0x00112233",
			want: Bytes{0x00, 0x11, 0x22, 0x33},
		},
		{
			arg:     "invalid",
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := BytesFromHex(tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_Bytes_PadLeft(t *testing.T) {
	tests := []struct {
		arg  Bytes
		len  int
		want Bytes
	}{
		{
			arg:  Bytes{0x01, 0x02},
			len:  0,
			want: Bytes{},
		},
		{
			arg:  Bytes{0x01, 0x02},
			len:  1,
			want: Bytes{0x02},
		},
		{
			arg:  Bytes{0x01, 0x02},
			len:  4,
			want: Bytes{0x00, 0x00, 0x01, 0x02},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.PadLeft(tt.len))
		})
	}
}

func Test_Bytes_PadRight(t *testing.T) {
	tests := []struct {
		arg  Bytes
		len  int
		want Bytes
	}{
		{
			arg:  Bytes{0x01, 0x02},
			len:  0,
			want: Bytes{},
		},
		{
			arg:  Bytes{0x01, 0x02},
			len:  1,
			want: Bytes{0x01},
		},
		{
			arg:  Bytes{0x01, 0x02},
			len:  4,
			want: Bytes{0x01, 0x02, 0x00, 0x00},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.PadRight(tt.len))
		})
	}
}

func Test_Bytes_Bytes(t *testing.T) {
	b := Bytes{0x01, 0x02, 0x03}
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, b.Bytes())
}

func Test_Bytes_String(t *testing.T) {
	b := Bytes{0x01, 0x02, 0x03}
	assert.Equal(t, "0x010203", b.String())
}

func Test_Bytes_MarshalJSON(t *testing.T) {
	b := Bytes{0x01, 0x02, 0x03}
	data, err := b.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"0x010203"`, string(data))
}

func Test_Bytes_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		arg     string
		want    Bytes
		wantErr bool
	}{
		{
			arg:  `"0x010203"`,
			want: Bytes{0x01, 0x02, 0x03},
		},
		{
			arg:  `"010203"`,
			want: Bytes{0x01, 0x02, 0x03},
		},
		{
			arg:     `0x010203`,
			wantErr: true,
		},
		{
			arg:     `"invalid"`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var b Bytes
			err := b.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, b)
			}
		})
	}
}

func Test_Bytes_MarshalText(t *testing.T) {
	b := Bytes{0x01, 0x02, 0x03}
	data, err := b.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "0x010203", string(data))
}

func Test_Bytes_UnmarshalText(t *testing.T) {
	tests := []struct {
		arg     string
		want    Bytes
		wantErr bool
	}{
		{
			arg:  `0x010203`,
			want: Bytes{0x01, 0x02, 0x03},
		},
		{
			arg:  `010203`,
			want: Bytes{0x01, 0x02, 0x03},
		},
		{
			arg:     `"0x010203"`,
			wantErr: true,
		},
		{
			arg:     `invalid`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var b Bytes
			err := b.UnmarshalText([]byte(tt.arg))
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, b)
			}
		})
	}
}
