package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/hexutil"
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
	)
	if err != nil {
		panic(err)
	}

	// Prepare the blobs.
	//
	// A blob is a fixed 128 KiB, so it is allocated on the heap rather than
	// as a local variable.
	//
	// The data is used verbatim. EIP-4844 requires every 32-byte field
	// element of a blob to be smaller than the BLS12-381 modulus, so
	// arbitrary bytes must be encoded before being placed in a blob.
	var blobs []types.BlobInfo
	for _, data := range []string{"hello world 1", "hello world 2"} {
		blob := new(crypto.KZGBlob)
		copy(blob[:], data)

		// NewBlobInfo computes the KZG commitment, the KZG proof, and the
		// versioned hash for the blob.
		info, err := types.NewBlobInfo(blob)
		if err != nil {
			panic(err)
		}
		blobs = append(blobs, info)
	}

	// Prepare a transaction.
	tx := types.NewTransactionBlob()
	tx.SetTo(types.MustAddressFromHex("0x69B352cbE6Fc5C130b6F62cc8f30b9d7B0DC27d0"))
	tx.SetBlobs(blobs)

	// The blob gas price is not estimated by any of the client options above,
	// so it has to be set explicitly.
	tx.SetMaxFeePerBlobGas(big.NewInt(1e10))

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		panic(err)
	}

	// Print the transaction hash and the versioned hash of each blob. The
	// versioned hashes are what the EVM sees; the blobs themselves are not
	// accessible to contracts and are discarded by the network after a few
	// weeks.
	fmt.Printf("Transaction hash: %s\n", txHash.String())
	for i, b := range blobs {
		fmt.Printf("Blob %d versioned hash: %s\n", i, hexutil.BytesToHex(b.Hash[:]))
	}
}

func keyPath() string {
	if _, err := os.Stat("./key.json"); err == nil {
		return "./key.json"
	}
	return "./examples/send-tx-blob/key.json"
}
