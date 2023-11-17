package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

type Point struct {
	X int
	Y int
}

func main() {
	// Add custom type.
	abi.Default.Types["Point"] = abi.MustParseStruct("struct {int256 x; int256 y;}")

	// Generate calldata.
	addTriangle := abi.MustParseMethod("addTriangle(Point a, Point b, Point c)")
	calldata := addTriangle.MustEncodeArgs(
		Point{X: 1, Y: 2},
		Point{X: 3, Y: 4},
		Point{X: 5, Y: 6},
	)

	// Print the calldata.
	fmt.Printf("Calldata: %s\n", hexutil.BytesToHex(calldata))
}
