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
	// Hash is the blob's versioned hash.
	Hash crypto.KZGHash

	// Sidecar contains the blob components. Nil when the sidecar is not available.
	Sidecar *BlobSidecar
}

// BlobSidecar contains the components of the blob stored by the consensus
// layer.
type BlobSidecar struct {
	// Blob is the blob data.
	Blob crypto.KZGBlob

	// Commitment is the KZG commitment for the blob.
	Commitment crypto.KZGCommitment

	// Proof is the KZG proof for the blob.
	Proof crypto.KZGProof
}

// ComputeHash computes the blob hash of the given blob sidecar.
func (sc *BlobSidecar) ComputeHash() crypto.KZGHash {
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
func NewBlobInfo(b *crypto.KZGBlob) (BlobInfo, error) {
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
	// Decoder is an optional transaction decoder. If nil, the default
	// decoder is used.
	Decoder JSONTransactionDecoder

	// Transaction is the decoded transaction data.
	Transaction Transaction

	// Hash is the transaction hash.
	Hash *Hash

	// BlockHash is the hash of the block containing this transaction.
	BlockHash *Hash

	// BlockNumber is the number of the block containing this transaction.
	BlockNumber *big.Int

	// TransactionIndex is the index of the transaction within the block.
	TransactionIndex *uint64
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
	return marshalJSONMerge(
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
	if ocd.BlockNumber != nil {
		t.BlockNumber = ocd.BlockNumber.Big()
	} else {
		t.BlockNumber = nil
	}
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
	// TransactionHash is the hash of the transaction.
	TransactionHash Hash

	// TransactionIndex is the index of the transaction within the block.
	TransactionIndex uint64

	// BlockHash is the hash of the block containing the transaction.
	BlockHash Hash

	// BlockNumber is the number of the block containing the transaction.
	BlockNumber *big.Int

	// From is the sender of the transaction.
	From Address

	// To is the recipient of the transaction.
	To Address

	// CumulativeGasUsed is the total gas used in the block up to and including
	// this transaction.
	CumulativeGasUsed uint64

	// EffectiveGasPrice is the actual gas price paid for this transaction.
	EffectiveGasPrice *big.Int

	// GasUsed is the gas used by this transaction.
	GasUsed uint64

	// ContractAddress is the address of the created contract, or nil if the
	// transaction was not a contract creation.
	ContractAddress *Address

	// Logs is the list of logs generated by the transaction.
	Logs []Log

	// LogsBloom is the bloom filter for the transaction's logs.
	LogsBloom []byte

	// Root is the post-transaction state root (pre-Byzantium only).
	Root *Hash

	// Status is 1 if the transaction succeeded, 0 if it reverted.
	Status *uint64
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
	// Number is the block number.
	Number *big.Int

	// Hash is the hash of the block.
	Hash Hash

	// ParentHash is the hash of the parent block.
	ParentHash Hash

	// StateRoot is the root hash of the state trie.
	StateRoot Hash

	// ReceiptsRoot is the root hash of the receipts trie.
	ReceiptsRoot Hash

	// TransactionsRoot is the root hash of the transactions trie.
	TransactionsRoot Hash

	// MixHash is the hash mixed with the nonce to prove proof-of-work.
	MixHash Hash

	// Sha3Uncles is the SHA3 hash of the uncle block headers.
	Sha3Uncles Hash

	// Nonce is the proof-of-work nonce.
	Nonce *big.Int

	// Miner is the address of the block's fee recipient.
	Miner Address

	// LogsBloom is the bloom filter for the block's logs.
	LogsBloom []byte

	// Difficulty is the proof-of-work difficulty for this block.
	Difficulty *big.Int

	// TotalDifficulty is the cumulative chain difficulty up to this block.
	TotalDifficulty *big.Int

	// Size is the size of the block in bytes.
	Size uint64

	// GasLimit is the maximum gas allowed in this block.
	GasLimit uint64

	// GasUsed is the total gas used by all transactions in this block.
	GasUsed uint64

	// Timestamp is the time at which the block was produced.
	Timestamp time.Time

	// Uncles is the list of uncle block hashes.
	Uncles []Hash

	// Transactions is the list of full transactions in the block.
	// Populated when the block was requested with full transaction objects.
	Transactions []TransactionOnChain

	// TransactionHashes is the list of transaction hashes in the block.
	// Populated when the block was requested with transaction hashes only.
	TransactionHashes []Hash

	// ExtraData is the arbitrary extra data field of this block.
	ExtraData []byte
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
	// OldestBlock is the oldest block number in the returned range.
	OldestBlock uint64

	// Reward contains the requested priority fee percentiles for each block.
	Reward [][]*big.Int

	// BaseFeePerGas is the base fee per gas for each block. The array has
	// blockCount+1 entries; the last entry is the next block's predicted fee.
	BaseFeePerGas []*big.Int

	// GasUsedRatio is the ratio of gas used to gas limit for each block.
	GasUsedRatio []float64

	// BaseFeePerBlobGas is the base fee per blob gas for each block (EIP-4844).
	// The array has blockCount+1 entries.
	BaseFeePerBlobGas []*big.Int

	// BlobGasUsedRatio is the ratio of blob gas used to the blob gas limit
	// for each block (EIP-4844).
	BlobGasUsedRatio []float64
}

// MarshalJSON implements the json.Marshaler interface.
func (f FeeHistory) MarshalJSON() ([]byte, error) {
	feeHistory := &jsonFeeHistory{
		OldestBlock:      NumberFromUint64(f.OldestBlock),
		GasUsedRatio:     f.GasUsedRatio,
		BlobGasUsedRatio: f.BlobGasUsedRatio,
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
	if len(f.BaseFeePerBlobGas) > 0 {
		feeHistory.BaseFeePerBlobGas = make([]Number, len(f.BaseFeePerBlobGas))
		for i, b := range f.BaseFeePerBlobGas {
			feeHistory.BaseFeePerBlobGas[i] = NumberFromBigInt(b)
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
	f.BaseFeePerBlobGas = make([]*big.Int, len(feeHistory.BaseFeePerBlobGas))
	for i, b := range feeHistory.BaseFeePerBlobGas {
		f.BaseFeePerBlobGas[i] = b.Big()
	}
	f.BlobGasUsedRatio = feeHistory.BlobGasUsedRatio
	return nil
}

// jsonFeeHistory is the JSON representation of a fee history.
type jsonFeeHistory struct {
	OldestBlock       Number     `json:"oldestBlock"`
	Reward            [][]Number `json:"reward,omitempty"`
	BaseFeePerGas     []Number   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio      []float64  `json:"gasUsedRatio,omitempty"`
	BaseFeePerBlobGas []Number   `json:"baseFeePerBlobGas,omitempty"`
	BlobGasUsedRatio  []float64  `json:"blobGasUsedRatio,omitempty"`
}

// AccessListResult represents the result of an eth_createAccessList call.
type AccessListResult struct {
	// AccessList is the list of addresses and storage keys accessed by
	// the transaction.
	AccessList AccessList

	// GasUsed is the gas used by the transaction with the given access list.
	GasUsed uint64

	// Error is a revert message if the transaction would revert.
	// Empty when the transaction succeeds.

	Error string
}

// MarshalJSON implements the json.Marshaler interface.
func (a AccessListResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAccessListResult{
		AccessList: a.AccessList,
		GasUsed:    NumberFromUint64(a.GasUsed),
		Error:      a.Error,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *AccessListResult) UnmarshalJSON(data []byte) error {
	j := &jsonAccessListResult{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	a.AccessList = j.AccessList
	a.GasUsed = j.GasUsed.Big().Uint64()
	a.Error = j.Error
	return nil
}

type jsonAccessListResult struct {
	AccessList AccessList `json:"accessList"`
	GasUsed    Number     `json:"gasUsed"`
	Error      string     `json:"error,omitempty"`
}

// StorageProof represents a single storage proof entry returned by
// eth_getProof.
type StorageProof struct {
	// Key is the storage slot key.
	Key Hash

	// Value is the storage slot value.
	Value *big.Int

	// Proof is the array of RLP-serialized Merkle-Patricia trie nodes
	// proving the storage slot value.
	Proof []Bytes
}

// MarshalJSON implements the json.Marshaler interface.
func (s StorageProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonStorageProof{
		Key:   s.Key,
		Value: NumberFromBigInt(s.Value),
		Proof: s.Proof,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *StorageProof) UnmarshalJSON(data []byte) error {
	j := &jsonStorageProof{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	s.Key = j.Key
	s.Value = j.Value.Big()
	s.Proof = j.Proof
	return nil
}

type jsonStorageProof struct {
	Key   Hash    `json:"key"`
	Value Number  `json:"value"`
	Proof []Bytes `json:"proof"`
}

// AccountProof represents the result of an eth_getProof call as defined in
// EIP-1186.
//
// https://eips.ethereum.org/EIPS/eip-1186
type AccountProof struct {
	// Address is the account address.
	Address Address

	// AccountProof is the array of RLP-serialized Merkle-Patricia trie nodes
	// from the state root to the account leaf.
	AccountProof []Bytes

	// Balance is the account balance in wei.
	Balance *big.Int

	// CodeHash is the Keccak-256 hash of the account's code.
	CodeHash Hash

	// Nonce is the account nonce.
	Nonce uint64

	// StorageHash is the Keccak-256 hash of the account's storage trie root.
	StorageHash Hash

	// StorageProof contains the proofs for each requested storage key.
	StorageProof []StorageProof
}

// MarshalJSON implements the json.Marshaler interface.
func (a AccountProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAccountProof{
		Address:      a.Address,
		AccountProof: a.AccountProof,
		Balance:      NumberFromBigInt(a.Balance),
		CodeHash:     a.CodeHash,
		Nonce:        NumberFromUint64(a.Nonce),
		StorageHash:  a.StorageHash,
		StorageProof: a.StorageProof,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (a *AccountProof) UnmarshalJSON(data []byte) error {
	j := &jsonAccountProof{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	a.Address = j.Address
	a.AccountProof = j.AccountProof
	a.Balance = j.Balance.Big()
	a.CodeHash = j.CodeHash
	a.Nonce = j.Nonce.Big().Uint64()
	a.StorageHash = j.StorageHash
	a.StorageProof = j.StorageProof
	return nil
}

type jsonAccountProof struct {
	Address      Address        `json:"address"`
	AccountProof []Bytes        `json:"accountProof"`
	Balance      Number         `json:"balance"`
	CodeHash     Hash           `json:"codeHash"`
	Nonce        Number         `json:"nonce"`
	StorageHash  Hash           `json:"storageHash"`
	StorageProof []StorageProof `json:"storageProof"`
}

// Log represents a contract log event.
type Log struct {
	// Address is the address of the contract that emitted the event.
	Address Address

	// Topics are the indexed event parameters.
	Topics []Hash

	// Data contains the non-indexed event parameters.
	Data []byte

	// BlockHash is the hash of the block containing this log.
	// Nil when the log is pending.
	BlockHash *Hash

	// BlockNumber is the number of the block containing this log.
	// Nil when the log is pending.
	BlockNumber *big.Int

	// TransactionHash is the hash of the transaction that emitted this log.
	// Nil when the log is pending.
	TransactionHash *Hash

	// TransactionIndex is the index of the transaction within the block.
	// Nil when the log is pending.
	TransactionIndex *uint64

	// LogIndex is the index of the log within the block.
	// Nil when the log is pending.
	LogIndex *uint64

	// Removed is true if this log was reverted by a chain reorganisation.
	Removed bool
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
