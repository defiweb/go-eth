package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func main() {
	// Create transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://ethereum.publicnode.com"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(rpc.WithTransport(t))
	if err != nil {
		panic(err)
	}

	// Parse method signature.
	balanceOf := abi.MustParseMethod("balanceOf(address)(uint256)")

	// Prepare a calldata.
	calldata := balanceOf.MustEncodeArgs("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")

	// Prepare a call.
	call := types.NewCall().
		SetTo(types.MustAddressFromHex("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")).
		SetInput(calldata)

	// Call balanceOf.
	b, _, err := c.Call(context.Background(), *call, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	// Decode the result.
	var balance *big.Int
	balanceOf.MustDecodeValues(b, &balance)

	// Print the result.
	fmt.Printf("Balance: %s\n", balance.String())
}
