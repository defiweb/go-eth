package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/kzg4844"
)

// AccessList represents an Ethereum access list as defined in EIP-2930.
//
// EIP-2930 introduced a new transaction type that includes an optional
// access list, which specifies a list of addresses and storage keys that the
// transaction plans to access. By declaring these accesses upfront,
// transactions can benefit from reduced gas costs for cold accesses, as
// the specified addresses and storage slots are warmed up ahead of execution.
//
// https://eips.ethereum.org/EIPS/eip-2930
type AccessList []AccessTuple

// AccessTuple is the element type of access list.
type AccessTuple struct {
	Address     Address `json:"address"`
	StorageKeys []Hash  `json:"storageKeys"`
}

// Copy creates a deep copy of the access list.
func (a *AccessList) Copy() AccessList {
	if a == nil {
		return nil
	}
	c := make(AccessList, len(*a))
	for i, tuple := range *a {
		c[i] = tuple.Copy()
	}
	return c
}

// EncodeRLP implements the rlp.Encoder interface.
func (a AccessList) EncodeRLP() ([]byte, error) {
	l := rlp.List{}
	for _, tuple := range a {
		tuple := tuple
		l.Add(&tuple)
	}
	return rlp.Encode(l)
}

// DecodeRLP implements the rlp.Decoder interface.
func (a *AccessList) DecodeRLP(data []byte) (int, error) {
	d, n, err := rlp.DecodeLazy(data)
	if err != nil {
		return 0, err
	}
	l, err := d.List()
	if err != nil {
		return 0, err
	}
	for _, tuple := range l {
		var t AccessTuple
		if err := tuple.Decode(&t); err != nil {
			return 0, err
		}
		*a = append(*a, t)
	}
	return n, nil
}

// Copy creates a deep copy of the access tuple.
func (a *AccessTuple) Copy() AccessTuple {
	keys := make([]Hash, len(a.StorageKeys))
	copy(keys, a.StorageKeys)
	return AccessTuple{
		Address:     a.Address,
		StorageKeys: keys,
	}
}

// EncodeRLP implements the rlp.Encoder interface.
func (a AccessTuple) EncodeRLP() ([]byte, error) {
	h := rlp.List{}
	for _, hash := range a.StorageKeys {
		hash := hash
		h.Add(&hash)
	}
	return rlp.Encode(rlp.List{a.Address, h})
}

// DecodeRLP implements the rlp.Decoder interface.
func (a *AccessTuple) DecodeRLP(data []byte) (int, error) {
	d, n, err := rlp.DecodeLazy(data)
	if err != nil {
		return n, err
	}
	l, err := d.List()
	if err != nil {
		return n, err
	}
	if len(l) != 2 {
		return n, fmt.Errorf("invalid access list tuple")
	}
	if err := l[0].Decode(&a.Address); err != nil {
		return n, err
	}
	h, err := l[1].List()
	if err != nil {
		return n, err
	}
	for _, item := range h {
		var hash Hash
		if err := item.Decode(&hash); err != nil {
			return n, err
		}
		a.StorageKeys = append(a.StorageKeys, hash)
	}
	return n, nil
}

// BlobInfo represents the information of an EIP-4844 blob carried in a
// transaction.
//
// EIP-4844 introduces "blob-carrying transactions" to Ethereum, which include
// a new type of data called "blobs". These blobs are large binary objects that
// are not directly accessible by the EVM but are committed to the consensus
// layer.
//
// https://eips.ethereum.org/EIPS/eip-4844
type BlobInfo struct {
	Hash    Hash         // Hash of the blob.
	Sidecar *BlobSidecar // Optional sidecar containing blob components.
}

// BlobSidecar contains the components of the blob stored by the consensus
// layer.
type BlobSidecar struct {
	Blob       kzg4844.Blob       // The actual blob data.
	Commitment kzg4844.Commitment // Commitment for the blob.
	Proof      kzg4844.Proof      // Proof for the blob.
}

// ComputeHash computes the blob hash of the given blob sidecar.
func (sc *BlobSidecar) ComputeHash() Hash {
	return crypto.KZGComputeBlobHashV1(sc.Commitment)
}

// NewBlobInfo creates a new EIP-4844 BlobInfo from the given blob, computing
// its hash, commitment, and proof.
//
// The provided blob must not be nil and must be a valid EIP-4844 blob
// of length 131072 bytes (4096 field elements of 32 bytes each).
// Each field element is a 32-byte big-endian integer not exceeding the
// BLS12-381 field modulus specified in EIP-4844.
//
// NewBlobInfo does not perform any encoding on the provided data.
//
// Returns an error if the blob is nil or if the commitment/proof computation
// fails.
func NewBlobInfo(b *kzg4844.Blob) (BlobInfo, error) {
	if b == nil {
		return BlobInfo{}, errors.New("blob is nil")
	}
	c, err := crypto.KZGBlobToCommitment(b)
	if err != nil {
		return BlobInfo{}, err
	}
	p, err := crypto.KZGComputeBlobProof(b, c)
	if err != nil {
		return BlobInfo{}, err
	}
	s := &BlobSidecar{
		Blob:       *b,
		Commitment: c,
		Proof:      p,
	}
	return BlobInfo{
		Hash:    s.ComputeHash(),
		Sidecar: s,
	}, nil
}

// TransactionOnChain represents a transaction on the blockchain.
type TransactionOnChain struct {
	Decoder          JSONTransactionDecoder // Decoder is an optional transaction decoder, if nil, the default decoder is used.
	Transaction      Transaction            // Transaction is the transaction data.
	Hash             *Hash                  // Hash of the transaction.
	BlockHash        *Hash                  // BlockHash is the hash of the block where this transaction was in.
	BlockNumber      *big.Int               // BlockNumber is the block number where this transaction was in.
	TransactionIndex *uint64                // TransactionIndex is the index of the transaction in the block.
}

// MarshalJSON implements the json.Marshaler interface.
func (t *TransactionOnChain) MarshalJSON() ([]byte, error) {
	ocd := &jsonOnChainTransaction{}
	ocd.Hash = t.Hash
	ocd.BlockHash = t.BlockHash
	ocd.BlockNumber = NumberFromBigIntPtr(t.BlockNumber)
	if t.TransactionIndex != nil {
		ocd.TransactionIndex = NumberFromUint64Ptr(*t.TransactionIndex)
	}
	return marshalJSONInline(
		t.Transaction,
		ocd,
	)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionOnChain) UnmarshalJSON(data []byte) error {
	ocd := &jsonOnChainTransaction{}
	if err := json.Unmarshal(data, ocd); err != nil {
		return err
	}
	t.Hash = ocd.Hash
	t.BlockHash = ocd.BlockHash
	t.BlockNumber = ocd.BlockNumber.Big()
	if ocd.TransactionIndex != nil {
		index := ocd.TransactionIndex.Big().Uint64()
		t.TransactionIndex = &index
	}
	dec := t.Decoder
	if dec == nil {
		dec = DefaultTransactionDecoder
	}
	tx, err := dec.DecodeJSON(data)
	if err != nil {
		return err
	}
	t.Transaction = tx
	return nil
}

type jsonOnChainTransaction struct {
	Hash             *Hash   `json:"hash,omitempty"`
	BlockHash        *Hash   `json:"blockHash,omitempty"`
	BlockNumber      *Number `json:"blockNumber,omitempty"`
	TransactionIndex *Number `json:"transactionIndex,omitempty"`
}

// TransactionReceipt represents transaction receipt.
type TransactionReceipt struct {
	TransactionHash   Hash     // TransactionHash is the hash of the transaction.
	TransactionIndex  uint64   // TransactionIndex is the index of the transaction in the block.
	BlockHash         Hash     // BlockHash is the hash of the block.
	BlockNumber       *big.Int // BlockNumber is the number of the block.
	From              Address  // From is the sender of the transaction.
	To                Address  // To is the recipient of the transaction.
	CumulativeGasUsed uint64   // CumulativeGasUsed is the total amount of gas used when this transaction was executed in the block.
	EffectiveGasPrice *big.Int // EffectiveGasPrice is the effective gas price of the transaction.
	GasUsed           uint64   // GasUsed is the amount of gas used by this specific transaction alone.
	ContractAddress   *Address // ContractAddress is the contract address created, if the transaction was a contract creation, otherwise nil.
	Logs              []Log    // Logs is the list of logs generated by the transaction.
	LogsBloom         []byte   // LogsBloom is the bloom filter for the logs of the transaction.
	Root              *Hash    // Root is the root of the state trie after the transaction.
	Status            *uint64  // Status is the status of the transaction.
}

// MarshalJSON implements the json.Marshaler interface.
func (t TransactionReceipt) MarshalJSON() ([]byte, error) {
	receipt := &jsonTransactionReceipt{
		TransactionHash:   t.TransactionHash,
		TransactionIndex:  NumberFromUint64(t.TransactionIndex),
		BlockHash:         t.BlockHash,
		BlockNumber:       NumberFromBigInt(t.BlockNumber),
		From:              t.From,
		To:                t.To,
		CumulativeGasUsed: NumberFromUint64(t.CumulativeGasUsed),
		EffectiveGasPrice: NumberFromBigInt(t.EffectiveGasPrice),
		GasUsed:           NumberFromUint64(t.GasUsed),
		ContractAddress:   t.ContractAddress,
		Logs:              t.Logs,
		LogsBloom:         t.LogsBloom,
		Root:              t.Root,
	}
	if t.Status != nil {
		status := NumberFromUint64(*t.Status)
		receipt.Status = &status
	}
	return json.Marshal(receipt)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionReceipt) UnmarshalJSON(data []byte) error {
	receipt := &jsonTransactionReceipt{}
	if err := json.Unmarshal(data, receipt); err != nil {
		return err
	}
	t.TransactionHash = receipt.TransactionHash
	t.TransactionIndex = receipt.TransactionIndex.Big().Uint64()
	t.BlockHash = receipt.BlockHash
	t.BlockNumber = receipt.BlockNumber.Big()
	t.From = receipt.From
	t.To = receipt.To
	t.CumulativeGasUsed = receipt.CumulativeGasUsed.Big().Uint64()
	t.EffectiveGasPrice = receipt.EffectiveGasPrice.Big()
	t.GasUsed = receipt.GasUsed.Big().Uint64()
	t.ContractAddress = receipt.ContractAddress
	t.Logs = receipt.Logs
	t.LogsBloom = receipt.LogsBloom
	t.Root = receipt.Root
	if receipt.Status != nil {
		status := receipt.Status.Big().Uint64()
		t.Status = &status
	}
	return nil
}

type jsonTransactionReceipt struct {
	TransactionHash   Hash     `json:"transactionHash"`
	TransactionIndex  Number   `json:"transactionIndex"`
	BlockHash         Hash     `json:"blockHash"`
	BlockNumber       Number   `json:"blockNumber"`
	From              Address  `json:"from"`
	To                Address  `json:"to"`
	CumulativeGasUsed Number   `json:"cumulativeGasUsed"`
	EffectiveGasPrice Number   `json:"effectiveGasPrice"`
	GasUsed           Number   `json:"gasUsed"`
	ContractAddress   *Address `json:"contractAddress"`
	Logs              []Log    `json:"logs"`
	LogsBloom         Bytes    `json:"logsBloom"`
	Root              *Hash    `json:"root"`
	Status            *Number  `json:"status"`
}

// Block represents a block on the blockchain.
type Block struct {
	Number            *big.Int             // Block is the block number.
	Hash              Hash                 // Hash is the hash of the block.
	ParentHash        Hash                 // ParentHash is the hash of the parent block.
	StateRoot         Hash                 // StateRoot is the root hash of the state trie.
	ReceiptsRoot      Hash                 // ReceiptsRoot is the root hash of the receipts trie.
	TransactionsRoot  Hash                 // TransactionsRoot is the root hash of the transactions trie.
	MixHash           Hash                 // MixHash is the hash of the seed used for the DAG.
	Sha3Uncles        Hash                 // Sha3Uncles is the SHA3 hash of the uncles data in the block.
	Nonce             *big.Int             // Nonce is the block's nonce.
	Miner             Address              // Miner is the address of the beneficiary to whom the mining rewards were given.
	LogsBloom         []byte               // LogsBloom is the bloom filter for the logs of the block.
	Difficulty        *big.Int             // Difficulty is the difficulty for this block.
	TotalDifficulty   *big.Int             // TotalDifficulty is the total difficulty of the chain until this block.
	Size              uint64               // Size is the size of the block in bytes.
	GasLimit          uint64               // GasLimit is the maximum gas allowed in this block.
	GasUsed           uint64               // GasUsed is the total used gas by all transactions in this block.
	Timestamp         time.Time            // Timestamp is the time at which the block was collated.
	Uncles            []Hash               // Uncles is the list of uncle hashes.
	Transactions      []TransactionOnChain // Transactions is the list of transactions in the block.
	TransactionHashes []Hash               // TransactionHashes is the list of transaction hashes in the block.
	ExtraData         []byte               // ExtraData is the "extra data" field of this block.
}

// MarshalJSON implements the json.Marshaler interface.
func (b Block) MarshalJSON() ([]byte, error) {
	block := &jsonBlock{
		Number:           NumberFromBigInt(b.Number),
		Hash:             b.Hash,
		ParentHash:       b.ParentHash,
		StateRoot:        b.StateRoot,
		ReceiptsRoot:     b.ReceiptsRoot,
		TransactionsRoot: b.TransactionsRoot,
		MixHash:          b.MixHash,
		Sha3Uncles:       b.Sha3Uncles,
		Nonce:            nonceFromBigInt(b.Nonce),
		Miner:            b.Miner,
		LogsBloom:        bloomFromBytes(b.LogsBloom),
		Difficulty:       NumberFromBigInt(b.Difficulty),
		TotalDifficulty:  NumberFromBigInt(b.TotalDifficulty),
		Size:             NumberFromUint64(b.Size),
		GasLimit:         NumberFromUint64(b.GasLimit),
		GasUsed:          NumberFromUint64(b.GasUsed),
		Timestamp:        NumberFromUint64(uint64(b.Timestamp.Unix())),
		Uncles:           b.Uncles,
		ExtraData:        b.ExtraData,
	}
	if len(b.Transactions) > 0 {
		block.Transactions.Objects = b.Transactions
	}
	if len(b.TransactionHashes) > 0 {
		block.Transactions.Hashes = b.TransactionHashes
	}
	return json.Marshal(block)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Block) UnmarshalJSON(data []byte) error {
	block := &jsonBlock{}
	if err := json.Unmarshal(data, block); err != nil {
		return err
	}
	b.Number = block.Number.Big()
	b.Hash = block.Hash
	b.ParentHash = block.ParentHash
	b.StateRoot = block.StateRoot
	b.ReceiptsRoot = block.ReceiptsRoot
	b.TransactionsRoot = block.TransactionsRoot
	b.MixHash = block.MixHash
	b.Sha3Uncles = block.Sha3Uncles
	b.Nonce = block.Nonce.Big()
	b.Miner = block.Miner
	b.LogsBloom = block.LogsBloom.Bytes()
	b.Difficulty = block.Difficulty.Big()
	b.TotalDifficulty = block.TotalDifficulty.Big()
	b.Size = block.Size.Big().Uint64()
	b.GasLimit = block.GasLimit.Big().Uint64()
	b.GasUsed = block.GasUsed.Big().Uint64()
	b.Timestamp = time.Unix(block.Timestamp.Big().Int64(), 0)
	b.Uncles = block.Uncles
	b.ExtraData = block.ExtraData
	b.Transactions = block.Transactions.Objects
	b.TransactionHashes = block.Transactions.Hashes
	return nil
}

type jsonBlock struct {
	Number           Number                `json:"number"`
	Hash             Hash                  `json:"hash"`
	ParentHash       Hash                  `json:"parentHash"`
	StateRoot        Hash                  `json:"stateRoot"`
	ReceiptsRoot     Hash                  `json:"receiptsRoot"`
	TransactionsRoot Hash                  `json:"transactionsRoot"`
	MixHash          Hash                  `json:"mixHash"`
	Sha3Uncles       Hash                  `json:"sha3Uncles"`
	Nonce            nonce                 `json:"nonce"`
	Miner            Address               `json:"miner"`
	LogsBloom        bloom                 `json:"logsBloom"`
	Difficulty       Number                `json:"difficulty"`
	TotalDifficulty  Number                `json:"totalDifficulty"`
	Size             Number                `json:"size"`
	GasLimit         Number                `json:"gasLimit"`
	GasUsed          Number                `json:"gasUsed"`
	Timestamp        Number                `json:"timestamp"`
	Uncles           []Hash                `json:"uncles"`
	ExtraData        Bytes                 `json:"extraData"`
	Transactions     jsonBlockTransactions `json:"transactions"`
}

type jsonBlockTransactions struct {
	Objects []TransactionOnChain
	Hashes  []Hash
}

func (b *jsonBlockTransactions) MarshalJSON() ([]byte, error) {
	if len(b.Objects) > 0 {
		return json.Marshal(b.Objects)
	}
	return json.Marshal(b.Hashes)
}

func (b *jsonBlockTransactions) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if bytes.IndexByte(data[1:], '{') >= 0 {
		return json.Unmarshal(data, &b.Objects)
	}
	return json.Unmarshal(data, &b.Hashes)
}

// FeeHistory contains information about the fee structure and gas usage
// over a range of blocks.
type FeeHistory struct {
	OldestBlock   uint64       // OldestBlock is the oldest block number for which the base fee and gas used are returned.
	Reward        [][]*big.Int // Reward is the reward for each block in the range [OldestBlock, LatestBlock].
	BaseFeePerGas []*big.Int   // BaseFeePerGas is the base fee per gas for each block in the range [OldestBlock, LatestBlock].
	GasUsedRatio  []float64    // GasUsedRatio is the gas used ratio for each block in the range [OldestBlock, LatestBlock].
}

// MarshalJSON implements the json.Marshaler interface.
func (f FeeHistory) MarshalJSON() ([]byte, error) {
	feeHistory := &jsonFeeHistory{
		OldestBlock:  NumberFromUint64(f.OldestBlock),
		GasUsedRatio: f.GasUsedRatio,
	}
	if len(f.Reward) > 0 {
		feeHistory.Reward = make([][]Number, len(f.Reward))
		for i, reward := range f.Reward {
			feeHistory.Reward[i] = make([]Number, len(reward))
			for j, r := range reward {
				feeHistory.Reward[i][j] = NumberFromBigInt(r)
			}
		}
	}
	if len(f.BaseFeePerGas) > 0 {
		feeHistory.BaseFeePerGas = make([]Number, len(f.BaseFeePerGas))
		for i, b := range f.BaseFeePerGas {
			feeHistory.BaseFeePerGas[i] = NumberFromBigInt(b)
		}
	}
	return json.Marshal(feeHistory)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *FeeHistory) UnmarshalJSON(input []byte) error {
	feeHistory := &jsonFeeHistory{}
	if err := json.Unmarshal(input, feeHistory); err != nil {
		return err
	}
	f.OldestBlock = feeHistory.OldestBlock.Big().Uint64()
	f.Reward = make([][]*big.Int, len(feeHistory.Reward))
	for i, reward := range feeHistory.Reward {
		f.Reward[i] = make([]*big.Int, len(reward))
		for j, r := range reward {
			f.Reward[i][j] = r.Big()
		}
	}
	f.BaseFeePerGas = make([]*big.Int, len(feeHistory.BaseFeePerGas))
	for i, b := range feeHistory.BaseFeePerGas {
		f.BaseFeePerGas[i] = b.Big()
	}
	f.GasUsedRatio = feeHistory.GasUsedRatio
	return nil
}

// jsonFeeHistory is the JSON representation of a fee history.
type jsonFeeHistory struct {
	OldestBlock   Number     `json:"oldestBlock"`
	Reward        [][]Number `json:"reward"`
	BaseFeePerGas []Number   `json:"baseFeePerGas"`
	GasUsedRatio  []float64  `json:"gasUsedRatio"`
}

// Log represents a contract log event.
type Log struct {
	Address          Address  // Address of the contract that generated the event
	Topics           []Hash   // Topics provide information about the event type.
	Data             []byte   // Data contains the non-indexed arguments of the event.
	BlockHash        *Hash    // BlockHash is the hash of the block where this log was in. Nil when pending.
	BlockNumber      *big.Int // BlockNumber is the block number where this log was in. Nil when pending.
	TransactionHash  *Hash    // TransactionHash is the hash of the transaction that generated this log. Nil when pending.
	TransactionIndex *uint64  // TransactionIndex is the index of the transaction in the block. Nil when pending.
	LogIndex         *uint64  // LogIndex is the index of the log in the block. Nil when pending.
	Removed          bool     // Removed is true if the log was reverted due to a chain reorganization. False if unknown.
}

// MarshalJSON implements the json.Marshaler interface.
func (l Log) MarshalJSON() ([]byte, error) {
	j := &jsonLog{}
	j.Address = l.Address
	j.Topics = l.Topics
	j.Data = l.Data
	j.BlockHash = l.BlockHash
	if l.BlockNumber != nil {
		j.BlockNumber = NumberFromBigIntPtr(l.BlockNumber)
	}
	j.TransactionHash = l.TransactionHash
	if l.TransactionIndex != nil {
		j.TransactionIndex = NumberFromUint64Ptr(*l.TransactionIndex)
	}
	if l.LogIndex != nil {
		j.LogIndex = NumberFromUint64Ptr(*l.LogIndex)
	}
	j.Removed = l.Removed
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Log) UnmarshalJSON(input []byte) error {
	log := &jsonLog{}
	if err := json.Unmarshal(input, log); err != nil {
		return err
	}
	l.Address = log.Address
	l.Topics = log.Topics
	l.Data = log.Data
	l.BlockHash = log.BlockHash
	if log.BlockNumber != nil {
		l.BlockNumber = log.BlockNumber.Big()
	}
	l.TransactionHash = log.TransactionHash
	if log.TransactionIndex != nil {
		l.TransactionIndex = new(uint64)
		*l.TransactionIndex = log.TransactionIndex.Big().Uint64()
	}
	if log.LogIndex != nil {
		l.LogIndex = new(uint64)
		*l.LogIndex = log.LogIndex.Big().Uint64()
	}
	l.Removed = log.Removed
	return nil
}

type jsonLog struct {
	Address          Address `json:"address"`
	Topics           []Hash  `json:"topics"`
	Data             Bytes   `json:"data"`
	BlockHash        *Hash   `json:"blockHash"`
	BlockNumber      *Number `json:"blockNumber"`
	TransactionHash  *Hash   `json:"transactionHash"`
	TransactionIndex *Number `json:"transactionIndex"`
	LogIndex         *Number `json:"logIndex"`
	Removed          bool    `json:"removed"`
}

// FilterLogsQuery represents a query to filter logs.
type FilterLogsQuery struct {
	Address   []Address
	FromBlock *BlockNumber
	ToBlock   *BlockNumber
	Topics    [][]Hash
	BlockHash *Hash
}

// NewFilterLogsQuery creates a new FilterLogsQuery.
func NewFilterLogsQuery() *FilterLogsQuery {
	return &FilterLogsQuery{}
}

// SetAddresses sets the addresses to filter logs.
func (q *FilterLogsQuery) SetAddresses(addresses ...Address) {
	q.Address = addresses
}

// AddAddresses adds addresses to filter logs.
func (q *FilterLogsQuery) AddAddresses(addresses ...Address) {
	q.Address = append(q.Address, addresses...)
}

// SetFromBlock sets the starting block number to filter logs.
func (q *FilterLogsQuery) SetFromBlock(fromBlock *BlockNumber) {
	q.FromBlock = fromBlock
}

// SetToBlock sets the ending block number to filter logs.
func (q *FilterLogsQuery) SetToBlock(toBlock *BlockNumber) {
	q.ToBlock = toBlock
}

// SetTopics sets the topics to filter logs.
func (q *FilterLogsQuery) SetTopics(topics ...[]Hash) {
	q.Topics = topics
}

// AddTopics adds topics to filter logs.
func (q *FilterLogsQuery) AddTopics(topics ...[]Hash) {
	q.Topics = append(q.Topics, topics...)
}

// SetBlockHash sets the block hash to filter logs.
func (q *FilterLogsQuery) SetBlockHash(blockHash *Hash) {
	q.BlockHash = blockHash
}

// MarshalJSON implements the json.Marshaler interface.
func (q FilterLogsQuery) MarshalJSON() ([]byte, error) {
	logsQuery := &jsonFilterLogsQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		BlockHash: q.BlockHash,
	}
	if len(q.Address) > 0 {
		logsQuery.Address = make([]Address, len(q.Address))
		copy(logsQuery.Address, q.Address)
	}
	if len(q.Topics) > 0 {
		logsQuery.Topics = make([]oneOrList[Hash], len(q.Topics))
		for i, t := range q.Topics {
			logsQuery.Topics[i] = make([]Hash, len(t))
			copy(logsQuery.Topics[i], t)
		}
	}
	return json.Marshal(logsQuery)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (q *FilterLogsQuery) UnmarshalJSON(input []byte) error {
	logsQuery := &jsonFilterLogsQuery{}
	if err := json.Unmarshal(input, logsQuery); err != nil {
		return err
	}
	q.FromBlock = logsQuery.FromBlock
	q.ToBlock = logsQuery.ToBlock
	q.BlockHash = logsQuery.BlockHash
	if len(logsQuery.Address) > 0 {
		q.Address = make([]Address, len(logsQuery.Address))
		copy(q.Address, logsQuery.Address)
	}
	if len(logsQuery.Topics) > 0 {
		q.Topics = make([][]Hash, len(logsQuery.Topics))
		for i, t := range logsQuery.Topics {
			q.Topics[i] = make([]Hash, len(t))
			copy(q.Topics[i], t)
		}
	}
	return nil
}

type jsonFilterLogsQuery struct {
	Address   oneOrList[Address] `json:"address"`
	FromBlock *BlockNumber       `json:"fromBlock,omitempty"`
	ToBlock   *BlockNumber       `json:"toBlock,omitempty"`
	Topics    []oneOrList[Hash]  `json:"topics"`
	BlockHash *Hash              `json:"blockhash,omitempty"`
}

// SyncStatus represents the sync status of a node.
type SyncStatus struct {
	StartingBlock BlockNumber `json:"startingBlock"`
	CurrentBlock  BlockNumber `json:"currentBlock"`
	HighestBlock  BlockNumber `json:"highestBlock"`
}
