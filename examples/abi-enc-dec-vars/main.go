package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

func main() {
	// Parse ABI type:
	dataABI := abi.MustParseStruct(`struct Data { int256 intVal; bool boolVal; string stringVal; }`)

	// Encode data:
	encodedData := abi.MustEncodeValues(dataABI, 42, true, "Hello, world!")

	// Print encoded data:
	fmt.Printf("Encoded data: %s\n", hexutil.BytesToHex(encodedData))

	// Decode data:
	var (
		intVal    int
		boolVal   bool
		stringVal string
	)
	abi.MustDecodeValues(dataABI, encodedData, &intVal, &boolVal, &stringVal)

	// Print decoded data:
	fmt.Printf("Decoded data: %d, %t, %s\n", intVal, boolVal, stringVal)
}
