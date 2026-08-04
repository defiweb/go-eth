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

	transfer := abi.MustParseEvent("Transfer(address indexed src, address indexed dst, uint256 wad)")

	// Create filter query.
	query := types.NewFilterLogsQuery()
	query.SetAddresses(types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"))
	query.SetFromBlock(types.BlockNumberFromUint64Ptr(16492400))
	query.SetToBlock(types.BlockNumberFromUint64Ptr(16492400))
	query.SetTopics([]types.Hash{transfer.Topic0()})

	// Fetch logs for WETH transfer events.
	logs, err := c.GetLogs(context.Background(), query)
	if err != nil {
		panic(err)
	}

	// Decode and print events.
	for _, log := range logs {
		var src, dst types.Address
		var wad *big.Int
		transfer.MustDecodeValues(log.Topics, log.Data, &src, &dst, &wad)
		fmt.Printf("Transfer: %s -> %s: %s\n", src.String(), dst.String(), wad.String())
	}
}
