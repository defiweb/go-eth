package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/txmodifier"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Load the private key.
	key, err := wallet.NewKeyFromJSON("./key.json", "test123")
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
		// uses it with SignTransaction, SendTransaction, and Sign methods
		// instead of relying on the node for signing.
		rpc.WithKeys(key),

		// Specify a default address for SendTransaction when the transaction
		// does not have a 'From' field set.
		rpc.WithDefaultAddress(key.Address()),

		// Specify a chain ID for SendTransaction when the transaction
		// does not have a 'ChainID' field set.
		rpc.WithChainID(1),

		// TX modifiers enable modifications to the transaction before signing
		// and sending to the node. While not mandatory, without them, transaction
		// parameters like gas limit, gas price, and nonce must be set manually.
		rpc.WithTXModifiers(
			// GasLimitEstimator automatically estimates the gas limit for the
			// transaction.
			txmodifier.NewGasLimitEstimator(txmodifier.GasLimitEstimatorOptions{
				Multiplier: 1.25,
			}),

			// GasFeeEstimator automatically estimates the gas price for the
			// transaction based on the current market conditions.
			txmodifier.NewEIP1559GasFeeEstimator(txmodifier.EIP1559GasFeeEstimatorOptions{
				GasPriceMultiplier:          1.25,
				PriorityFeePerGasMultiplier: 1.25,
			}),

			// NonceProvider automatically sets the nonce for the transaction.
			txmodifier.NewNonceProvider(txmodifier.NonceProviderOptions{
				UsePendingBlock: false,
			}),
		),
	)
	if err != nil {
		panic(err)
	}

	// Parse method signature.
	transfer := abi.MustParseMethod("transfer(address, uint256)(bool)")

	// Prepare a calldata for transfer call.
	calldata := transfer.MustEncodeArgs("0xd8da6bf26964af9d7eed9e03e53415d37aa96045", new(big.Int).Mul(big.NewInt(100), big.NewInt(1e6)))

	// Prepare a transaction.
	tx := types.NewTransaction().
		SetTo(types.MustAddressFromHex("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")).
		SetInput(calldata)

	txHash, _, err := c.SendTransaction(context.Background(), *tx)
	if err != nil {
		panic(err)
	}

	// Print the transaction hash.
	fmt.Printf("Transaction hash: %s\n", txHash.String())
}
