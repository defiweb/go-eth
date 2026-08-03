package types

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionOnChain_JSON(t *testing.T) {
	tests := []struct {
		tx       *TransactionOnChain
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "to": "0x2222222222222222222222222222222222222222",
				  "gas": "0x186a0",
				  "gasPrice": "0x3b9aca00",
				  "input": "0x01020304",
				  "nonce": "0x1",
				  "value": "0xde0b6b3a7640000",
				  "v": "0x6f",
				  "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
				  "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
				  "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
				  "blockHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
				  "blockNumber": "0x3",
				  "transactionIndex": "0x4"
				}
			`,
			tx: &TransactionOnChain{
				Transaction: &TransactionLegacy{
					SigningData: SigningData{
						Nonce:     ptr(uint64(1)),
						Signature: MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"),
					},
					CallLegacy: CallLegacy{
						ExecutionData: ExecutionData{
							To:       MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
							Value:    big.NewInt(1000000000000000000),
							GasLimit: ptr(uint64(100000)),
							Input:    []byte{1, 2, 3, 4},
						},
						LegacyFeeData: LegacyFeeData{
							GasPrice: big.NewInt(1000000000),
						},
					},
				},
				Hash:             MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				BlockHash:        MustHashFromHexPtr("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
				BlockNumber:      big.NewInt(3),
				TransactionIndex: ptr(uint64(4)),
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			tx := &TransactionOnChain{}

			err := tx.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.tx, tx)

			j, err := tx.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestTransactionReceipt_JSON(t *testing.T) {
	tests := []struct {
		receipt  *TransactionReceipt
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "transactionHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
				  "transactionIndex": "0x4",
				  "blockHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
				  "blockNumber": "0x3",
				  "from": "0x3333333333333333333333333333333333333333",
				  "to": "0x4444444444444444444444444444444444444444",
				  "cumulativeGasUsed": "0x186a0",
				  "effectiveGasPrice": "0x3b9aca00",
				  "gasUsed": "0x5208",
				  "contractAddress": "0x5555555555555555555555555555555555555555",
				  "logs": [],
				  "logsBloom": "0x01020304",
				  "root": "0x6666666666666666666666666666666666666666666666666666666666666666",
				  "status": "0x1"
				}
			`,
			receipt: &TransactionReceipt{
				TransactionHash:   MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				TransactionIndex:  4,
				BlockHash:         MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
				BlockNumber:       big.NewInt(3),
				From:              MustAddressFromHex("0x3333333333333333333333333333333333333333"),
				To:                MustAddressFromHex("0x4444444444444444444444444444444444444444"),
				CumulativeGasUsed: 100000,
				EffectiveGasPrice: big.NewInt(1000000000),
				GasUsed:           21000,
				ContractAddress:   MustAddressFromHexPtr("0x5555555555555555555555555555555555555555"),
				Logs:              []Log{},
				LogsBloom:         []byte{1, 2, 3, 4},
				Root:              MustHashFromHexPtr("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
				Status:            ptr(uint64(1)),
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			receipt := &TransactionReceipt{}

			err := receipt.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.receipt, receipt)

			j, err := receipt.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestBlock_JSON(t *testing.T) {
	tests := []struct {
		block    *Block
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "number": "0x1",
				  "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
				  "parentHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
				  "stateRoot": "0x3333333333333333333333333333333333333333333333333333333333333333",
				  "receiptsRoot": "0x4444444444444444444444444444444444444444444444444444444444444444",
				  "transactionsRoot": "0x5555555555555555555555555555555555555555555555555555555555555555",
				  "mixHash": "0x6666666666666666666666666666666666666666666666666666666666666666",
				  "sha3Uncles": "0x7777777777777777777777777777777777777777777777777777777777777777",
				  "nonce": "0x0000000000000042",
				  "miner": "0x8888888888888888888888888888888888888888",
				  "logsBloom": "` + testBloomHex + `",
				  "difficulty": "0x2",
				  "totalDifficulty": "0x3",
				  "size": "0x220",
				  "gasLimit": "0x1c9c380",
				  "gasUsed": "0x5208",
				  "timestamp": "0x64000000",
				  "uncles": ["0x9999999999999999999999999999999999999999999999999999999999999999"],
				  "extraData": "0x01020304",
				  "transactions": ["0xaaaa000000000000000000000000000000000000000000000000000000000000"]
				}
			`,
			block: &Block{
				Number:           big.NewInt(1),
				Hash:             MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				ParentHash:       MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
				StateRoot:        MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", PadNone),
				ReceiptsRoot:     MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
				TransactionsRoot: MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
				MixHash:          MustHashFromHex("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
				Sha3Uncles:       MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", PadNone),
				Nonce:            big.NewInt(0x42),
				Miner:            MustAddressFromHex("0x8888888888888888888888888888888888888888"),
				LogsBloom:        testBloom(),
				Difficulty:       big.NewInt(2),
				TotalDifficulty:  big.NewInt(3),
				Size:             544,
				GasLimit:         30000000,
				GasUsed:          21000,
				Timestamp:        time.Unix(0x64000000, 0),
				Uncles:           []Hash{MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", PadNone)},
				TransactionHashes: []Hash{
					MustHashFromHex("0xaaaa000000000000000000000000000000000000000000000000000000000000", PadNone),
				},
				ExtraData: []byte{1, 2, 3, 4},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			block := &Block{}

			err := block.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.block, block)

			j, err := block.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestBlock_JSON_FullTransactions(t *testing.T) {
	wantJSON := `
		{
		  "number": "0x1",
		  "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
		  "parentHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
		  "stateRoot": "0x3333333333333333333333333333333333333333333333333333333333333333",
		  "receiptsRoot": "0x4444444444444444444444444444444444444444444444444444444444444444",
		  "transactionsRoot": "0x5555555555555555555555555555555555555555555555555555555555555555",
		  "mixHash": "0x6666666666666666666666666666666666666666666666666666666666666666",
		  "sha3Uncles": "0x7777777777777777777777777777777777777777777777777777777777777777",
		  "nonce": "0x0000000000000042",
		  "miner": "0x8888888888888888888888888888888888888888",
		  "logsBloom": "` + testBloomHex + `",
		  "difficulty": "0x2",
		  "totalDifficulty": "0x3",
		  "size": "0x220",
		  "gasLimit": "0x1c9c380",
		  "gasUsed": "0x5208",
		  "timestamp": "0x64000000",
		  "uncles": [],
		  "extraData": "0x01020304",
		  "transactions": [
		    {
		      "to": "0x2222222222222222222222222222222222222222",
		      "gas": "0x186a0",
		      "gasPrice": "0x3b9aca00",
		      "input": "0x01020304",
		      "nonce": "0x1",
		      "value": "0xde0b6b3a7640000",
		      "v": "0x6f",
		      "r": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad91490",
		      "s": "0x8051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd84",
		      "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
		      "blockHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
		      "blockNumber": "0x3",
		      "transactionIndex": "0x4"
		    }
		  ]
		}
	`
	block := &Block{}

	err := block.UnmarshalJSON([]byte(wantJSON))
	require.NoError(t, err)
	require.Len(t, block.Transactions, 1)
	assert.Empty(t, block.TransactionHashes)

	j, err := block.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, wantJSON, string(j))
}

func TestFeeHistory_JSON(t *testing.T) {
	tests := []struct {
		feeHistory *FeeHistory
		wantJSON   string
	}{
		{
			wantJSON: `
				{
				  "oldestBlock": "0x10",
				  "reward": [["0x1", "0x2"], ["0x3", "0x4"]],
				  "baseFeePerGas": ["0x5", "0x6", "0x7"],
				  "gasUsedRatio": [0.25, 0.5],
				  "baseFeePerBlobGas": ["0x8", "0x9", "0xa"],
				  "blobGasUsedRatio": [0.75, 1]
				}
			`,
			feeHistory: &FeeHistory{
				OldestBlock: 16,
				Reward: [][]*big.Int{
					{big.NewInt(1), big.NewInt(2)},
					{big.NewInt(3), big.NewInt(4)},
				},
				BaseFeePerGas:     []*big.Int{big.NewInt(5), big.NewInt(6), big.NewInt(7)},
				GasUsedRatio:      []float64{0.25, 0.5},
				BaseFeePerBlobGas: []*big.Int{big.NewInt(8), big.NewInt(9), big.NewInt(10)},
				BlobGasUsedRatio:  []float64{0.75, 1},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			feeHistory := &FeeHistory{}

			err := feeHistory.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.feeHistory, feeHistory)

			j, err := feeHistory.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestAccessListResult_JSON(t *testing.T) {
	tests := []struct {
		result   *AccessListResult
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "accessList": [
				    {
				      "address": "0x1111111111111111111111111111111111111111",
				      "storageKeys": [
				        "0x2222222222222222222222222222222222222222222222222222222222222222"
				      ]
				    }
				  ],
				  "gasUsed": "0x5208",
				  "error": "execution reverted"
				}
			`,
			result: &AccessListResult{
				AccessList: AccessList{
					{
						Address: MustAddressFromHex("0x1111111111111111111111111111111111111111"),
						StorageKeys: []Hash{
							MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
						},
					},
				},
				GasUsed: 21000,
				Error:   "execution reverted",
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			result := &AccessListResult{}

			err := result.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.result, result)

			j, err := result.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestStorageProof_JSON(t *testing.T) {
	tests := []struct {
		proof    *StorageProof
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "key": "0x1111111111111111111111111111111111111111111111111111111111111111",
				  "value": "0x2a",
				  "proof": ["0x01020304", "0x05060708"]
				}
			`,
			proof: &StorageProof{
				Key:   MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", PadNone),
				Value: big.NewInt(42),
				Proof: []Bytes{{1, 2, 3, 4}, {5, 6, 7, 8}},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			proof := &StorageProof{}

			err := proof.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.proof, proof)

			j, err := proof.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestAccountProof_JSON(t *testing.T) {
	tests := []struct {
		proof    *AccountProof
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "address": "0x1111111111111111111111111111111111111111",
				  "accountProof": ["0x01020304"],
				  "balance": "0xde0b6b3a7640000",
				  "codeHash": "0x2222222222222222222222222222222222222222222222222222222222222222",
				  "nonce": "0x7",
				  "storageHash": "0x3333333333333333333333333333333333333333333333333333333333333333",
				  "storageProof": [
				    {
				      "key": "0x4444444444444444444444444444444444444444444444444444444444444444",
				      "value": "0x2a",
				      "proof": ["0x05060708"]
				    }
				  ]
				}
			`,
			proof: &AccountProof{
				Address:      MustAddressFromHex("0x1111111111111111111111111111111111111111"),
				AccountProof: []Bytes{{1, 2, 3, 4}},
				Balance:      big.NewInt(1000000000000000000),
				CodeHash:     MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
				Nonce:        7,
				StorageHash:  MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", PadNone),
				StorageProof: []StorageProof{
					{
						Key:   MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
						Value: big.NewInt(42),
						Proof: []Bytes{{5, 6, 7, 8}},
					},
				},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			proof := &AccountProof{}

			err := proof.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.proof, proof)

			j, err := proof.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestLog_JSON(t *testing.T) {
	tests := []struct {
		log      *Log
		wantJSON string
	}{
		{
			wantJSON: `
				{
				  "address": "0x1111111111111111111111111111111111111111",
				  "topics": [
				    "0x2222222222222222222222222222222222222222222222222222222222222222",
				    "0x3333333333333333333333333333333333333333333333333333333333333333"
				  ],
				  "data": "0x01020304",
				  "blockHash": "0x4444444444444444444444444444444444444444444444444444444444444444",
				  "blockNumber": "0x5",
				  "transactionHash": "0x6666666666666666666666666666666666666666666666666666666666666666",
				  "transactionIndex": "0x7",
				  "logIndex": "0x8",
				  "removed": true
				}
			`,
			log: &Log{
				Address: MustAddressFromHex("0x1111111111111111111111111111111111111111"),
				Topics: []Hash{
					MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone),
					MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", PadNone),
				},
				Data:             []byte{1, 2, 3, 4},
				BlockHash:        MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
				BlockNumber:      big.NewInt(5),
				TransactionHash:  MustHashFromHexPtr("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
				TransactionIndex: ptr(uint64(7)),
				LogIndex:         ptr(uint64(8)),
				Removed:          true,
			},
		},
		{
			// A pending log: every nullable field is absent.
			wantJSON: `
				{
				  "address": "0x1111111111111111111111111111111111111111",
				  "topics": [],
				  "data": "0x",
				  "blockHash": null,
				  "blockNumber": null,
				  "transactionHash": null,
				  "transactionIndex": null,
				  "logIndex": null,
				  "removed": false
				}
			`,
			log: &Log{
				Address: MustAddressFromHex("0x1111111111111111111111111111111111111111"),
				Topics:  []Hash{},
				Data:    []byte{},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			log := &Log{}

			err := log.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.log, log)

			j, err := log.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestFilterLogsQuery_JSON(t *testing.T) {
	tests := []struct {
		query    *FilterLogsQuery
		wantJSON string
	}{
		{
			// oneOrList collapses a single-element list to a bare value on
			// output, so wantJSON is written in that canonical form. See
			// TestFilterLogsQuery_JSON_OneOrList for the equivalence of the
			// two input spellings.
			wantJSON: `
				{
				  "address": [
				    "0x1111111111111111111111111111111111111111",
				    "0x2222222222222222222222222222222222222222"
				  ],
				  "fromBlock": "0x1",
				  "toBlock": "latest",
				  "topics": [
				    "0x3333333333333333333333333333333333333333333333333333333333333333",
				    [
				      "0x4444444444444444444444444444444444444444444444444444444444444444",
				      "0x5555555555555555555555555555555555555555555555555555555555555555"
				    ]
				  ],
				  "blockhash": "0x6666666666666666666666666666666666666666666666666666666666666666"
				}
			`,
			query: &FilterLogsQuery{
				Address: []Address{
					MustAddressFromHex("0x1111111111111111111111111111111111111111"),
					MustAddressFromHex("0x2222222222222222222222222222222222222222"),
				},
				FromBlock: ptr(BlockNumberFromUint64(1)),
				ToBlock:   ptr(LatestBlockNumber),
				Topics: [][]Hash{
					{MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", PadNone)},
					{
						MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", PadNone),
						MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", PadNone),
					},
				},
				BlockHash: MustHashFromHexPtr("0x6666666666666666666666666666666666666666666666666666666666666666", PadNone),
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			query := &FilterLogsQuery{}

			err := query.UnmarshalJSON([]byte(tt.wantJSON))
			require.NoError(t, err)
			assert.Equal(t, tt.query, query)

			j, err := query.MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(j))
		})
	}
}

func TestFilterLogsQuery_JSON_OneOrList(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "bare values",
			json: `
				{
				  "address": "0x1111111111111111111111111111111111111111",
				  "topics": ["0x2222222222222222222222222222222222222222222222222222222222222222"]
				}
			`,
		},
		{
			name: "lists",
			json: `
				{
				  "address": ["0x1111111111111111111111111111111111111111"],
				  "topics": [["0x2222222222222222222222222222222222222222222222222222222222222222"]]
				}
			`,
		},
	}
	want := &FilterLogsQuery{
		Address: []Address{MustAddressFromHex("0x1111111111111111111111111111111111111111")},
		Topics: [][]Hash{
			{MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", PadNone)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := &FilterLogsQuery{}

			err := query.UnmarshalJSON([]byte(tt.json))
			require.NoError(t, err)
			assert.Equal(t, want, query)
		})
	}
}

func testBloom() []byte {
	b := make([]byte, bloomLength)
	b[bloomLength-3] = 1
	b[bloomLength-2] = 2
	b[bloomLength-1] = 3
	return b
}

var testBloomHex = func() string {
	s := "0x"
	for _, b := range testBloom() {
		s += fmt.Sprintf("%02x", b)
	}
	return s
}()
