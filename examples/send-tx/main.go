package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Load the private key.
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

		// Simulate transactions before sending them to the node.
		rpc.WithSimulate(),
	)
	if err != nil {
		panic(err)
	}

	// Parse method signature.
	transfer := abi.MustParseMethod("transfer(address, uint256)(bool)")

	// Prepare a calldata for transfer call.
	calldata := transfer.MustEncodeArgs("0xd8da6bf26964af9d7eed9e03e53415d37aa96045", new(big.Int).Mul(big.NewInt(100), big.NewInt(1e6)))

	// Prepare a transaction.
	tx := types.NewTransactionLegacy()
	tx.SetTo(types.MustAddressFromHex("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"))
	tx.SetInput(calldata)

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	// Print the transaction hash.
	fmt.Printf("Transaction hash: %s\n", txHash.String())
}

func keyPath() string {
	if _, err := os.Stat("./key.json"); err == nil {
		return "./key.json"
	}
	return "./examples/send-tx/key.json"
}
