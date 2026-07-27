[![Run Tests](https://github.com/defiweb/go-eth/actions/workflows/test.yml/badge.svg)](https://github.com/defiweb/go-eth/actions/workflows/test.yml)

# go-eth

This library is a comprehensive Go package designed to interact with the Ethereum blockchain. It provides robust tools for
connecting to Ethereum nodes, sending transactions, and handling smart contract events. Whether you're developing a
decentralized application or conducting blockchain analysis, this package offers everything you need.

Key features include:

- **RPC Client**: Supports HTTP, WebSocket, and IPC transports with advanced middleware capabilities
- **ABI Package**: Comprehensive ABI encoder/decoder with Human-Readable ABI support
- **Transaction Support**: All transaction types including Legacy, EIP-1559, EIP-2930, and EIP-4844 (blob transactions)
- **Wallet Management**: Support for JSON keystore, HD wallets, and mnemonic phrases
- **Cryptographic Functions**: Built-in support for ECDSA signing, Keccak hashing, and KZG4844 proofs
- **Extensible Design**: Modular architecture with customizable transport hijackers and middleware

<!-- TOC -->

- [go-eth](#go-eth)
    - [Installation](#installation)
    - [Quick start](#quick-start)
        - [Connecting to a node](#connecting-to-a-node)
        - [Calling a contract method](#calling-a-contract-method)
        - [Calling a contract method using a Human-Readable ABI](#calling-a-contract-method-using-a-human-readable-abi)
        - [Sending a transaction](#sending-a-transaction)
        - [Subscribing to events](#subscribing-to-events)
    - [Transports](#transports)
    - [Wallets](#wallets)
    - [Working with ABI](#working-with-abi)
        - [Mapping rules](#mapping-rules)
        - [Encoding and Decoding Methods](#encoding-and-decoding-methods)
            - [Encoding method arguments](#encoding-method-arguments)
            - [Decoding method return values](#decoding-method-return-values)
        - [Events / Logs](#events--logs)
            - [Decoding events](#decoding-events)
        - [Contract ABI](#contract-abi)
            - [JSON-ABI](#json-abi)
            - [Human-Readable ABI](#human-readable-abi)
        - [Errors](#errors)
        - [Reverts](#reverts)
        - [Panics](#panics)
        - [Signature parser syntax](#signature-parser-syntax)
        - [Custom types](#custom-types)
            - [Simple types](#simple-types)
            - [Advanced types](#advanced-types)
    - [Client Configuration](#client-configuration)
        - [Available Options](#available-options)
        - [Transport Hijacking](#transport-hijacking)
    - [Cryptographic Functions](#cryptographic-functions)
    - [Utility Packages](#utility-packages)
    - [Additional tools](#additional-tools)
    - [Documentation](#documentation)

<!-- TOC -->

## Installation

```bash
go get -u github.com/defiweb/go-eth
```

## Quick start

The examples below provide a comprehensive overview of the `go-eth` package capabilities.

### Connecting to a node

The `go-eth` package offers a JSON-RPC client that can connect to Ethereum nodes using various transport methods. The example below
demonstrates how to connect using HTTP transport.

<!-- examples/connect/main.go -->

```go
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
```

### Calling a contract method

The following example demonstrates how to call the `balanceOf` method on an ERC-20 contract.

<!-- examples/call/main.go -->

```go
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
	call := types.NewCall()
	call.SetTo(types.MustAddressFromHex("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"))
	call.SetInput(calldata)

	// Call balanceOf.
	b, err := c.Call(context.Background(), call, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	// Decode the result.
	var balance *big.Int
	balanceOf.MustDecodeValues(b, &balance)

	// Print the result.
	fmt.Printf("Balance: %s\n", balance.String())
}
```

### Calling a contract method using a Human-Readable ABI

The following example shows how to call a contract method using a Human-Readable ABI. It uses the popular
[Multicall3](https://www.multicall3.com) contract as an example.

<!-- examples/call-abi/main.go -->

```go
package main

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type Call3 struct {
	Target       types.Address `abi:"target"`
	AllowFailure bool          `abi:"allowFailure"`
	CallData     []byte        `abi:"callData"`
}

type Result struct {
	Success    bool   `abi:"success"`
	ReturnData []byte `abi:"returnData"`
}

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

	// Parse contract ABI.
	multicall := abi.MustParseSignatures(
		"struct Call { address target; bytes callData; }",
		"struct Call3 { address target; bool allowFailure; bytes callData; }",
		"struct Call3Value { address target; bool allowFailure; uint256 value; bytes callData; }",
		"struct Result { bool success; bytes returnData; }",
		"function aggregate(Call[] calldata calls) public payable returns (uint256 blockNumber, bytes[] memory returnData)",
		"function aggregate3(Call3[] calldata calls) public payable returns (Result[] memory returnData)",
		"function aggregate3Value(Call3Value[] calldata calls) public payable returns (Result[] memory returnData)",
		"function blockAndAggregate(Call[] calldata calls) public payable returns (uint256 blockNumber, bytes32 blockHash, Result[] memory returnData)",
		"function getBasefee() view returns (uint256 basefee)",
		"function getBlockHash(uint256 blockNumber) view returns (bytes32 blockHash)",
		"function getBlockNumber() view returns (uint256 blockNumber)",
		"function getChainId() view returns (uint256 chainid)",
		"function getCurrentBlockCoinbase() view returns (address coinbase)",
		"function getCurrentBlockDifficulty() view returns (uint256 difficulty)",
		"function getCurrentBlockGasLimit() view returns (uint256 gaslimit)",
		"function getCurrentBlockTimestamp() view returns (uint256 timestamp)",
		"function getEthBalance(address addr) view returns (uint256 balance)",
		"function getLastBlockHash() view returns (bytes32 blockHash)",
		"function tryAggregate(bool requireSuccess, Call[] calldata calls) public payable returns (Result[] memory returnData)",
		"function tryBlockAndAggregate(bool requireSuccess, Call[] calldata calls) public payable returns (uint256 blockNumber, bytes32 blockHash, Result[] memory returnData)",
	)

	// Prepare a calldata.
	// In this example we will call the `getCurrentBlockGasLimit` and `getCurrentBlockTimestamp` methods
	// on the Multicall3 contract.
	calldata := multicall.Methods["aggregate3"].MustEncodeArgs([]Call3{
		{
			Target:   types.MustAddressFromHex("0xcA11bde05977b3631167028862bE2a173976CA11"),
			CallData: multicall.Methods["getCurrentBlockGasLimit"].MustEncodeArgs(),
		},
		{
			Target:   types.MustAddressFromHex("0xcA11bde05977b3631167028862bE2a173976CA11"),
			CallData: multicall.Methods["getCurrentBlockTimestamp"].MustEncodeArgs(),
		},
	})

	// Prepare a call.
	call := types.NewCall()
	call.SetTo(types.MustAddressFromHex("0xcA11bde05977b3631167028862bE2a173976CA11"))
	call.SetInput(calldata)

	// Call the contract.
	b, err := c.Call(context.Background(), call, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	// Decode the result.
	var (
		results   []Result
		gasLimit  uint64
		timestamp uint64
	)
	multicall.Methods["aggregate3"].MustDecodeValues(b, &results)
	multicall.Methods["getCurrentBlockGasLimit"].MustDecodeValues(results[0].ReturnData, &gasLimit)
	multicall.Methods["getCurrentBlockTimestamp"].MustDecodeValues(results[1].ReturnData, &timestamp)

	// Print the result.
	fmt.Println("Gas limit:", gasLimit)
	fmt.Println("Timestamp:", timestamp)
}
```

### Sending a transaction

The following example demonstrates how to execute an ERC-20 token transfer transaction. It also illustrates the
use of client options to automatically set gas, nonce, and other transaction parameters.

<!-- examples/send-tx/main.go -->

```go
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
```

### Subscribing to events

The following example shows how to subscribe to WETH transfer events.

<!-- examples/subscription/main.go -->

```go
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
	query := types.NewFilterLogsQuery()
	query.SetAddresses(types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"))
	query.SetTopics([]types.Hash{transfer.Topic0()})

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
```

## Transports

To connect to a node, you need to choose an appropriate transport method. The transport handles low-level communication
with the node. The `go-eth` package offers the following transport options:

| Transport | Description                                                                                | Subscriptions   |
| --------- | ------------------------------------------------------------------------------------------ | --------------- |
| HTTP      | Connects to a node using the HTTP protocol. Most widely supported.                         | No              |
| WebSocket | Connects to a node using the WebSocket protocol. Supports real-time subscriptions.         | Yes             |
| IPC       | Connects to a node using Inter-Process Communication (Unix sockets/named pipes).           | Yes             |
| Retry     | Wraps a transport and retries requests in case of errors with exponential backoff.         | Yes<sup>2</sup> |
| Combined  | Wraps two transports and uses one for methods and the other for subscriptions.<sup>1</sup> | Yes             |

1. Recommended by some RPC providers to use HTTP for methods and WebSocket for subscriptions for better performance.
2. Only if the underlying transport supports subscriptions.

Transports can be created using the `transport.New*` functions. You can also create custom transports by
implementing the `transport.Transport` interface or the `transport.SubscriptionTransport` interface for subscription support.

## Wallets

The `go-eth` package provides comprehensive support for various wallet types and key management:

| Description                       | Example                                                                     |
| --------------------------------- | --------------------------------------------------------------------------- |
| Random key generation             | `key := wallet.NewRandomKey()`                                              |
| Private key from bytes            | `key, err := wallet.NewKeyFromBytes(privateKey)`                            |
| JSON keystore file<sup>1</sup>    | `key, err := wallet.NewKeyFromJSON(path, password)`                         |
| JSON keystore content<sup>1</sup> | `key, err := wallet.NewKeyFromJSONContent(jsonContent, password)`           |
| HD wallet from mnemonic           | `key, err := wallet.NewKeyFromMnemonic(mnemonic, password, account, index)` |
| Remote RPC wallet                 | `key := wallet.NewKeyRPC(client, address)`                                  |

1. Only Ethereum JSON V3 keystores are supported.

You can also create wallets using custom derivation paths. For example, the following code creates a wallet using the
`m/44'/60'/0'/10/10` derivation path:

<!-- examples/key-mnemonic/main.go -->

```go
package main

import (
	"fmt"

	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Parse mnemonic.
	mnemonic, err := wallet.NewMnemonic("gravity trophy shrimp suspect sheriff avocado label trust dove tragic pitch title network myself spell task protect smooth sword diary brain blossom under bulb", "")
	if err != nil {
		panic(err)
	}

	// Parse derivation path.
	path, err := wallet.ParseDerivationPath("m/44'/60'/0'/10/10")
	if err != nil {
		panic(err)
	}

	// Derive private key.
	key, err := mnemonic.Derive(path)
	if err != nil {
		panic(err)
	}

	// Print the address of the derived private key.
	fmt.Println("Private key:", key.Address().String())
}
```

## Client Configuration

The RPC client can be configured with various options to automatically handle transaction parameters, implement middleware, and customize behavior.

### Available Options

| Option                   | Description                               |
| ------------------------ | ----------------------------------------- |
| `WithTransport`          | Sets the transport for communication      |
| `WithKeys`               | Adds private keys for transaction signing |
| `WithDefaultAddress`     | Sets default sender address               |
| `WithChainID`            | Auto-sets chain ID for transactions       |
| `WithNonce`              | Auto-manages transaction nonces           |
| `WithGasLimit`           | Auto-estimates gas limits                 |
| `WithLegacyGasFee`       | Auto-estimates legacy gas prices          |
| `WithDynamicGasFee`      | Auto-estimates EIP-1559 gas fees          |
| `WithSimulate`           | Simulates transactions before sending     |
| `WithTransactionDecoder` | Custom transaction decoder                |

## Working with ABI

The `go-eth` package offers a comprehensive ABI encoder and decoder for working with smart contract data. The package includes a signature
parser for parsing method, event, and error signatures, as well as support for custom types and structs.

The following example shows how to encode and decode data:

<!-- examples/abi-enc-dec-struct/main.go -->

```go
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
```

In the example above, data is encoded and decoded using a struct. The `abi` tags map the struct fields to the
corresponding tuple or struct fields. These tags are optional. If absent, fields are mapped by name, with the first
consecutive uppercase letters converted to lowercase.

It is also possible to encode and decode values to separate variables:

<!-- examples/abi-enc-dec-vars/main.go -->

```go
package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

func main() {
	// Parse ABI type:
	dataABI := abi.MustParseStruct(`struct Data { int256 intVal; bool boolVal; string stringVal; }`)

	// Encode data:
	encodedData := abi.MustEncodeValues(dataABI, 42, true, "Hello, world!")

	// Print encoded data:
	fmt.Printf("Encoded data: %s\n", hexutil.BytesToHex(encodedData))

	// Decode data:
	var (
		intVal    int
		boolVal   bool
		stringVal string
	)
	abi.MustDecodeValues(dataABI, encodedData, &intVal, &boolVal, &stringVal)

	// Print decoded data:
	fmt.Printf("Decoded data: %d, %t, %s\n", intVal, boolVal, stringVal)
}
```

**Note: In both examples above, similarly named functions are used to encode and decode data. The only difference is
that the second example uses the plural form of the function names.** The plural form encodes/decodes data from
separate variables, while the singular form works with structs or maps. This is a consistent pattern throughout the `go-eth`
package.

Finally, instead of using the signature parser, you can create types manually, which may be useful for creating
custom types programmatically:

<!-- examples/abi-enc-dec-prog/main.go -->

```go
package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

func main() {
	// Create ABI type:
	dataABI := abi.NewTupleType(
		abi.TupleTypeElem{
			Name: "intVal",
			Type: abi.NewIntType(256),
		},
		abi.TupleTypeElem{
			Name: "boolVal",
			Type: abi.NewBoolType(),
		},
		abi.TupleTypeElem{
			Name: "stringVal",
			Type: abi.NewStringType(),
		},
	)

	// Encode data:
	encodedData := abi.MustEncodeValues(dataABI, 42, true, "Hello, world!")

	// Print encoded data:
	fmt.Printf("Encoded data: %s\n", hexutil.BytesToHex(encodedData))

	// Decode data:
	var (
		intVal    int
		boolVal   bool
		stringVal string
	)
	abi.MustDecodeValues(dataABI, encodedData, &intVal, &boolVal, &stringVal)

	// Print decoded data:
	fmt.Printf("Decoded data: %d, %t, %s\n", intVal, boolVal, stringVal)
}
```

### Mapping rules

When mapping between Go and Solidity types, the following rules apply:

| Go type \ Solidity type | `intX`           | `uintX`            | `bool` | `string` | `bytes`       | `bytesX`         | `address`       |
| ----------------------- | ---------------- | ------------------ | ------ | -------- | ------------- | ---------------- | --------------- |
| `intX`                  | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3,6</sup>  | ✗               |
| `uintX`                 | ✓<sup>1,2</sup>  | ✓<sup>1</sup>      | ✗      | ✗        | ✗             | ✓<sup>3,6</sup>  | ✗               |
| `bool`                  | ✗                | ✗                  | ✓      | ✗        | ✗             | ✗                | ✗               |
| `string`                | ✓<sup>5</sup>    | ✓<sup>5,6</sup>    | ✗      | ✓        | ✓<sup>7</sup> | ✓<sup>7,8</sup>  | ✓<sup>7,9</sup> |
| `[]byte`                | ✗                | ✗                  | ✗      | ✓        | ✓             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `[X]byte`               | ✗                | ✗                  | ✗      | ✗        | ✗             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `big.Int`               | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3,6</sup>  | ✗               |
| `types.Address`         | ✗                | ✗                  | ✗      | ✗        | ✓             | ✓<sup>4</sup>    | ✓               |
| `types.Hash`            | ✗                | ✗                  | ✗      | ✗        | ✓             | ✓<sup>3</sup>    | ✗               |
| `types.Bytes`           | ✗                | ✗                  | ✗      | ✓        | ✓             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `types.Number`          | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3,6</sup>  | ✗               |
| `types.BlockNumber`     | ✓<sup>1,10</sup> | ✓<sup>1,2,10</sup> | ✗      | ✗        | ✗             | ✓<sup>3,10</sup> | ✗               |

- ✓ - Supported
- ✗ - Not supported

1. Destination type must be able to hold the value of the source type. Otherwise, the mapping will result in an error.
   For example, `uint16` can be mapped to `uint8`, but only if the value is less than 256.
2. Mapping of negative values is supported only if both types support negative values.
3. Only mapping from/to `bytes32` is supported.
4. Only mapping from/to `bytes20` is supported.
5. String representation of the number is assumed to be in hexadecimal format. When string is used as a source value,
   the "0x" prefix is optional. Negative values are prefixed with a minus sign, e.g. "-0x123".
6. Negative values are not supported.
7. String representation is assumed to be in hexadecimal format.
8. When mapping to `bytesX`, the length of the data must be the same as the length of the destination type.
9. When mapping to `address`, length of the data must be 20 bytes.
10. Mapping of latest, earliest, and pending block numbers is not supported.

Note: The Go type `[X]byte` represents a fixed-size byte array, such as `[20]byte`. Solidity types `intX`, `uintX`,
and `bytesX` are also fixed-size types, such as `uint32`.

The general rule for type mapping is that the destination type must be capable of holding the value of the source type,
the conversion must be unambiguous, and the mapping must be reversible. Mapping from larger to smaller types is
supported because Solidity contracts often use `uint256` for all numbers, even when the value is known to be much smaller
than 256 bits.

### Encoding and Decoding Methods

To work with methods, an `abi.Method` structure needs to be created. Methods can be created using different approaches:

- `abi.ParseMethod` / `abi.MustParseMethod` - creates a new method by parsing a method signature.
- `abi.NewMethod(name, inputs, outputs, mutability)` - creates a new method using provided arguments.
- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

#### Encoding method arguments

To encode method arguments, the `abi.Method.EncodeArg` or `abi.Method.EncodeArgs` functions can be used. The first
function encodes a struct, the second function encodes consecutive variables.

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func main() {
	// Parse method signature.
	transfer := abi.MustParseMethod("transfer(address, uint256) returns (bool)")

	// Encode method arguments.
	abiData, err := transfer.EncodeArgs(
		types.MustAddressFromHex("0x1234567890123456789012345678901234567890"),
		big.NewInt(100),
	)
	if err != nil {
		panic(err)
	}

	// Print encoded data.
	fmt.Printf("Encoded data: %s\n", hexutil.BytesToHex(abiData))
}
```

#### Decoding method return values

To decode method return values, the `abi.Method.DecodeValue` or `abi.Method.DecodeValues` functions can be used. The first
function decodes returned values to a struct, the second function decodes returned values to consecutive variables.

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

func main() {
	abiData := hexutil.MustHexToBytes("0x00000000000000000000000000000000000000000000000002b5e3af16b1880000")

	// Parse method signature.
	balanceOf := abi.MustParseMethod("balanceOf(address) returns (uint256)")

	// Decode method return values.
	var balance big.Int
	err := balanceOf.DecodeValues(abiData, &balance)
	if err != nil {
		panic(err)
	}

	// Prints decoded data.
	fmt.Printf("Balance: %s\n", balance.String())
}
```

### Events / Logs

To decode contract events, an `abi.Event` structure needs to be created. Events can be created using different approaches:

- `abi.ParseEvent` / `abi.MustParseEvent` - creates a new event by parsing an event signature.
- `abi.NewEvent(name, inputs, anonymous)` - creates a new event using provided arguments.
- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

#### Decoding events

<!-- examples/events/main.go -->

```go
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
```

### Contract ABI

The `abi.Contract` structure is a utility that provides an interface to a contract. It can be created using a JSON-ABI
file or by supplying a list of signatures (also known as a Human-Readable ABI).

To create a contract struct, the following methods can be used:

- `abi.LoadJSON` / `abi.MustLoadJSON` - creates a new contract by loading a JSON-ABI file.
- `abi.ParseJSON` / `abi.MustParseJSON` - creates a new contract by parsing a JSON-ABI string.
- `abi.ParseSignatures` / `abi.MustParseSignatures` - creates a new contract by parsing a list of signatures (
  Human-Readable ABI).

#### JSON-ABI

<!-- examples/contract-json-abi/main.go -->

```go
package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	erc20, err := abi.LoadJSON(abiPath())
	if err != nil {
		panic(err)
	}

	transfer := erc20.Methods["transfer"]
	calldata, err := transfer.EncodeArgs(
		"0x1234567890123456789012345678901234567890",
		big.NewInt(1e18),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer calldata: 0x%x\n", calldata)
}

func abiPath() string {
	if _, err := os.Stat("./erc20.json"); err == nil {
		return "./erc20.json"
	}
	return "./examples/contract-json-abi/erc20.json"
}
```

#### Human-Readable ABI

<!-- examples/contract-hra-abi/main.go -->

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	erc20, err := abi.ParseSignatures(
		"function name() public view returns (string)",
		"function symbol() public view returns (string)",
		"function decimals() public view returns (uint8)",
		"function totalSupply() public view returns (uint256)",
		"function balanceOf(address _owner) public view returns (uint256 balance)",
		"function transfer(address _to, uint256 _value) public returns (bool success)",
		"function transferFrom(address _from, address _to, uint256 _value) public returns (bool success)",
		"function approve(address _spender, uint256 _value) public returns (bool success)",
		"function allowance(address _owner, address _spender) public view returns (uint256 remaining)",
		"event Transfer(address indexed _from, address indexed _to, uint256 _value)",
		"event Approval(address indexed _owner, address indexed _spender, uint256 _value)",
	)
	if err != nil {
		panic(err)
	}

	transfer := erc20.Methods["transfer"]
	calldata, err := transfer.EncodeArgs(
		"0x1234567890123456789012345678901234567890",
		big.NewInt(1e18),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer calldata: 0x%x\n", calldata)
}
```

### Errors

To decode custom contract errors, first an `abi.Error` struct must be created. Errors can be created using different approaches:

- `abi.ParseError` / `abi.MustParseError` - creates a new error by parsing an error signature.
- `abi.NewError(name, inputs)` - creates a new error using provided arguments.
- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

Custom errors may be decoded from errors returned by the `Call` function using the `abi.Error.HandleError` method.

When using an `abi.Contract`, errors can be decoded from call errors using the `abi.Contract.HandleError` method. This
method will attempt to decode the error using all errors defined in the contract, including reverts and panics.

### Reverts

Reverts are special errors returned by the EVM when a contract call fails. Reverts are ABI-encoded errors with
the `Error(string)` signature. The `abi.DecodeRevert` function can be used to decode reverts. The `abi`
package also provides `abi.Revert`, a predefined error type for decoding reverts.

To verify if an error is a revert, use the `abi.IsRevert` function.

### Panics

Similar to reverts, panics are special errors returned by the EVM when a contract call fails. Panics are ABI-encoded
errors with the `Panic(uint256)` signature. The `abi.DecodePanic` function can be used to decode panics. The
`abi` package also provides `abi.Panic`, a predefined error type for decoding panics.

To verify if an error is a panic, use the `abi.IsPanic` function.

### Signature parser syntax

The parser is based on Solidity grammar but allows for the omission of argument names, as well as the `returns`
and `function` keywords. This means it can parse full Solidity signatures as well as short signatures, such
as `bar(uint256,bytes32)`. Tuples are represented as a list of parameters, for example, `(uint256,bytes32)`. The list
can optionally be prefixed with the `tuple` keyword, for example, `tuple(uint256,bytes32)`.

Examples of signatures accepted by the parser:

- `getPrice(string)`
- `getPrice(string)((uint256,uint256))`
- `getPrice(string symbol) returns ((uint256 price, uint256 timestamp) result)`
- `function getPrice(string calldata symbol) external view returns (tuple(uint256 price, uint256 timestamp) result)`
- `event PriceUpdated(string indexed symbol, uint256 price)`
- `error PriceExpired(string symbol, uint256 timestamp)`

### Custom types

The `go-eth` package allows for the creation of custom types that can be used with the ABI encoder and decoder and with
the signature parser.

#### Simple types

The simplest way to create a custom type is to use the `abi.ParseType`, `abi.ParseStruct`, `abi.MustParseType`, or
`abi.MustParseStruct` functions, which parse a type signature and return a `Type` struct. This method can be used to
create custom types for commonly used structs.

<!-- examples/custom-type-simple/main.go -->

```go
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
```

#### Advanced types

More complex types can be created by implementing the `abi.Type` and `abi.Value` interfaces. The `abi.Type` interface
provides basic information about the type, while the `abi.Value` interface includes methods for encoding and decoding
values and holds the value itself. Optionally, the `abi.MapTo` and `abi.MapFrom` methods can be implemented to support
advanced mapping logic.

The example below demonstrates how to create a custom type that represents a 256-bit boolean array stored in a
single `bytes32` value.

<!-- examples/custom-type-advanced/main.go -->

```go
package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

// BoolFlagsType is a custom type that represents a 256-bit bitfield.
//
// It must implement the abi.Type interface.
type BoolFlagsType struct{}

// IsDynamic returns true if the type is dynamic-length, like string or bytes.
func (b BoolFlagsType) IsDynamic() bool {
	return false
}

// CanonicalType is the type as it would appear in the ABI.
// It must only use the types defined in the ABI specification:
// https://docs.soliditylang.org/en/latest/abi-spec.html
func (b BoolFlagsType) CanonicalType() string {
	return "bytes32"
}

// String returns the custom type name.
func (b BoolFlagsType) String() string {
	return "BoolFlags"
}

// Value returns the zero value for this type.
func (b BoolFlagsType) Value() abi.Value {
	return &BoolFlagsValue{}
}

// BoolFlagsValue is the value of the custom type.
//
// It must implement the abi.Value interface.
type BoolFlagsValue [256]bool

// IsDynamic returns true if the type is dynamic-length, like string or bytes.
func (b BoolFlagsValue) IsDynamic() bool {
	return false
}

// EncodeABI encodes the value to the ABI format.
func (b BoolFlagsValue) EncodeABI() (abi.Words, error) {
	var w abi.Word
	for i, v := range b {
		if v {
			w[i/8] |= 1 << uint(i%8)
		}
	}
	return abi.Words{w}, nil
}

// DecodeABI decodes the value from the ABI format.
func (b *BoolFlagsValue) DecodeABI(words abi.Words) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("abi: cannot decode BoolFlags from empty data")
	}
	for i, v := range words[0] {
		for j := 0; j < 8; j++ {
			b[i*8+j] = v&(1<<uint(j)) != 0
		}
	}
	return 1, nil
}

// MapFrom and MapTo are optional methods that allow mapping between different
// types.
//
// The abi.Mapper is the instance of the internal mapper that is used to
// perform the mapping. It can be used to map nested types.
//
// Note that you may want to use reflection to implement the following methods,
// as it would let you write more generic mapping code.

// MapFrom maps value from a different type.
func (b *BoolFlagsValue) MapFrom(_ abi.Mapper, src any) error {
	switch src := src.(type) {
	case [256]bool:
		*b = src
	case []bool:
		if len(src) > 256 {
			return fmt.Errorf("abi: cannot map []bool of length %d to BoolFlags", len(src))
		}
		for i, v := range src {
			b[i] = v
		}
	}
	return nil
}

// MapTo maps value to a different type.
func (b *BoolFlagsValue) MapTo(_ abi.Mapper, dst any) error {
	switch dst := dst.(type) {
	case *[256]bool:
		*dst = *b
	case *[]bool:
		*dst = make([]bool, 256)
		for i, v := range b {
			(*dst)[i] = v
		}
	}
	return nil
}

func main() {
	// Add custom type.
	abi.Default.Types["BoolFlags"] = &BoolFlagsType{}

	// Generate calldata.
	setFlags := abi.MustParseMethod("setFlags(BoolFlags flags)")
	calldata, _ := setFlags.EncodeArgs(
		[]bool{true, false, true, true, false, true, false, true},
	)

	// Print the calldata.
	fmt.Printf("Calldata: %s\n", hexutil.BytesToHex(calldata))
}
```

Please note that adding a custom type to the `abi.Default.Types` map will affect all instances of the `abi` package in
the current process. If you want to add a custom type to a single `abi` instance, you can create a new instance using
the `abi.NewABI` function.

## Cryptographic Functions

The `go-eth` package includes a comprehensive crypto module with support for various cryptographic operations:

### ECDSA Operations

- **Key Recovery**: `crypto.ECRecoverHash()` and `crypto.ECRecoverMessage()`
- **Digital Signatures**: `crypto.ECSignHash()` and `crypto.ECSignMessage()`
- **Address Generation**: `crypto.ECPublicKeyToAddress()`

### Hashing

- **Keccak-256**: `crypto.Keccak256()` - The primary hash function used in Ethereum

### KZG4844 (EIP-4844 Blob Transactions)

The package includes full support for KZG commitments and proofs used in blob transactions:

- **Blob to Commitment**: `crypto.KZGBlobToCommitment()`
- **Proof Generation**: `crypto.KZGComputeProof()` and `crypto.KZGComputeBlobProof()`
- **Proof Verification**: `crypto.KZGVerifyProof()` and `crypto.KZGVerifyBlobProof()`
- **Blob Hash Computation**: `crypto.KZGComputeBlobHashV1()`

Example KZG usage:

```go
import "github.com/defiweb/go-eth/crypto/kzg4844"

// Create a blob (128KB of data)
var blob kzg4844.Blob
copy(blob[:], "your data here...")

// Generate commitment
commitment, err := kzg4844.BlobToCommitment(&blob)
if err != nil {
    panic(err)
}

// Generate proof for a specific point
var point kzg4844.Point
proof, claim, err := kzg4844.ComputeProof(&blob, point)
```

### Transaction Signing

The crypto module also provides transaction signing capabilities through the `txsign` subpackage, supporting all Ethereum transaction types including legacy, EIP-1559, EIP-2930, and EIP-4844.

## Utility Packages

### HexUtil Package

The `hexutil` package provides convenient functions for working with hexadecimal data:

```go
import "github.com/defiweb/go-eth/hexutil"

// Convert between big integers and hex strings
hexStr := hexutil.BigIntToHex(big.NewInt(12345))  // "0x3039"
bigInt := hexutil.MustHexToBigInt("0x3039")       // 12345

// Convert between byte slices and hex strings
hexStr = hexutil.BytesToHex([]byte{0x12, 0x34})  // "0x1234"
bytes := hexutil.MustHexToBytes("0x1234")        // []byte{0x12, 0x34}

// Utility functions
hasPrefix := hexutil.Has0xPrefix("0x1234")       // true
```

### Types Package

The `types` package provides Ethereum-specific data types:

- **Address**: 20-byte Ethereum addresses with validation
- **Hash**: 32-byte hashes with utilities
- **BlockNumber**: Special handling for latest/earliest/pending blocks
- **Transaction Types**: Support for all transaction types (Legacy, AccessList, DynamicFee, Blob)
- **Call Types**: Various call structures for different use cases
- **Filter Queries**: For event log filtering

Example:

```go
import "github.com/defiweb/go-eth/types"

// Create address from hex string
addr := types.MustAddressFromHex("0x742d35Cc6606e7d8123456789012345678901234")

// Create various transaction types
legacyTx := types.NewTransactionLegacy()
eip1559Tx := types.NewTransactionDynamicFee()
eip2930Tx := types.NewTransactionAccessList()
blobTx := types.NewTransactionBlob()

// Create call for contract interaction
call := types.NewCall()
call.SetTo(addr)
call.SetInput(calldata)
```

## Additional tools

You may also find the following related tools interesting:

- [go-rlp](https://github.com/defiweb/go-rlp) - RLP serialization/deserialization library.
- [go-sigparser](https://github.com/defiweb/go-sigparser) - Solidity-compatible signature parser.
- [go-anymapper](https://github.com/defiweb/go-anymapper) - Data mapper used by this package.

## Documentation

[https://pkg.go.dev/github.com/defiweb/go-eth](https://pkg.go.dev/github.com/defiweb/go-eth)
