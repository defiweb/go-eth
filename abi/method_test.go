package abi

import (
	"testing"

	"github.com/defiweb/go-eth/hexutil"
)

type Data struct {
	A []Data2
}

type Data2 struct {
	A string
	B int
}

func TestA(t *testing.T) {
	m, err := ParseMethod("transfer((address a,uint256 b)[] a)")
	if err != nil {
		t.Fatal(err)
	}

	abi, err := m.Encode(Data{A: []Data2{{A: "0x1234567890123456789012345678901234567890", B: 123}}})
	if err != nil {
		t.Fatal(err)
	}

	println(hexutil.BytesToHex(abi))
	println(m.Signature())

}
