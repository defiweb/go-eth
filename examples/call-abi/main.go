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
