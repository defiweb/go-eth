package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	erc20, err := abi.LoadJSON(abiPath())
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

func abiPath() string {
	if _, err := os.Stat("./erc20.json"); err == nil {
		return "./erc20.json"
	}
	return "./examples/contract-json-abi/erc20.json"
}
