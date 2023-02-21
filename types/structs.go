package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/defiweb/go-rlp"
)

type Call struct {
	From       Address    // the sender of the 'transaction'.
	To         *Address   // the destination contract (nil for contract creation).
	Gas        uint64     // if 0, the call executes with near-infinite gas.
	GasPrice   *big.Int   // wei <-> gas exchange ratio.
	GasFeeCap  *big.Int   // EIP-1559 fee cap per gas.
	GasTipCap  *big.Int   // EIP-1559 tip per gas.
	Value      *big.Int   // amount of wei sent along with the call.
	Data       []byte     // input data, usually an ABI-encoded contract method invocation.
	AccessList AccessList // EIP-2930 access list.
}

func (c Call) MarshalJSON() ([]byte, error) {
	call := &jsonCall{
		From:       c.From,
		To:         c.To,
		Data:       c.Data,
		AccessList: c.AccessList,
	}
	if c.Gas != 0 {
		call.Gas = Uint64ToNumberPtr(c.Gas)
	}
	if c.GasPrice != nil {
		gasPrice := BigIntToNumber(c.GasPrice)
		call.GasPrice = &gasPrice
	}
	if c.GasFeeCap != nil {
		gasFeeCap := BigIntToNumber(c.GasFeeCap)
		call.GasFeeCap = &gasFeeCap
	}
	if c.GasTipCap != nil {
		gasTipCap := BigIntToNumber(c.GasTipCap)
		call.GasTipCap = &gasTipCap
	}
	if c.Value != nil {
		value := BigIntToNumber(c.Value)
		call.Value = &value
	}
	return json.Marshal(call)
}

func (c *Call) UnmarshalJSON(data []byte) error {
	call := &jsonCall{}
	if err := json.Unmarshal(data, call); err != nil {
		return err
	}
	c.From = call.From
	c.To = call.To
	c.Gas = call.Gas.Big().Uint64()
	c.Data = call.Data
	c.AccessList = call.AccessList
	if call.GasPrice != nil {
		c.GasPrice = call.GasPrice.Big()
	}
	if call.GasFeeCap != nil {
		c.GasFeeCap = call.GasFeeCap.Big()
	}
	if call.GasTipCap != nil {
		c.GasTipCap = call.GasTipCap.Big()
	}
	if call.Value != nil {
		c.Value = call.Value.Big()
	}
	return nil
}

type jsonCall struct {
	From       Address    `json:"from"`
	To         *Address   `json:"to,omitempty"`
	Gas        *Number    `json:"gas"`
	GasPrice   *Number    `json:"gasPrice,omitempty"`
	GasFeeCap  *Number    `json:"maxFeePerGas,omitempty"`
	GasTipCap  *Number    `json:"maxPriorityFeePerGas,omitempty"`
	Value      *Number    `json:"value,omitempty"`
	Data       Bytes      `json:"data"`
	AccessList AccessList `json:"accessList,omitempty"`
}

type TransactionType uint64

// Transaction types.
const (
	LegacyTxType TransactionType = iota
	AccessListTxType
	DynamicFeeTxType
)

// Transaction represents a transaction.
type Transaction struct {
	Type      TransactionType
	From      *Address
	To        *Address
	Gas       uint64
	GasPrice  *big.Int
	Input     []byte
	Nonce     *big.Int
	Value     *big.Int
	Signature Signature

	// On-chain fields
	Hash             Hash
	BlockHash        *Hash
	BlockNumber      *uint64
	TransactionIndex uint64

	// EIP-2930 fields
	ChainID    *big.Int
	AccessList AccessList

	// EIP-1559 fields
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

// Raw returns the raw transaction data that could be sent to the network.
func (t Transaction) Raw() ([]byte, error) {
	return t.EncodeRLP()
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	transaction := &jsonTransaction{
		Hash:             t.Hash,
		BlockHash:        t.BlockHash,
		TransactionIndex: Uint64ToNumber(t.TransactionIndex),
		From:             t.From,
		To:               t.To,
		Gas:              Uint64ToNumber(t.Gas),
		GasPrice:         BigIntToNumber(t.GasPrice),
		Input:            t.Input,
		Nonce:            BigIntToNumber(t.Nonce),
		Value:            BigIntToNumber(t.Value),
		V:                BigIntToNumber(t.Signature.BigV()),
		R:                BigIntToNumber(t.Signature.BigR()),
		S:                BigIntToNumber(t.Signature.BigS()),
	}
	if t.BlockNumber != nil {
		blockNumber := Uint64ToNumber(*t.BlockNumber)
		transaction.BlockNumber = &blockNumber
	}
	return json.Marshal(transaction)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	transaction := &jsonTransaction{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return err
	}
	signature, err := BigIntToSignature(transaction.V.Big(), transaction.R.Big(), transaction.S.Big())
	if err != nil {
		return err
	}
	t.Hash = transaction.Hash
	t.BlockHash = transaction.BlockHash
	t.TransactionIndex = transaction.TransactionIndex.Big().Uint64()
	t.From = transaction.From
	t.To = transaction.To
	t.Gas = transaction.Gas.Big().Uint64()
	t.GasPrice = transaction.GasPrice.Big()
	t.Input = transaction.Input
	t.Nonce = transaction.Nonce.Big()
	t.Value = transaction.Value.Big()
	t.Signature = signature
	if transaction.BlockNumber != nil {
		blockNumber := transaction.BlockNumber.Big().Uint64()
		t.BlockNumber = &blockNumber
	}
	return nil
}

func (t Transaction) EncodeRLP() ([]byte, error) {
	l := rlp.NewList()
	if t.Type != LegacyTxType {
		l.Append(rlp.NewBigInt(t.ChainID))
	}
	l.Append(rlp.NewBigInt(t.Nonce))
	if t.Type == DynamicFeeTxType {
		l.Append(rlp.NewBigInt(t.MaxPriorityFeePerGas))
		l.Append(rlp.NewBigInt(t.MaxFeePerGas))
	} else {
		l.Append(rlp.NewBigInt(t.GasPrice))
	}
	l.Append(rlp.NewUint(t.Gas))
	l.Append(t.To)
	l.Append(rlp.NewBigInt(t.Value))
	l.Append(rlp.NewBytes(t.Input))
	if t.Type != LegacyTxType {
		l.Append(&t.AccessList)
	}
	l.Append(rlp.NewBigInt(t.Signature.BigV()))
	l.Append(rlp.NewBigInt(t.Signature.BigR()))
	l.Append(rlp.NewBigInt(t.Signature.BigS()))
	b, err := rlp.Encode(l)
	if err != nil {
		return nil, err
	}
	if t.Type == AccessListTxType {
		b = append([]byte{byte(AccessListTxType)}, b...)
	}
	if t.Type == DynamicFeeTxType {
		b = append([]byte{byte(DynamicFeeTxType)}, b...)
	}
	return b, nil
}

func (t *Transaction) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	typ := TransactionType(data[0])
	var (
		elemNum int
		elemIdx int
	)
	switch typ {
	default:
		t.Type = LegacyTxType
		elemNum = 9
	case AccessListTxType:
		t.Type = AccessListTxType
		elemNum = 11
		data = data[1:]
	case DynamicFeeTxType:
		t.Type = DynamicFeeTxType
		elemNum = 12
		data = data[1:]
	}
	d, n, err := rlp.Decode(data)
	if err != nil {
		return 0, err
	}
	l, err := d.GetList()
	if err != nil {
		return 0, err
	}
	if len(l) != elemNum {
		return 0, errors.New("invalid transaction RLP")
	}
	if t.Type != LegacyTxType {
		if t.ChainID, err = l[elemIdx].GetBigInt(); err != nil {
			return 0, err
		}
		elemIdx++
	}
	if t.Nonce, err = l[elemIdx].GetBigInt(); err != nil {
		return 0, err
	}
	elemIdx++
	if t.Type == DynamicFeeTxType {
		if t.MaxPriorityFeePerGas, err = l[elemIdx].GetBigInt(); err != nil {
			return 0, err
		}
		elemIdx++
		if t.MaxFeePerGas, err = l[elemIdx].GetBigInt(); err != nil {
			return 0, err
		}
		elemIdx++
	} else {
		if t.GasPrice, err = l[elemIdx].GetBigInt(); err != nil {
			return 0, err
		}
		elemIdx++
	}
	if t.Gas, err = l[elemIdx].GetUint(); err != nil {
		return 0, err
	}
	elemIdx++
	if err := l[elemIdx].Get(&Address{}, func(i rlp.Item) { t.To = i.(*Address) }); err != nil {
		return 0, err
	}
	elemIdx++
	if t.Value, err = l[elemIdx].GetBigInt(); err != nil {
		return 0, err
	}
	elemIdx++
	if t.Input, err = l[elemIdx].GetBytes(); err != nil {
		return 0, err
	}
	elemIdx++
	if t.Type != LegacyTxType {
		if err := l[elemIdx].Get(&AccessList{}, func(i rlp.Item) { t.AccessList = *i.(*AccessList) }); err != nil {
			return 0, err
		}
		elemIdx++
	}
	var v, r, s *big.Int
	if v, err = l[elemIdx].GetBigInt(); err != nil {
		return 0, err
	}
	elemIdx++
	if r, err = l[elemIdx].GetBigInt(); err != nil {
		return 0, err
	}
	elemIdx++
	if s, err = l[elemIdx].GetBigInt(); err != nil {
		return 0, err
	}
	sig, err := BigIntToSignature(v, r, s)
	if err != nil {
		return 0, err
	}
	t.Signature = sig
	if t.Type == LegacyTxType {
		return n, nil
	}
	return n + 1, nil
}

// SigningHash returns the transaction hash to be signed by the sender.
func (t Transaction) SigningHash(h HashFunc) (Hash, error) {
	l := rlp.NewList()
	if t.Type != LegacyTxType {
		l.Append(rlp.NewBigInt(t.ChainID))
	}
	l.Append(rlp.NewBigInt(t.Nonce))
	if t.Type == DynamicFeeTxType {
		l.Append(rlp.NewBigInt(t.MaxPriorityFeePerGas))
		l.Append(rlp.NewBigInt(t.MaxFeePerGas))
	} else {
		l.Append(rlp.NewBigInt(t.GasPrice))
	}
	l.Append(rlp.NewUint(t.Gas))
	l.Append(t.To)
	l.Append(rlp.NewBigInt(t.Value))
	l.Append(rlp.NewBytes(t.Input))
	if t.Type != LegacyTxType {
		l.Append(&t.AccessList)
	}
	// EIP-155 replay-protection
	if t.ChainID != nil && t.ChainID.Sign() != 0 && t.Type == LegacyTxType {
		l.Append(rlp.NewBigInt(t.ChainID))
		l.Append(rlp.NewBigInt(big.NewInt(0)))
		l.Append(rlp.NewBigInt(big.NewInt(0)))
	}
	b, err := rlp.Encode(l)
	if err != nil {
		return ZeroHash, err
	}
	if t.Type == AccessListTxType {
		b = append([]byte{byte(AccessListTxType)}, b...)
	}
	if t.Type == DynamicFeeTxType {
		b = append([]byte{byte(DynamicFeeTxType)}, b...)
	}
	return h(b), nil
}

type jsonTransaction struct {
	Hash             Hash     `json:"hash"`
	BlockHash        *Hash    `json:"blockHash"`
	BlockNumber      *Number  `json:"blockNumber"`
	TransactionIndex Number   `json:"transactionIndex"`
	From             *Address `json:"from"`
	To               *Address `json:"to"`
	Gas              Number   `json:"gas"`
	GasPrice         Number   `json:"gasPrice"`
	Input            Bytes    `json:"input"`
	Nonce            Number   `json:"nonce"`
	Value            Number   `json:"value"`
	V                Number   `json:"v"`
	R                Number   `json:"r"`
	S                Number   `json:"s"`
}

// AccessList is an EIP-2930 access list.
type AccessList []AccessTuple

// AccessTuple is the element type of access list.
type AccessTuple struct {
	Address     Address `json:"address"`
	StorageKeys []Hash  `json:"storageKeys"`
}

func (a AccessList) EncodeRLP() ([]byte, error) {
	l := rlp.NewList()
	for _, tuple := range a {
		tuple := tuple
		l.Append(&tuple)
	}
	return rlp.Encode(l)
}

func (a *AccessList) DecodeRLP(data []byte) (int, error) {
	d, n, err := rlp.Decode(data)
	if err != nil {
		return 0, err
	}
	l, err := d.GetList()
	if err != nil {
		return 0, err
	}
	for _, tuple := range l {
		var t AccessTuple
		if err := tuple.DecodeInto(&t); err != nil {
			return 0, err
		}
		*a = append(*a, t)
	}
	return n, nil
}

func (a AccessTuple) EncodeRLP() ([]byte, error) {
	h := rlp.NewList()
	for _, hash := range a.StorageKeys {
		hash := hash
		h.Append(&hash)
	}
	return rlp.Encode(rlp.NewList(&a.Address, h))
}

func (a *AccessTuple) DecodeRLP(data []byte) (int, error) {
	d, n, err := rlp.Decode(data)
	if err != nil {
		return n, err
	}
	l, err := d.GetList()
	if err != nil {
		return n, err
	}
	if len(l) != 2 {
		return n, fmt.Errorf("invalid access list tuple")
	}
	if err := l[0].DecodeInto(&a.Address); err != nil {
		return n, err
	}
	h, err := l[1].GetList()
	if err != nil {
		return n, err
	}
	for _, item := range h {
		var hash Hash
		if err := item.DecodeInto(&hash); err != nil {
			return n, err
		}
		a.StorageKeys = append(a.StorageKeys, hash)
	}
	return n, nil
}

// TransactionReceipt represents transaction receipt.
type TransactionReceipt struct {
	TransactionHash   Hash
	TransactionIndex  uint64
	BlockHash         Hash
	BlockNumber       uint64
	From              Address
	To                Address
	CumulativeGasUsed uint64
	EffectiveGasPrice *big.Int
	GasUsed           uint64
	ContractAddress   *Address
	Logs              []Log
	LogsBloom         []byte
	Root              *Hash
	Status            *uint64
}

func (t TransactionReceipt) MarshalJSON() ([]byte, error) {
	receipt := &jsonTransactionReceipt{
		TransactionHash:   t.TransactionHash,
		TransactionIndex:  Uint64ToNumber(t.TransactionIndex),
		BlockHash:         t.BlockHash,
		BlockNumber:       Uint64ToNumber(t.BlockNumber),
		From:              t.From,
		To:                t.To,
		CumulativeGasUsed: Uint64ToNumber(t.CumulativeGasUsed),
		EffectiveGasPrice: BigIntToNumber(t.EffectiveGasPrice),
		GasUsed:           Uint64ToNumber(t.GasUsed),
		ContractAddress:   t.ContractAddress,
		Logs:              t.Logs,
		LogsBloom:         t.LogsBloom,
		Root:              t.Root,
	}
	if t.Status != nil {
		status := Uint64ToNumber(*t.Status)
		receipt.Status = &status
	}
	return json.Marshal(receipt)
}

func (t *TransactionReceipt) UnmarshalJSON(data []byte) error {
	receipt := &jsonTransactionReceipt{}
	if err := json.Unmarshal(data, receipt); err != nil {
		return err
	}
	t.TransactionHash = receipt.TransactionHash
	t.TransactionIndex = receipt.TransactionIndex.Big().Uint64()
	t.BlockHash = receipt.BlockHash
	t.BlockNumber = receipt.BlockNumber.Big().Uint64()
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

// SignTransaction represents a transaction to be signed.
type SignTransaction struct {
	From     Address
	To       *Address
	Gas      *uint64
	GasPrice *big.Int
	Data     []byte
	Nonce    *big.Int
	Value    *big.Int
}

func (t SignTransaction) MarshalJSON() ([]byte, error) {
	transaction := &jsonSignTransaction{
		From: t.From,
		To:   t.To,
		Data: t.Data,
	}
	if t.Gas != nil {
		gas := Uint64ToNumber(*t.Gas)
		transaction.Gas = &gas
	}
	if t.GasPrice != nil {
		gasPrice := BigIntToNumber(t.GasPrice)
		transaction.GasPrice = &gasPrice
	}
	if t.Nonce != nil {
		nonce := BigIntToNumber(t.Nonce)
		transaction.Nonce = &nonce
	}
	if t.Value != nil {
		value := BigIntToNumber(t.Value)
		transaction.Value = &value
	}
	return json.Marshal(transaction)
}

func (t *SignTransaction) UnmarshalJSON(data []byte) error {
	transaction := &jsonSignTransaction{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return err
	}
	t.From = transaction.From
	t.To = transaction.To
	t.Data = transaction.Data
	if transaction.Gas != nil {
		gas := transaction.Gas.Big().Uint64()
		t.Gas = &gas
	}
	if transaction.GasPrice != nil {
		gasPrice := transaction.GasPrice.Big()
		t.GasPrice = gasPrice
	}
	if transaction.Nonce != nil {
		nonce := transaction.Nonce.Big()
		t.Nonce = nonce
	}
	if transaction.Value != nil {
		value := transaction.Value.Big()
		t.Value = value
	}
	return nil
}

type jsonSignTransaction struct {
	From     Address  `json:"from"`
	To       *Address `json:"to,omitempty"`
	Gas      *Number  `json:"gas,omitempty"`
	GasPrice *Number  `json:"gasPrice,omitempty"`
	Data     Bytes    `json:"data"`
	Nonce    *Number  `json:"nonce,omitempty"`
	Value    *Number  `json:"value,omitempty"`
}

type Block struct {
	Number            uint64
	Hash              Hash
	ParentHash        Hash
	StateRoot         Hash
	ReceiptsRoot      Hash
	TransactionsRoot  Hash
	MixHash           Hash
	Sha3Uncles        Hash
	Nonce             *big.Int
	Miner             Address
	LogsBloom         []byte
	Difficulty        *big.Int
	TotalDifficulty   *big.Int
	Size              uint64
	GasLimit          uint64
	GasUsed           uint64
	Timestamp         time.Time
	Uncles            []Hash
	Transactions      []Transaction
	TransactionHashes []Hash
	ExtraData         []byte
}

func (b Block) MarshalJSON() ([]byte, error) {
	block := &jsonBlock{
		Number:           Uint64ToNumber(b.Number),
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
		Difficulty:       BigIntToNumber(b.Difficulty),
		TotalDifficulty:  BigIntToNumber(b.TotalDifficulty),
		Size:             Uint64ToNumber(b.Size),
		GasLimit:         Uint64ToNumber(b.GasLimit),
		GasUsed:          Uint64ToNumber(b.GasUsed),
		Timestamp:        Uint64ToNumber(uint64(b.Timestamp.Unix())),
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

func (b *Block) UnmarshalJSON(data []byte) error {
	block := &jsonBlock{}
	if err := json.Unmarshal(data, block); err != nil {
		return err
	}
	b.Number = block.Number.Big().Uint64()
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
	Nonce            hexNonce              `json:"nonce"`
	Miner            Address               `json:"miner"`
	LogsBloom        hexBloom              `json:"logsBloom"`
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
	Objects []Transaction
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

// FeeHistory represents the result of the feeHistory Client call.
type FeeHistory struct {
	OldestBlock   uint64
	Reward        [][]*big.Int
	BaseFeePerGas []*big.Int
	GasUsedRatio  []float64
}

func (f FeeHistory) MarshalJSON() ([]byte, error) {
	feeHistory := &jsonFeeHistory{
		OldestBlock:  Uint64ToNumber(f.OldestBlock),
		GasUsedRatio: f.GasUsedRatio,
	}
	if len(f.Reward) > 0 {
		feeHistory.Reward = make([][]Number, len(f.Reward))
		for i, reward := range f.Reward {
			feeHistory.Reward[i] = make([]Number, len(reward))
			for j, r := range reward {
				feeHistory.Reward[i][j] = BigIntToNumber(r)
			}
		}
	}
	if len(f.BaseFeePerGas) > 0 {
		feeHistory.BaseFeePerGas = make([]Number, len(f.BaseFeePerGas))
		for i, b := range f.BaseFeePerGas {
			feeHistory.BaseFeePerGas[i] = BigIntToNumber(b)
		}
	}
	return json.Marshal(feeHistory)
}

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
	Address     Address
	Topics      []Hash
	Data        []byte
	BlockHash   Hash
	BlockNumber uint64
	TxHash      Hash
	TxIndex     uint64
	LogIndex    uint64
	Removed     bool
}

func (l Log) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonLog{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockHash:   l.BlockHash,
		BlockNumber: Uint64ToNumber(l.BlockNumber),
		TxHash:      l.TxHash,
		TxIndex:     Uint64ToNumber(l.TxIndex),
		LogIndex:    Uint64ToNumber(l.LogIndex),
		Removed:     l.Removed,
	})
}

func (l *Log) UnmarshalJSON(input []byte) error {
	log := &jsonLog{}
	if err := json.Unmarshal(input, log); err != nil {
		return err
	}
	l.Address = log.Address
	l.Topics = log.Topics
	l.Data = log.Data
	l.BlockHash = log.BlockHash
	l.BlockNumber = log.BlockNumber.Big().Uint64()
	l.TxHash = log.TxHash
	l.TxIndex = log.TxIndex.Big().Uint64()
	l.LogIndex = log.LogIndex.Big().Uint64()
	l.Removed = log.Removed
	return nil
}

type jsonLog struct {
	Address     Address `json:"address"`
	Topics      []Hash  `json:"topics"`
	Data        Bytes   `json:"data"`
	BlockHash   Hash    `json:"blockHash"`
	BlockNumber Number  `json:"blockNumber"`
	TxHash      Hash    `json:"transactionHash"`
	TxIndex     Number  `json:"transactionIndex"`
	LogIndex    Number  `json:"logIndex"`
	Removed     bool    `json:"removed"`
}

// FilterLogsQuery represents a query to filter logs.
type FilterLogsQuery struct {
	Address   []Address
	FromBlock *BlockNumber
	ToBlock   *BlockNumber
	Topics    [][]Hash
	BlockHash *Hash
}

func (q FilterLogsQuery) MarshalJSON() ([]byte, error) {
	logsQuery := &jsonFilterLogsQuery{
		FromBlock: q.FromBlock,
		ToBlock:   q.ToBlock,
		BlockHash: q.BlockHash,
	}
	if len(q.Address) > 0 {
		logsQuery.Address = make([]Address, len(q.Address))
		for i, a := range q.Address {
			logsQuery.Address[i] = a
		}
	}
	if len(q.Topics) > 0 {
		logsQuery.Topics = make([]hashList, len(q.Topics))
		for i, t := range q.Topics {
			logsQuery.Topics[i] = make([]Hash, len(t))
			for j, h := range t {
				logsQuery.Topics[i][j] = h
			}
		}
	}
	return json.Marshal(logsQuery)
}

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
		for i, a := range logsQuery.Address {
			q.Address[i] = a
		}
	}
	if len(logsQuery.Topics) > 0 {
		q.Topics = make([][]Hash, len(logsQuery.Topics))
		for i, t := range logsQuery.Topics {
			q.Topics[i] = make([]Hash, len(t))
			for j, h := range t {
				q.Topics[i][j] = h
			}
		}
	}
	return nil
}

type jsonFilterLogsQuery struct {
	Address   addressList  `json:"address"`
	FromBlock *BlockNumber `json:"fromBlock,omitempty"`
	ToBlock   *BlockNumber `json:"toBlock,omitempty"`
	Topics    []hashList   `json:"topics"`
	BlockHash *Hash        `json:"blockhash,omitempty"`
}
