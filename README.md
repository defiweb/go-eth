[![Run Tests](https://github.com/defiweb/go-eth/actions/workflows/test.yml/badge.svg)](https://github.com/defiweb/go-eth/actions/workflows/test.yml)

# go-eth

This library is a Go package designed to interact with the Ethereum blockchain. This package provides robust tools for
connecting to Ethereum nodes, sending transactions, and handling smart contract events. Whether you're developing a
decentralized application or conducting blockchain analysis.

Some of key features include:

* An RPC client that supports HTTP, WebSocket and IPC transports.
* An ABI package allowing developers to easily interact with smart contracts.
* An extendable and easy to use ABI encoder and decoder.
* Support for JSON and HD wallets.

<!-- TOC -->

* [go-eth](#go-eth)
    * [Installation](#installation)
    * [Quick start](#quick-start)
        * [Connecting to a node](#connecting-to-a-node)
        * [Calling a contract method](#calling-a-contract-method)
        * [Calling a contract method using a Human-Readable ABI](#calling-a-contract-method-using-a-human-readable-abi)
        * [Sending a transaction](#sending-a-transaction)
        * [Subscribing to events](#subscribing-to-events)
    * [Transports](#transports)
    * [Wallets](#wallets)
    * [Working with ABI](#working-with-abi)
        * [Mapping rules](#mapping-rules)
        * [Encoding and Decoding Methods](#encoding-and-decoding-methods)
            * [Encoding method arguments](#encoding-method-arguments)
            * [Decoding method return values](#decoding-method-return-values)
        * [Events / Logs](#events--logs)
            * [Decoding events](#decoding-events)
        * [Errors](#errors)
        * [Reverts](#reverts)
        * [Panics](#panics)
        * [Contract ABI](#contract-abi)
            * [JSON-ABI](#json-abi)
            * [Human-Readable ABI](#human-readable-abi)
        * [Signature parser syntax](#signature-parser-syntax)
        * [Custom types](#custom-types)
            * [Simple types](#simple-types)
            * [Advanced types](#advanced-types)
    * [Additional tools](#additional-tools)
    * [Documentation](#documentation)

<!-- TOC -->

## Installation

```bash
go get -u github.com/defiweb/go-eth
```

## Quick start

The examples below provide a glimpse into the usage of the `go-eth` package.

### Connecting to a node

The `go-eth` package offers a JSON-RPC client that can be used to establish a connection with a node. The example below
demonstrates how to connect to a node using HTTP transport method.

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

The example demonstrates how to call the `balanceOf` method on a contract.

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
```

### Calling a contract method using a Human-Readable ABI

Following example shows how to call a contract method using a Human-Readable ABI. It uses popular
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
	call := types.NewCall().
		SetTo(types.MustAddressFromHex("0xcA11bde05977b3631167028862bE2a173976CA11")).
		SetInput(calldata)

	// Call the contract.
	b, _, err := c.Call(context.Background(), *call, types.LatestBlockNumber)
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

The following example demonstrates how to execute an ERC20 token transfer transaction. Additionally, it illustrates the
use of TX modifiers to simplify the transaction creation process.

<!-- examples/send-tx/main.go -->

```go
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
```

### Subscribing to events

Following example shows how to subscribe to WETH transfer events.

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
	query := types.NewFilterLogsQuery().
		SetAddresses(types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")).
		SetTopics([]types.Hash{transfer.Topic0()})

	// Fetch logs for WETH transfer events.
	logs, err := c.SubscribeLogs(ctx, *query)
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

To connect to a node, it is necessary to choose a suitable transport method. The transport is responsible for executing
a low-level communication protocol with the node. The `go-eth` package offers the following transport options:

| Transport | Description                                                                                 | Subscriptions   |
|-----------|---------------------------------------------------------------------------------------------|-----------------|
| HTTP      | Connects to a node using the HTTP protocol.                                                 | No              |
| WebSocket | Connects to a node using the WebSocket protocol.                                            | Yes             |
| IPC       | Connects to a node using the IPC protocol.                                                  | Yes             |
| Retry     | Wraps a transport and retries requests in case of an error.                                 | Yes<sup>2</sup> |
| Combined  | Wraps two transports and uses one for requests and the other for subscriptions.<sup>1</sup> | Yes             |

1. It is recommended by some RPC providers to use HTTP for requests and WebSocket for subscriptions.
2. Only if the underlying transport supports subscriptions.

Transports can be created using the `transport.New*` functions. It is also possible to create custom transport by
implementing the `transport.Transport` interface or `transport.SubscriptionTransport` interface.

## Wallets

The `go-eth` package provides support for the following wallet types:

| Description                  | Example                                                                     |
|------------------------------|-----------------------------------------------------------------------------|
| A random key                 | `key := wallet.NewRandomKey()`                                              |
| Private key                  | `key, err := wallet.NewKeyFromBytes(privateKey)`                            |
| JSON key file<sup>1</sup>    | `key, err := wallet.NewKeyFromJSON(path, password)`                         |
| JSON key content<sup>1</sup> | `key, err := wallet.NewKeyFromJSONContent(jsonContent, password)`           |
| Mnemonic                     | `key, err := wallet.NewKeyFromMnemonic(mnemonic, password, account, index)` |

1. Only V3 JSON keys are supported.

Wallets can be also created using custom derivation paths. For example, the following code creates a wallet using the
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

## Working with ABI

The `go-eth` package offers an ABI encoder and decoder for working with ABI data. The package also includes a signature
parser for parsing method, event, and error signatures, as well as custom types and structs.

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
consecutive uppercase letters converted to lowercase. For instance, the `Number` struct field maps to the `number`
field, and the `DAPPName` field maps to the `dappName` field.

It is also possible to encode and decode values to a separate variables:

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

Note that in both examples above, similarly named functions are used to encode and decode data. The only difference is
that the second example uses the plural form of the function. The plural form is used to encode and decode data from
separate variables, while the singular form is used for structs or maps. This is a common pattern in the `go-eth`
package.

Finally, instead of using signature parser, it is possible to create types manually which may be useful to create
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
|-------------------------|------------------|--------------------|--------|----------|---------------|------------------|-----------------|
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

* ✓ - Supported
* ✗ - Not supported

1. Destination type must be able to hold the value of the source type. For example, `uint16` can be mapped to `uint8`,
   but only if the value is less than 256.
2. Mapping of negative values is supported only if both types support negative values.
3. Only mapping from/to `bytes32` is supported.
4. Only mapping from/to `bytes20` is supported.
5. String representation of the number is assumed to be in hexadecimal format. When string is used as a source value,
   the "0x" prefix is optional. Negative values are prefixed with a minus sign, e.g. "-0x123".
6. Negative values are not supported.
7. String representation is assumed to be in hexadecimal format.
8. When mapping to `bytesX`, length of the data must the same as the length of the destination type.
9. When mapping to `address`, length of the data must be 20 bytes.
10. Mapping latest, earliest and pending block numbers is not supported.

Note: Go type `[X]byte` represents a fixed-size byte array, such as `[20]byte`. Solidity types `intX`, `uintX`,
and `bytesX` are also fixed-size types, such as, `uint32`.

The general rule for mapping types is that the destination type must be capable of holding the value of the source type,
the conversion must be unambiguous, and the mapping must be reversible. Mapping from larger to smaller types is
supported because often Solidity contracts use `uint256` for all numbers, even when the value is known to be much less
than 256 bits.

### Encoding and Decoding Methods

To work with methods, the `abi.Method` structure needs to be created. Methods may be created using different methods:

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

To decode method arguments, the `abi.Method.DecodeArg` or `abi.Method.DecodeArgs` functions can be used. The first
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

	// Encode method arguments.
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

To decode contract events, the `abi.Event` structure needs to be created. Events may be created using different methods:

- `abi.ParseEvent` / `abi.MustParseEvent` - creates a new event by parsing an event signature.
-
    - `abi.NewEvent(name, inputs)` - creates a new event using provided arguments.
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
	query := types.NewFilterLogsQuery().
		SetAddresses(types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")).
		SetFromBlock(types.BlockNumberFromUint64Ptr(16492400)).
		SetToBlock(types.BlockNumberFromUint64Ptr(16492400)).
		SetTopics([]types.Hash{transfer.Topic0()})

	// Fetch logs for WETH transfer events.
	logs, err := c.GetLogs(context.Background(), *query)
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

### Errors

To decode custom contract errors, first a `abi.Error` struct must be created. Errors may be created using different
methods:

- `abi.ParseError` / `abi.MustParseError` - creates a new error by parsing an error signature.
- `abi.NewError(name, inputs)` - creates a new error using provided arguments.
- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

Custom errors may be decoded from errors returned by the `Call` function using the `abi.Error.HandleError` method.

When using a `abi.Contract`, errors may be decoded from call errors using the `abi.Contract.HandleError` method. This
method will try to decode the error using all errors defined in the contract, also including reverts and panics.

### Reverts

Reverts are special errors returned by the EVM when a contract call fails. Reverts are ABI-encoded errors with
the `Error(string)` signature. The `abi.DecodeRevert` function can be used to decode reverts. Optionally, the `abi`
package provides `abi.Revert`, a predefined error type that can be used to decode reverts.

To verify if an error is a revert, use the `abi.IsRevert` function.

### Panics

Similar to reverts, panics are special errors returned by the EVM when a contract call fails. Panics are ABI-encoded
errors with the `Panic(uint256)` signature. The `abi.DecodePanic` function can be used to decode panics. Optionally, the
`abi` package also provides `abi.Panic`, a predefined error type that can be used to decode panics.

To verify if an error is a panic, use the `abi.IsPanic` function.

### Contract ABI

The `abi.Contract` structure is a utility that provides an interface to a contract. It can be created using a JSON-ABI
file or by supplying a list of signatures (also known as a Human-Readable ABI).

#### JSON-ABI

<!-- examples/contract-json-abi/main.go -->

```go
package main

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	erc20, err := abi.LoadJSON("erc20.json")
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

### Signature parser syntax

The parser is based on Solidity grammar, but it allows for the omission of argument names, as well as the `returns`
and `function` keywords. This means it can parse full Solidity signatures as well as short signatures, such
as `bar(uint256,bytes32)`. Tuples are represented as a list of parameters, for example, `(uint256,bytes32)`. The list
can be optionally prefixed with the `tuple` keyword, for example, `tuple(uint256,bytes32)`.

Examples of signatures that are accepted by the parser:

- `getPrice(string)`
- `getPrice(string)((uint256,unit256))`
- `getPrice(string symbol) returns ((uint256 price, unit256 timestamp) result)`
- `function getPrice(string calldata symbol) external view returns (tuple(uint256 price, uint256 timestamp) result)`
- `event PriceUpated(string indexed symbol, uint256 price)`
- `error PriceExpired(string symbol, uint256 timestamp)`

### Custom types

It is possible to add custom types to the `abi` package.

#### Simple types

The simplest way to create a custom type is to use the `abi.ParseType`, `abi.ParseStruct`, `abi.MustParseType`,
`abi.MustParseStruct` functions, which parses a type signature and returns a `Type` struct. This method can be used to
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

The example below demonstrates how to create a custom type that represents a 32-byte boolean array stored in a
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
		return 0, fmt.Errorf("abi: cannot decode BytesFlags from empty data")
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

// MapFrom maps value from a different type.
func (b *BoolFlagsValue) MapFrom(_ abi.Mapper, src any) error {
	switch src := src.(type) {
	case [256]bool:
		*b = src
	case []bool:
		if len(src) > 256 {
			return fmt.Errorf("abi: cannot map []bool of length %d to BytesFlags", len(src))
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

## Additional tools

You may be also find the following tools interesting:

* [go-rlp](https://github.com/defiweb/go-rlp) - RLP serialization/deserialization library.
* [go-sigparser](https://github.com/defiweb/go-sigparser) - Solidity-compatible signature parser.
* [go-anymapper](https://github.com/defiweb/go-anymapper) - Data mapper used by this package.

## Documentation

[https://pkg.go.dev/github.com/defiweb/go-eth](https://pkg.go.dev/github.com/defiweb/go-eth)
