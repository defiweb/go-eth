package main

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
)

func main() {
	// Create transport.
	//
	// There are several other transports available:
	// - HTTP (NewHTTP)
	// - WebSocket (NewWebsocket)
	// - IPC (NewIPC)
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://ethereum.publicnode.com"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(rpc.WithTransport(t))
	if err != nil {
		panic(err)
	}

	// Get the latest block number.
	b, err := c.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Latest block number:", b)
}
