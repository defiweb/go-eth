package abi

import "github.com/defiweb/go-eth/hexutil"

type FourBytes [4]byte

func (f FourBytes) Bytes() []byte {
	return f[:]
}

func (f FourBytes) Hex() string {
	return hexutil.BytesToHex(f[:])
}

func (f FourBytes) String() string {
	return f.Hex()
}

func (f FourBytes) Match(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return f == FourBytes{data[0], data[1], data[2], data[3]}
}
