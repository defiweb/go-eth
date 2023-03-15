# go-eth

---

**This software is in alpha stage and is subject to change.**

---

The `go-eth` package is a suite of tools for interacting with Ethereum-based blockchains.

Some of key features include:

* An RPC client that supports HTTP, WebSocket and IPC transports.
* An ABI package allowing developers to easily interact with
  smart contracts.
* An extendable ABI encoder and decoder that allows user to easily interact with smart contracts.
* Support for JSON and HD wallets.

<!-- TOC -->

* [Installation](#installation)
* [Basic usage](#basic-usage)
    * [Connecting to a node](#connecting-to-a-node)
    * [Calling a contract method](#calling-a-contract-method)
    * [Sending a transaction](#sending-a-transaction)
* [Transports](#transports)
* [Wallets](#wallets)
* [Working with ABI](#working-with-abi)
    * [Methods](#methods)
        * [Encoding method arguments](#encoding-method-arguments)
        * [Decoding method arguments](#decoding-method-arguments)
    * [Events / Logs](#events--logs)
        * [Decoding events](#decoding-events)
    * [Errors](#errors)
    * [Reverts](#reverts)
    * [Panics](#panics)
    * [Contract ABI](#contract-abi)
        * [JSON-ABI](#json-abi)
        * [Human-Readable ABI](#human-readable-abi)
    * [Mapping rules](#mapping-rules)
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

## Basic usage

The examples below provide a glimpse into the usage of the `go-eth` package.

### Connecting to a node

The `go-eth` package provides an JSON-RPC client that can be used to connect to a node. In order to connect to a node,
you need to choose a transport and create a client. The following example shows how to connect to a node using the HTTP
transport:

```go
package main

import (
	"context"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
)

func main() {
	// Create a transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "http://example.com/rpc-node"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c := rpc.NewClient(t)

	// Get the latest block number.
	b, err := c.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	println(b)
}
```

### Calling a contract method

Calling a `balanceOf` method on a contract:

```go
package main

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func main() {
	// Create a transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://example.com/rpc-node"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c := rpc.NewClient(t)

	// Parse method signature.
	balanceOf := abi.MustParseMethod("balanceOf(address)(uint256)")

	// Prepare a calldata.
	calldata, err := balanceOf.EncodeArgs("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")
	if err != nil {
		panic(err)
	}

	// Call balanceOf.
	b, err := c.Call(context.Background(), types.Call{
		To:   types.MustHexToAddressPtr("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		Data: calldata,
	}, types.LatestBlockNumber)
	if err != nil {
		panic(err)
	}

	// Decode the result.
	var balance *big.Int
	err = balanceOf.DecodeValues(b, &balance)
	if err != nil {
		panic(err)
	}

	// Print the result.
	println(balance.String())
}
```

### Sending a transaction

Sending an ERC20 token transfer transaction:

```go
package main

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Load the private key.
	key, err := wallet.NewKeyFromJSON("./examples/keys/key.json", "test123")
	if err != nil {
		panic(err)
	}

	// Create a transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://example.com/rpc-node"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c, err := rpc.NewClient(
		// Transport is always required.
		rpc.WithTransport(t),

		// You can specify a key to sign transactions. If provided, the client will
		// use it to with SignTransaction, SendTransaction, and Sign methods instead
		// of making RPC calls.
		rpc.WithKeys(key),

		// You can specify a default address to use with SendTransaction if the
		// transaction doesn't have a "from" field set.
		rpc.WithDefaultAddress(key.Address()),

		// You can specify a chain ID to use with SendTransaction if the transaction
		// doesn't have a "chainID" field set.
		rpc.WithChainID(1),
	)
	if err != nil {
		panic(err)
	}

	transfer := abi.MustParseMethod("transfer(address, uint256)(bool)")

	// Prepare a calldata for transfer call.
	calldata, err := transfer.EncodeArgs("0xd8da6bf26964af9d7eed9e03e53415d37aa96045", new(big.Int).Mul(big.NewInt(100), big.NewInt(1e6)))
	if err != nil {
		panic(err)
	}

	// Prepare a transaction.
	tx := (&types.Transaction{}).
		SetType(types.DynamicFeeTxType).
		SetTo(types.MustAddressFromHex("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")).
		SetInput(calldata).
		SetNonce(0).
		SetMaxPriorityFeePerGas(big.NewInt(1 * 1e9)).
		SetMaxFeePerGas(big.NewInt(20 * 1e9))

	txHash, err := c.SendTransaction(context.Background(), *tx)
	if err != nil {
		panic(err)
	}

	// Print the transaction hash.
	println(txHash.String())
}
```

## Transports

To connect to a node, it is necessary to choose a suitable transport method. The transport is responsible for executing
a low-level communication protocol with the node. The `go-eth` package offers the following transport options:

| Transport | Description                                                                                 | Subscriptions |
|-----------|---------------------------------------------------------------------------------------------|---------------|
| HTTP      | Connects to a node using the HTTP protocol.                                                 | No            |
| WebSocket | Connects to a node using the WebSocket protocol.                                            | Yes           |
| IPC       | Connects to a node using the IPC protocol.                                                  | Yes           |
| Retry     | Wraps a transport and retries requests in case of an error.                                 | Yes           |
| Combined  | Wraps two transports and uses one for requests and the other for subscriptions.<sup>1</sup> | Yes           |

1. It is recommended by some RPC providers to use HTTP for requests and WebSocket for subscriptions.

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

```go
package main

import (
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
	println(key.Address().String())

}
```

## Working with ABI

The `abi` package is used for encoding and decoding ABI data. Internally each Solidity type is represented by the
two structures that implement the `abi.Type` and `abi.Value` interfaces. The `abi.Type` is used to represent a type of
Solidity variable, e.g. `uint256`, `address`, `bytes32`, etc. The `abi.Value` is used to represent a value of a Solidity
variable. It is similar to the `reflect.Type` and `reflect.Value` types in the standard library.

For example, the following code encodes an `uint256` value:

```go
package main

import (
	"github.com/defiweb/go-eth/abi"
)

func main() {
	u256Typ := abi.NewUintType(256)

	// Encode an uint256 value.
	u256ValEnc := u256Typ.Value().(*abi.UintValue)
	u256ValEnc.SetUint64(100)
	abiData, err := u256ValEnc.EncodeABI()
	if err != nil {
		panic(err)
	}

	// Decode an uint256 value.
	u256ValDec := u256Typ.Value().(*abi.UintValue)
	if _, err = u256ValDec.DecodeABI(abiData); err != nil {
		panic(err)
	}

	// Print the decoded value.
	println(u256ValDec.Uint64())
}

```

The example above gives an insight into the inner workings of the package, but this is not how the package is usually
used. Although this method is slightly faster, so it can be useful in some situations.

To make it easier to work with ABI data, the package provides a human-readable signature parser and a JSON ABI parser
(described later) to simplify creating types and a value mapper that helps to map ABI values to Go values.

The above example can be rewritten as follows:

```go
package main

import (
	"math/big"

	"github.com/defiweb/go-eth/abi"
)

func main() {
	// Create a new uint256 type using signature parser.
	u256Typ := abi.MustParseType("uint256")

	// Encode an uint256 value from an int type.
	abiData, err := abi.EncodeValue(u256Typ, 100)
	if err != nil {
		panic(err)
	}

	// Decode an uint256 value to big.Int.
	var u256Val big.Int
	if err = abi.DecodeValue(u256Typ, abiData, &u256Val); err != nil {
		panic(err)
	}

	// Print the decoded value.
	println(u256Val.Uint64())
}
```

In the example above, first the type is created using the `abi.MustParseType` function. Then the `abi.EncodeValue` and
`abi.DecodeValue` functions are used to encode and decode the value using a value mapper.

The `abi.MustParseType` could also parse a tuple type, e.g. `(uint256, address)`:

```go
package main

import (
	"github.com/defiweb/go-eth/abi"
)

type Data struct {
	Number  uint64 `abi:"num"`
	Address string `abi:"addr"`
}

func main() {
	// Create a new uint256 type using signature parser.
	tuple := abi.MustParseType("(uint256 num, address addr)")

	// Encode an uint256 value from an int type.
	abiData, err := abi.EncodeValue(tuple, &Data{
		Number:  100,
		Address: "0x1234567890123456789012345678901234567890",
	})
	if err != nil {
		panic(err)
	}

	// Decode an uint256 value to big.Int.
	var data Data
	if err = abi.DecodeValue(tuple, abiData, &data); err != nil {
		panic(err)
	}

	// Print the decoded value.
	println(data.Number)
	println(data.Address)
}
```

In the above example, the data is encoded and decoded using a struct. The struct `abi` tags are used to map the struct
fields to the tuple fields. These tags are optional; if they are not present, fields are mapped by their names with the
first consecutive uppercase letters being lowercased. For example, the `Number` field is mapped to the `number` field,
the `DAPPName` field is mapped to the `dappName` field, etc. If no names are specified for the tuple types, the default
names `arg0`, `arg1`, etc. are used.

Instead of using structs, it is also possible to encode and decode tuples to consecutive variables by using
the `abi.EncodeValues` and `abi.DecodeValues` functions (plural). Note that these plural versions of the encode/decode
functions can be used only with tuples.

```go
package main

import (
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/types"
)

func main() {
	// Create a new uint256 type using signature parser.
	tuple := abi.MustParseType("(uint256,address)")

	// Encode an uint256 value from an int type.
	abiData, err := abi.EncodeValues(tuple, 100, "0x1234567890123456789012345678901234567890")
	if err != nil {
		panic(err)
	}

	// Decode an uint256 value to big.Int.
	var u256Val big.Int
	var addrVal types.Address
	if err = abi.DecodeValues(tuple, abiData, &u256Val, &addrVal); err != nil {
		panic(err)
	}

	// Print the decoded value.
	println(u256Val.Uint64())
	println(addrVal.String())
}
```

### Methods

To work with methods, the `abi.Method` structure needs to be created. To create a method, the following methods can be
used:

- `abi.NewMethod(name, inputs, outputs)` - creates a new method with the given name, inputs and outputs types.

```go
package main

import "github.com/defiweb/go-eth/abi"

func main() {
	transfer := abi.NewMethod("transfer",
		abi.NewTupleType(
			abi.TupleTypeElem{Type: abi.NewAddressType()},
			abi.TupleTypeElem{Type: abi.NewUintType(256)},
		),
		abi.NewTupleType(
			abi.TupleTypeElem{Type: abi.NewBoolType()},
		),
	)
	// ...
}
```

- `abi.ParseMethod` / `abi.MustParseMethod` - creates a new method by parsing a method signature.

```go
package main

import "github.com/defiweb/go-eth/abi"

func main() {
	transfer := abi.MustParseMethod("transfer(address, uint256) returns (bool)")
	// ...
}
```

- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

The `abi.Method` structure allows to encode and decode method arguments and return values, calculate the method ID and
generate a method signature.

#### Encoding method arguments

To encode method arguments, the `abi.Method.EncodeArg` or `abi.Method.EncodeArgs` functions can be used. The first
function encodes a struct, the second function encodes consecutive variables.

```go
package main

import (
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

	// Prints: 0xa9059cbb00000000000000000000000012345678901234567890123456789012345678900000000000000000000000000000000000000000000000000000000000000064
	println(hexutil.BytesToHex(abiData))
}
```

#### Decoding method arguments

To decode method arguments, the `abi.Method.DecodeArg` or `abi.Method.DecodeArgs` functions can be used. The first
function decodes returned values to a struct, the second function decodes returned values to consecutive variables.

```go
package main

import (
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

func main() {
	abiData := hexutil.MustHexToBytes("0x0000000000000000000000000000000000000000000000002b5e3af16b1880000")

	// Parse method signature.
	balanceOf := abi.MustParseMethod("balanceOf(address) returns (uint256)")

	// Encode method arguments.
	var balance big.Int
	err := balanceOf.DecodeValues(abiData, &balance)
	if err != nil {
		panic(err)
	}

	// Prints: 195312500000000000
	println(balance.String())
}
```

### Events / Logs

To decode contract events, first a `abi.Event` struct must be created. Events may be created using different methods:

- `abi.NewEvent(name, inputs)` - creates a new event with the given name and inputs types.
- `abi.ParseEvent` / `abi.MustParseEvent` - creates a new event by parsing an event signature.
- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

#### Decoding events

```go
package main

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func main() {
	// Create a transport.
	t, err := transport.NewHTTP(transport.HTTPOptions{URL: "https://example.com/rpc-node"})
	if err != nil {
		panic(err)
	}

	// Create a JSON-RPC client.
	c := rpc.NewClient(t)

	transfer := abi.MustParseEvent("Transfer(address indexed src, address indexed dst, uint256 wad)")

	// Fetch logs for WETH transfer events.
	logs, err := c.GetLogs(context.Background(), types.FilterLogsQuery{
		Address:   []types.Address{types.MustAddressFromHex("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")},
		FromBlock: types.BlockNumberFromUint64Ptr(16492400),
		ToBlock:   types.BlockNumberFromUint64Ptr(16492400),
		Topics:    [][]types.Hash{{transfer.Topic0()}},
	})
	if err != nil {
		panic(err)
	}

	// Decode and print the logs.
	for _, log := range logs {
		var src, dst types.Address
		var wad *big.Int
		if err := transfer.DecodeValues(log.Topics, log.Data, &src, &dst, &wad); err != nil {
			panic(err)
		}
		println(src.String(), dst.String(), wad.String())
	}
}
```

### Errors

To decode contract errors, first a `abi.Error` struct must be created. Errors may be created using different methods:

- `abi.NewError(name, inputs)` - creates a new error with the given name and inputs types.

```go
package main

import "github.com/defiweb/go-eth/abi"

func main() {
	error := abi.NewError(
		"InsufficientBalance",
		abi.NewTupleType(
			abi.TupleTypeElem{Name: "available", Type: abi.NewUintType(256)},
			abi.TupleTypeElem{Name: "required", Type: abi.NewUintType(256)},
		),
	)
	// ...
}
```

- `abi.ParseError` / `abi.MustParseError` - creates a new error by parsing an error signature.

```go
package main

import "github.com/defiweb/go-eth/abi"

func main() {
	error := abi.MustParseError("InsufficientBalance(uint256 available, uint256 required)")
	// ...
}
```

- Using the `abi.Contract` struct (see [Contract ABI](#contract-abi) section).

### Reverts

Reverts are special errors that are returned by the EVM when a contract call fails. Reverts are ABI-encoded errors
with the `Error(string)` signature. To decode reverts, the `abi.DecodeRevert` function can be used. Optionally, the
`abi` package provides a `abi.Revert` that is a predefined error type that can be used to decode reverts.

To verify if an error is a revert, the `abi.IsRevert` function can be used.

### Panics

Similar to reverts, panics are special errors that are returned by the EVM when a contract call fails. Panics are
ABI-encoded errors with the `Panic(uint256)` signature. To decode panics, the `abi.DecodePanic` function can be used.
Optionally, the `abi` package provides a `abi.Panic` that is a predefined error type that can be used to decode panics.

To verify if an error is a panic, the `abi.IsPanic` function can be used.

### Contract ABI

The `abi.Contract` struct is a helper struct that provides an interface to a contract's ABI. It can be created using
a JSON-ABI file or by providing a list of signatures.

#### JSON-ABI

```go
package main

import (
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

	// ...
}
```

#### Human-Readable ABI

```go
package main

import (
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

	// ...
}
```

### Mapping rules

When mapping between Go and Solidity types, the following rules apply:

| Go type \ Solidity type | `intX`           | `uintX`            | `bool` | `string` | `bytes`       | `bytesX`         | `address`       |
|-------------------------|------------------|--------------------|--------|----------|---------------|------------------|-----------------|
| `intX`                  | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3</sup>    | ✗               |
| `uintX`                 | ✓<sup>1,2</sup>  | ✓<sup>1</sup>      | ✗      | ✗        | ✗             | ✓<sup>3</sup>    | ✗               |
| `bool`                  | ✗                | ✗                  | ✓      | ✗        | ✗             | ✗                | ✗               |
| `string`                | ✓<sup>5</sup>    | ✓<sup>5,6</sup>    | ✗      | ✓        | ✓<sup>7</sup> | ✓<sup>7,8</sup>  | ✓<sup>7,9</sup> |
| `[]byte`                | ✗                | ✗                  | ✗      | ✓        | ✓             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `[X]byte`               | ✗                | ✗                  | ✗      | ✗        | ✗             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `big.Int`               | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3</sup>    | ✗               |
| `types.Address`         | ✗                | ✗                  | ✗      | ✗        | ✓             | ✓<sup>4</sup>    | ✓               |
| `types.Hash`            | ✗                | ✗                  | ✗      | ✗        | ✓             | ✓<sup>3</sup>    | ✗               |
| `types.Bytes`           | ✗                | ✗                  | ✗      | ✓        | ✓             | ✓<sup>8</sup>    | ✓<sup>9</sup>   |
| `types.Number`          | ✓<sup>1</sup>    | ✓<sup>1,2</sup>    | ✗      | ✗        | ✗             | ✓<sup>3</sup>    | ✗               |
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

Note: `[X]byte` is a fixed-size byte array, e.g. `[20]byte`. `intX`, `uintX` and `bytesX` are fixed-size types,
e.g. `uint32`.

General rule for mapping rules is that the destination type must be able to hold the value of the source type,
conversion must be non-ambiguous, and mapping must be reversible. Mapping from larger to smaller types is supported
because very often Solidity contracts use `uint256` for all numbers, even if the value is known to be much less than
`2^256`.

### Signature parser syntax

The parser is based on the Solidity grammar, but allows to omit argument names, and the `returns` and `function`
keywords, so it can parse full Solidity signatures as well as short signatures like: `bar(uint256,bytes32)`.
Tuples are represented as a list of parameters, e.g. `(uint256,bytes32)`. The list can be optionally prefixed with
`tuple` keyword, e.g. `tuple(uint256,bytes32)`.

Examples of signatures that are accepted by the parser:

- `getPrice(string)`
- `getPrice(string)((uint256,unit256))`
- `getPrice(string symbol) returns ((uint256 price, unit256 timestamp) result)`
- `function getPrice(string calldata symbol) external view returns (tuple(uint256 price, uint256 timestamp) result)`
- `event PriceUpated(string indexed symbol, uint256 price)`
- `error PriceExpired(string symbol, uint256 timestamp)`

### Custom types

It is possible to add custom types to the `abi` package. Custom types are recognized by the signature parser.

#### Simple types

The simples way to create a custom type is to use `abi.ParseType` function that parses a type signature and returns
a `Type` struct. This method may be used to create custom types for commonly used structs, e.g.:

```go
package main

import (
	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

type Point struct {
	X int
	Y int
}

func main() {
	// Add custom type.
	abi.Default.Types["Point"] = abi.MustParseType("(int256 x, int256 y)")

	// Generate calldata.
	addTriangle := abi.MustParseMethod("addTriangle(Point a, Point b, Point c)")
	calldata, _ := addTriangle.EncodeArgs(
		Point{X: 1, Y: 2},
		Point{X: 3, Y: 4},
		Point{X: 5, Y: 6},
	)

	// Print the calldata.
	println(hexutil.BytesToHex(calldata))
}
```

#### Advanced types

More complex types can be created by implementing the `abi.Type` and `abi.Value` interfaces. The `abi.Type` interface
contains basic information about the type, and the `abi.Value` interface contains methods for encoding and decoding
values. It can optionally implement `abi.MapTo` and `abi.MapFrom` methods to support mapping to and from other types.

The following example shows how to create a custom type that represents a 32 byte bool array that is stored in a
single `bytes32` value:

```go
package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

type BoolFlagsType struct{}

func (b BoolFlagsType) CanonicalType() string {
	return "bytes32"
}

func (b BoolFlagsType) String() string {
	return "BoolFlags"
}

func (b BoolFlagsType) Value() abi.Value {
	return &BoolFlagsValue{}
}

type BoolFlagsValue [256]bool

func (b BoolFlagsValue) IsDynamic() bool {
	return false
}

func (b BoolFlagsValue) EncodeABI() (abi.Words, error) {
	var w abi.Word
	for i, v := range b {
		if v {
			w[i/8] |= 1 << uint(i%8)
		}
	}
	return abi.Words{w}, nil
}

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

func main() {
	// Add custom type.
	abi.Default.Types["BytesFlags"] = &BoolFlagsType{}

	// Generate calldata.
	setFlags := abi.MustParseMethod("setFlags(BytesFlags flags)")
	calldata, _ := setFlags.EncodeArgs(
		&BoolFlagsValue{true, false, true, true, false, true, false, true},
	)

	// Print the calldata.
	println(hexutil.BytesToHex(calldata))
}
```

## Additional tools

You may be also find the following tools interesting:

* [go-rlp](https://github.com/defiweb/go-rlp) - RLP serialization/deserialization library.
* [go-sigparser](https://github.com/defiweb/go-sigparser) - Solidity-compatible signature parser.
* [go-anymapper](https://github.com/defiweb/go-anymapper) - Data mapper used by this package.

## Documentation

[https://pkg.go.dev/github.com/defiweb/go-eth](https://pkg.go.dev/github.com/defiweb/go-eth)
