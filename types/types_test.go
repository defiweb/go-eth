package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddressType_Unmarshal(t *testing.T) {
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
			arg:     `"00112233445566778899aabbccddeeff0011223344"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff0011223344"`,
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

func Test_AddressType_Marshal(t *testing.T) {
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
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_hashType_Unmarshal(t *testing.T) {
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
			arg:     `"00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"`,
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

func Test_hashType_Marshal(t *testing.T) {
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
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_hashesType_Unmarshal(t *testing.T) {
	tests := []struct {
		arg     string
		want    hashList
		wantErr bool
	}{
		{
			arg:  `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
			want: (hashList)([]Hash{{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}),
		},
		{
			arg:  `"00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"`,
			want: (hashList)([]Hash{{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}}),
		},
		{
			arg: `["0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff", "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"]`,
			want: (hashList)([]Hash{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			}),
		},
		{
			arg:     `"00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00"`,
			wantErr: true,
		},
		{
			arg:     `"""`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &hashList{}
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

func Test_hashesType_Marshal(t *testing.T) {
	tests := []struct {
		arg  hashList
		want string
	}{
		{
			arg: (hashList)([]Hash{
				HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
			}),
			want: `"0x1111111111111111111111111111111111111111111111111111111111111111"`,
		},
		{
			arg: (hashList)([]Hash{
				HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111"),
				HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222"),
			}),
			want: `["0x1111111111111111111111111111111111111111111111111111111111111111","0x2222222222222222222222222222222222222222222222222222222222222222"]`,
		},
		{
			arg:  hashList{},
			want: `[]`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_AddressesType_Unmarshal(t *testing.T) {
	tests := []struct {
		arg     string
		want    addressList
		wantErr bool
	}{
		{
			arg:  `"0x00112233445566778899aabbccddeeff00112233"`,
			want: (addressList)([]Address{{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}}),
		},
		{
			arg:  `"00112233445566778899aabbccddeeff00112233"`,
			want: (addressList)([]Address{{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33}}),
		},
		{
			arg: `["0x00112233445566778899aabbccddeeff00112233", "0x00112233445566778899aabbccddeeff00112233"]`,
			want: (addressList)([]Address{
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33},
				{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33},
			}),
		},
		{
			arg:     `"00112233445566778899aabbccddeeff0011223344"`,
			wantErr: true,
		},
		{
			arg:     `"0x00112233445566778899aabbccddeeff0011223344"`,
			wantErr: true,
		},
		{
			arg:     `"""`,
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &addressList{}
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

func Test_AddressesType_Marshal(t *testing.T) {
	tests := []struct {
		arg  addressList
		want string
	}{
		{
			arg: (addressList)([]Address{
				MustHexToAddress("0x1111111111111111111111111111111111111111"),
			}),
			want: `"0x1111111111111111111111111111111111111111"`,
		},
		{
			arg: (addressList)([]Address{
				MustHexToAddress("0x1111111111111111111111111111111111111111"),
				MustHexToAddress("0x2222222222222222222222222222222222222222"),
			}),
			want: `["0x1111111111111111111111111111111111111111","0x2222222222222222222222222222222222222222"]`,
		},
		{
			arg:  addressList{},
			want: `[]`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_BlockNumberType_Unmarshal(t *testing.T) {
	tests := []struct {
		arg        string
		want       BlockNumber
		wantErr    bool
		isTag      bool
		isEarliest bool
		isLatest   bool
		isPending  bool
	}{
		{arg: `"0x0"`, want: Uint64ToBlockNumber(0)},
		{arg: `"0xF"`, want: Uint64ToBlockNumber(15)},
		{arg: `"0"`, want: Uint64ToBlockNumber(0)},
		{arg: `"F"`, want: Uint64ToBlockNumber(15)},
		{arg: `"earliest"`, want: EarliestBlockNumber, isTag: true, isEarliest: true},
		{arg: `"latest"`, want: LatestBlockNumber, isTag: true, isLatest: true},
		{arg: `"pending"`, want: PendingBlockNumber, isTag: true, isPending: true},
		{arg: `"foo"`, wantErr: true},
		{arg: `"0xZ"`, wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &BlockNumber{}
			err := v.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, *v)
				assert.Equal(t, tt.isTag, v.IsTag())
				assert.Equal(t, tt.isEarliest, v.IsEarliest())
				assert.Equal(t, tt.isLatest, v.IsLatest())
				assert.Equal(t, tt.isPending, v.IsPending())
			}
		})
	}
}

func Test_BlockNumberType_Marshal(t *testing.T) {
	tests := []struct {
		arg  BlockNumber
		want string
	}{
		{arg: Uint64ToBlockNumber(0), want: `"0x0"`},
		{arg: Uint64ToBlockNumber(15), want: `"0xf"`},
		{arg: EarliestBlockNumber, want: `"earliest"`},
		{arg: LatestBlockNumber, want: `"latest"`},
		{arg: PendingBlockNumber, want: `"pending"`},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_BytesType_Unmarshal(t *testing.T) {
	tests := []struct {
		arg     string
		want    Bytes
		wantErr bool
	}{
		{arg: `"0xDEADBEEF"`, want: (Bytes)([]byte{0xDE, 0xAD, 0xBE, 0xEF})},
		{arg: `"DEADBEEF"`, want: (Bytes)([]byte{0xDE, 0xAD, 0xBE, 0xEF})},
		{arg: `"0x"`, want: (Bytes)([]byte{})},
		{arg: `""`, want: (Bytes)([]byte{})},
		{arg: `"0x0"`, want: (Bytes)([]byte{0x0})},
		{arg: `"foo"`, wantErr: true},
		{arg: `"0xZZ"`, wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Bytes{}
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

func Test_BytesType_Marshal(t *testing.T) {
	tests := []struct {
		arg  Bytes
		want string
	}{
		{arg: (Bytes)([]byte{0xDE, 0xAD, 0xBE, 0xEF}), want: `"0xdeadbeef"`},
		{arg: (Bytes)([]byte{}), want: `"0x"`},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			j, err := tt.arg.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}
