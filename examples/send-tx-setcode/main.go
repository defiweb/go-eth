package main

import (
	"context"
	"fmt"
	"os"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Load the private key.
	//
	// The key is both the authority (the EOA whose code is being set) and the
	// sponsor (the account paying for gas).
	k, err := wallet.NewKeyFromJSON(keyPath(), "test123")
	if err != nil {
		panic(err)
	}

	// Create transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://ethereum.publicnode.com"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(
		// Transport is always required.
		rpc.WithTransport(t),

		// Specify a key for signing transactions. If provided, the client
		// will sign transactions before sending them to the node.
		rpc.WithKeys(k),

		// Specify the default "from" address for transactions.
		rpc.WithDefaultAddress(rpc.AddressOptions{
			Address: k.Address(),
		}),

		// Estimate gas limit for transactions if not provided explicitly.
		rpc.WithGasLimit(rpc.GasLimitOptions{
			Multiplier: 1.25,
		}),

		// Estimate gas price for transactions if not provided explicitly.
		rpc.WithDynamicGasFee(rpc.DynamicGasFeeOptions{
			GasPriceMultiplier:          1.25,
			PriorityFeePerGasMultiplier: 1.25,
		}),

		// Automatically set the chain ID for transactions.
		rpc.WithChainID(rpc.ChainIDOptions{}),

		// Automatically set the nonce for transactions.
		rpc.WithNonce(rpc.NonceOptions{}),
	)
	if err != nil {
		panic(err)
	}

	// The contract whose code the EOA will adopt.
	contractAddr := types.MustAddressFromHex("0x0000000000000000000000000000000000000000")

	// Fetch the authority's current nonce. This is the nonce that will be
	// consumed by the authorization, not the transaction nonce.
	authorityNonce, err := c.GetTransactionCount(context.Background(), k.Address(), types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	// Build the authorization tuple.
	auth := types.Authorization{
		ChainID: 1,
		Address: contractAddr,
		Nonce:   authorityNonce,
	}

	// Sign the authorization.
	if err := auth.Sign(context.Background(), k); err != nil {
		panic(err)
	}

	// Prepare a set code transaction.
	tx := types.NewTransactionSetCode()
	tx.SetTo(k.Address())
	tx.AddAuthorization(auth)

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	// Print the transaction hash.
	//
	// After this transaction is mined, any call to k.Address() will execute
	// through the adopted contract code. The delegation persists until the EOA
	// sends another set code transaction with a different address or the zero
	// address to clear it.
	fmt.Printf("Transaction hash: %s\n", txHash.String())
}

func keyPath() string {
	if _, err := os.Stat("./key.json"); err == nil {
		return "./key.json"
	}
	return "./examples/send-tx-setcode/key.json"
}
