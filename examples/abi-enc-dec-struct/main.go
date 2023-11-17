package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

// Data is a struct that represents the data we want to encode and decode.
type Data struct {
	IntValue    int    `abi:"intVal"`
	BoolValue   bool   `abi:"boolVal"`
	StringValue string `abi:"stringVal"`
}

func main() {
	// Parse ABI type:
	dataABI := abi.MustParseStruct(`struct Data { int256 intVal; bool boolVal; string stringVal; }`)

	// Encode data:
	encodedData := abi.MustEncodeValue(dataABI, Data{
		IntValue:    42,
		BoolValue:   true,
		StringValue: "Hello, world!",
	})

	// Print encoded data:
	fmt.Printf("Encoded data: %s\n", hexutil.BytesToHex(encodedData))

	// Decode data:
	var decodedData Data
	abi.MustDecodeValue(dataABI, encodedData, &decodedData)

	// Print decoded data:
	fmt.Printf("Decoded data: %+v\n", decodedData)
}
