package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer ctxCancel()

	// Create transport.
	t, err := transport.NewWebsocket(transport.WebsocketOptions{
		Context: ctx,
		URL:     "wss://ethereum.publicnode.com",
	})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(rpc.WithTransport(t))
	if err != nil {
		panic(err)
	}

	// Parse event signature.
	transfer := abi.MustParseEvent("event Transfer(address indexed src, address indexed dst, uint256 wad)")

	// Create a filter query.
	query := types.NewFilterLogsQuery().
		SetAddresses(types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")).
		SetTopics([]types.Hash{transfer.Topic0()})

	// Fetch logs for WETH transfer events.
	logs, err := c.SubscribeLogs(ctx, query)
	if err != nil {
		panic(err)
	}

	// Decode and print events.
	for log := range logs {
		var (
			src types.Address
			dst types.Address
			wad *big.Int
		)
		transfer.MustDecodeValues(log.Topics, log.Data, &src, &dst, &wad)
		fmt.Printf("Transfer: %s -> %s: %s\n", src.String(), dst.String(), wad.String())
	}
}
