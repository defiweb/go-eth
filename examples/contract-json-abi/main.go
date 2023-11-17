package main

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	erc20, err := abi.LoadJSON("erc20.json")
	if err != nil {
		panic(err)
	}

	transfer := erc20.Methods["transfer"]
	calldata, err := transfer.EncodeArgs(
		"0x1234567890123456789012345678901234567890",
		big.NewInt(1e18),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer calldata: 0x%x\n", calldata)
}
