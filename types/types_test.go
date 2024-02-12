package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
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
			assert.Equal(t, tt.addr, MustAddressFromHex(tt.addr).Checksum(keccak256))
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
				MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
			}),
			want: `"0x1111111111111111111111111111111111111111111111111111111111111111"`,
		},
		{
			arg: (hashList)([]Hash{
				MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
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
				MustAddressFromHex("0x1111111111111111111111111111111111111111"),
			}),
			want: `"0x1111111111111111111111111111111111111111"`,
		},
		{
			arg: (addressList)([]Address{
				MustAddressFromHex("0x1111111111111111111111111111111111111111"),
				MustAddressFromHex("0x2222222222222222222222222222222222222222"),
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
		{arg: `"0x0"`, want: BlockNumberFromUint64(0)},
		{arg: `"0xF"`, want: BlockNumberFromUint64(15)},
		{arg: `"0"`, want: BlockNumberFromUint64(0)},
		{arg: `"F"`, want: BlockNumberFromUint64(15)},
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
		{arg: BlockNumberFromUint64(0), want: `"0x0"`},
		{arg: BlockNumberFromUint64(15), want: `"0xf"`},
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

func Test_SignatureType_Unmarshal(t *testing.T) {
	tests := []struct {
		arg     string
		want    Signature
		wantErr bool
	}{
		{
			arg: `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			want: Signature{
				V: big.NewInt(0),
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
		},
		{
			arg: `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
			want: Signature{
				V: big.NewInt(0),
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
		},
		{
			arg: `"0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000021b"`,
			want: Signature{
				V: big.NewInt(27),
				R: big.NewInt(1),
				S: big.NewInt(2),
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			v := &Signature{}
			err := v.UnmarshalJSON([]byte(tt.arg))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.want.Equal(*v))
			}
		})
	}
}

func Test_SignatureType_Marshal(t *testing.T) {
	tests := []struct {
		signature Signature
		want      string
		wantErr   bool
	}{
		{
			signature: Signature{},
			want:      `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
		},
		{
			signature: Signature{
				V: big.NewInt(0),
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
			want: `"0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`,
		},
		{
			signature: Signature{
				V: big.NewInt(27),
				R: big.NewInt(1),
				S: big.NewInt(2),
			},
			want: `"0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000021b"`,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			j, err := tt.signature.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(j))
		})
	}
}

func Test_SignatureType_Equal(t *testing.T) {
	tests := []struct {
		a, b Signature
		want bool
	}{
		{
			a:    Signature{},
			b:    Signature{},
			want: true,
		},
		{
			a: Signature{},
			b: Signature{
				V: big.NewInt(0),
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
			want: true,
		},
		{
			a: Signature{
				V: big.NewInt(0),
				R: nil,
				S: big.NewInt(0),
			},
			b: Signature{
				V: nil,
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
			want: true,
		},
		{
			a: Signature{
				V: big.NewInt(27),
				R: big.NewInt(1),
				S: big.NewInt(2),
			},
			b: Signature{
				V: big.NewInt(27),
				R: big.NewInt(1),
				S: big.NewInt(2),
			},
			want: true,
		},
		{
			a: Signature{
				V: big.NewInt(27),
				R: nil,
				S: big.NewInt(2),
			},
			b: Signature{
				V: nil,
				R: big.NewInt(2),
				S: big.NewInt(2),
			},
			want: false,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Equal(tt.b))
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
			want: Hash{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
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
			want: Hash{0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
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

func keccak256(data ...[]byte) Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range data {
		h.Write(i)
	}
	return MustHashFromBytes(h.Sum(nil), PadNone)
}
