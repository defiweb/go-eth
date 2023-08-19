package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/defiweb/go-rlp"
)

// Call represents a call to a contract or a contract creation if To is nil.
type Call struct {
	From     *Address // From is the sender address.
	To       *Address // To is the recipient address. nil means contract creation.
	GasLimit *uint64  // GasLimit is the gas limit, if 0, there is no limit.
	GasPrice *big.Int // GasPrice is the gas price in wei per gas unit.
	Value    *big.Int // Value is the amount of wei to send.
	Input    []byte   // Input is the input data.

	// EIP-2930 fields:
	AccessList AccessList // AccessList is the list of addresses and storage keys that the transaction can access.

	// EIP-1559 fields:
	MaxPriorityFeePerGas *big.Int // MaxPriorityFeePerGas is the maximum priority fee per gas the sender is willing to pay.
	MaxFeePerGas         *big.Int // MaxFeePerGas is the maximum fee per gas the sender is willing to pay.
}

func (c *Call) SetFrom(from Address) *Call {
	c.From = &from
	return c
}

func (c *Call) SetTo(to Address) *Call {
	c.To = &to
	return c
}

func (c *Call) SetGasLimit(gasLimit uint64) *Call {
	c.GasLimit = &gasLimit
	return c
}

func (c *Call) SetGasPrice(gasPrice *big.Int) *Call {
	c.GasPrice = gasPrice
	return c
}

func (c *Call) SetValue(value *big.Int) *Call {
	c.Value = value
	return c
}

func (c *Call) SetInput(input []byte) *Call {
	c.Input = input
	return c
}

func (c *Call) SetAccessList(accessList AccessList) *Call {
	c.AccessList = accessList
	return c
}

func (c *Call) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) *Call {
	c.MaxPriorityFeePerGas = maxPriorityFeePerGas
	return c
}

func (c *Call) SetMaxFeePerGas(maxFeePerGas *big.Int) *Call {
	c.MaxFeePerGas = maxFeePerGas
	return c
}

func (c Call) MarshalJSON() ([]byte, error) {
	call := &jsonCall{
		From:       c.From,
		To:         c.To,
		Data:       c.Input,
		AccessList: c.AccessList,
	}
	if c.GasLimit != nil {
		call.GasLimit = NumberFromUint64Ptr(*c.GasLimit)
	}
	if c.GasPrice != nil {
		call.GasPrice = NumberFromBigIntPtr(c.GasPrice)
	}
	if c.MaxFeePerGas != nil {
		call.MaxFeePerGas = NumberFromBigIntPtr(c.MaxFeePerGas)
	}
	if c.MaxPriorityFeePerGas != nil {
		call.MaxPriorityFeePerGas = NumberFromBigIntPtr(c.MaxPriorityFeePerGas)
	}
	if c.Value != nil {
		value := NumberFromBigInt(c.Value)
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
	if call.GasLimit != nil {
		gas := call.GasLimit.Big().Uint64()
		c.GasLimit = &gas
	}
	if call.GasPrice != nil {
		c.GasPrice = call.GasPrice.Big()
	}
	if call.MaxFeePerGas != nil {
		c.MaxFeePerGas = call.MaxFeePerGas.Big()
	}
	if call.MaxPriorityFeePerGas != nil {
		c.MaxPriorityFeePerGas = call.MaxPriorityFeePerGas.Big()
	}
	if call.Value != nil {
		c.Value = call.Value.Big()
	}
	c.Input = call.Data
	c.AccessList = call.AccessList
	return nil
}

type jsonCall struct {
	From                 *Address   `json:"from,omitempty"`
	To                   *Address   `json:"to,omitempty"`
	GasLimit             *Number    `json:"gas,omitempty"`
	GasPrice             *Number    `json:"gasPrice,omitempty"`
	MaxFeePerGas         *Number    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *Number    `json:"maxPriorityFeePerGas,omitempty"`
	Value                *Number    `json:"value,omitempty"`
	Data                 Bytes      `json:"data,omitempty"`
	AccessList           AccessList `json:"accessList,omitempty"`
}

// TransactionType is the type of transaction.
type TransactionType uint64

// Transaction types.
const (
	LegacyTxType TransactionType = iota
	AccessListTxType
	DynamicFeeTxType
)

// Transaction represents a transaction.
type Transaction struct {
	Call

	// Transaction data:
	Type      TransactionType // Type is the transaction type.
	Nonce     *uint64         // Nonce is the number of transactions made by the sender prior to this one.
	Signature *Signature      // Signature of the transaction.

	// EIP-2930 fields:
	ChainID *uint64 // ChainID is the chain ID of the transaction.
}

func (t *Transaction) SetFrom(from Address) *Transaction {
	t.From = &from
	return t
}

func (t *Transaction) SetTo(to Address) *Transaction {
	t.To = &to
	return t
}

func (t *Transaction) SetGasLimit(gasLimit uint64) *Transaction {
	t.GasLimit = &gasLimit
	return t
}

func (t *Transaction) SetGasPrice(gasPrice *big.Int) *Transaction {
	t.GasPrice = gasPrice
	return t
}

func (t *Transaction) SetValue(value *big.Int) *Transaction {
	t.Value = value
	return t
}

func (t *Transaction) SetInput(input []byte) *Transaction {
	t.Input = input
	return t
}

func (t *Transaction) SetAccessList(accessList AccessList) *Transaction {
	t.AccessList = accessList
	return t
}

func (t *Transaction) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) *Transaction {
	t.MaxPriorityFeePerGas = maxPriorityFeePerGas
	return t
}

func (t *Transaction) SetMaxFeePerGas(maxFeePerGas *big.Int) *Transaction {
	t.MaxFeePerGas = maxFeePerGas
	return t
}

func (t *Transaction) SetType(transactionType TransactionType) *Transaction {
	t.Type = transactionType
	return t
}

func (t *Transaction) SetNonce(nonce uint64) *Transaction {
	t.Nonce = &nonce
	return t
}

func (t *Transaction) SetSignature(signature Signature) *Transaction {
	t.Signature = &signature
	return t
}

func (t *Transaction) SetChainID(chainID uint64) *Transaction {
	t.ChainID = &chainID
	return t
}

// Raw returns the raw transaction data that could be sent to the network.
func (t Transaction) Raw() ([]byte, error) {
	return t.EncodeRLP()
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	transaction := &jsonTransaction{}
	transaction.To = t.To
	transaction.From = t.From
	if t.GasLimit != nil {
		transaction.GasLimit = NumberFromUint64Ptr(*t.GasLimit)
	}
	if t.GasPrice != nil {
		transaction.GasPrice = NumberFromBigIntPtr(t.GasPrice)
	}
	if t.MaxFeePerGas != nil {
		transaction.MaxFeePerGas = NumberFromBigIntPtr(t.MaxFeePerGas)
	}
	if t.MaxPriorityFeePerGas != nil {
		transaction.MaxPriorityFeePerGas = NumberFromBigIntPtr(t.MaxPriorityFeePerGas)
	}
	transaction.Input = t.Input
	if t.Nonce != nil {
		transaction.Nonce = NumberFromUint64Ptr(*t.Nonce)
	}
	if t.Value != nil {
		transaction.Value = NumberFromBigIntPtr(t.Value)
	}
	transaction.AccessList = t.AccessList
	if t.Signature != nil {
		transaction.V = NumberFromBigIntPtr(t.Signature.V)
		transaction.R = NumberFromBigIntPtr(t.Signature.R)
		transaction.S = NumberFromBigIntPtr(t.Signature.S)
	}
	return json.Marshal(transaction)
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	transaction := &jsonTransaction{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return err
	}
	t.To = transaction.To
	t.From = transaction.From
	if transaction.GasLimit != nil {
		gas := transaction.GasLimit.Big().Uint64()
		t.GasLimit = &gas
	}
	if transaction.GasPrice != nil {
		t.GasPrice = transaction.GasPrice.Big()
	}
	if transaction.MaxFeePerGas != nil {
		t.MaxFeePerGas = transaction.MaxFeePerGas.Big()
	}
	if transaction.MaxPriorityFeePerGas != nil {
		t.MaxPriorityFeePerGas = transaction.MaxPriorityFeePerGas.Big()
	}
	t.Input = transaction.Input
	if transaction.Nonce != nil {
		nonce := transaction.Nonce.Big().Uint64()
		t.Nonce = &nonce
	}
	if transaction.Value != nil {
		t.Value = transaction.Value.Big()
	}
	t.AccessList = transaction.AccessList
	if transaction.V != nil && transaction.R != nil && transaction.S != nil {
		t.Signature = SignatureFromVRSPtr(transaction.V.Big(), transaction.R.Big(), transaction.S.Big())
	}
	return nil
}

//nolint:funlen
func (t Transaction) EncodeRLP() ([]byte, error) {
	var (
		chainID              = uint64(1)
		nonce                = uint64(0)
		gasPrice             = big.NewInt(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (AccessList)(nil)
		v                    = big.NewInt(0)
		r                    = big.NewInt(0)
		s                    = big.NewInt(0)
	)
	if t.ChainID != nil {
		chainID = *t.ChainID
	}
	if t.Nonce != nil {
		nonce = *t.Nonce
	}
	if t.GasPrice != nil {
		gasPrice = t.GasPrice
	}
	if t.GasLimit != nil {
		gasLimit = *t.GasLimit
	}
	if t.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = t.MaxPriorityFeePerGas
	}
	if t.MaxFeePerGas != nil {
		maxFeePerGas = t.MaxFeePerGas
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = t.Value
	}
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	if t.Signature != nil {
		v = t.Signature.V
		r = t.Signature.R
		s = t.Signature.S
	}
	switch t.Type {
	case LegacyTxType:
		return rlp.NewList(
			rlp.NewUint(nonce),
			rlp.NewBigInt(gasPrice),
			rlp.NewUint(gasLimit),
			rlp.NewBytes(to),
			rlp.NewBigInt(value),
			rlp.NewBytes(t.Input),
			rlp.NewBigInt(v),
			rlp.NewBigInt(r),
			rlp.NewBigInt(s),
		).EncodeRLP()
	case AccessListTxType:
		bin, err := rlp.NewList(
			rlp.NewUint(chainID),
			rlp.NewUint(nonce),
			rlp.NewBigInt(gasPrice),
			rlp.NewUint(gasLimit),
			rlp.NewBytes(to),
			rlp.NewBigInt(value),
			rlp.NewBytes(t.Input),
			&t.AccessList,
			rlp.NewBigInt(v),
			rlp.NewBigInt(r),
			rlp.NewBigInt(s),
		).EncodeRLP()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(t.Type)}, bin...), nil
	case DynamicFeeTxType:
		bin, err := rlp.NewList(
			rlp.NewUint(chainID),
			rlp.NewUint(nonce),
			rlp.NewBigInt(maxPriorityFeePerGas),
			rlp.NewBigInt(maxFeePerGas),
			rlp.NewUint(gasLimit),
			rlp.NewBytes(to),
			rlp.NewBigInt(value),
			rlp.NewBytes(t.Input),
			&accessList,
			rlp.NewBigInt(v),
			rlp.NewBigInt(r),
			rlp.NewBigInt(s),
		).EncodeRLP()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(t.Type)}, bin...), nil
	default:
		return nil, fmt.Errorf("unknown transaction type: %d", t.Type)
	}
}

//nolint:funlen
func (t *Transaction) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	var (
		list                 = (*rlp.ListItem)(nil)
		chainID              = &rlp.UintItem{}
		nonce                = &rlp.UintItem{}
		gasPrice             = &rlp.BigIntItem{}
		gasLimit             = &rlp.UintItem{}
		maxPriorityFeePerGas = &rlp.BigIntItem{}
		maxFeePerGas         = &rlp.BigIntItem{}
		to                   = &rlp.StringItem{}
		value                = &rlp.BigIntItem{}
		input                = &rlp.StringItem{}
		accessList           = &AccessList{}
		v                    = &rlp.BigIntItem{}
		r                    = &rlp.BigIntItem{}
		s                    = &rlp.BigIntItem{}
	)
	switch {
	case data[0] >= 0x80: // LegacyTxType
		t.Type = LegacyTxType
		list = rlp.NewList(
			nonce,
			gasPrice,
			gasLimit,
			to,
			value,
			input,
			v,
			r,
			s,
		)
	case data[0] == byte(AccessListTxType):
		t.Type = AccessListTxType
		data = data[1:]
		list = rlp.NewList(
			chainID,
			nonce,
			gasPrice,
			gasLimit,
			to,
			value,
			input,
			accessList,
			v,
			r,
			s,
		)
	case data[0] == byte(DynamicFeeTxType):
		t.Type = DynamicFeeTxType
		data = data[1:]
		list = rlp.NewList(
			chainID,
			nonce,
			maxPriorityFeePerGas,
			maxFeePerGas,
			gasLimit,
			to,
			value,
			input,
			accessList,
			v,
			r,
			s,
		)
	default:
		return 0, fmt.Errorf("invalid transaction type: %d", data[0])
	}
	if _, err := rlp.DecodeTo(data, list); err != nil {
		return 0, err
	}
	t.ChainID = &chainID.X
	t.Nonce = &nonce.X
	t.GasPrice = gasPrice.X
	t.GasLimit = &gasLimit.X
	t.MaxPriorityFeePerGas = maxPriorityFeePerGas.X
	t.MaxFeePerGas = maxFeePerGas.X
	t.To = AddressFromBytesPtr(to.Bytes())
	t.Value = value.X
	if len(input.Bytes()) > 0 {
		t.Input = input.Bytes()
	}
	if len(*accessList) > 0 {
		t.AccessList = *accessList
	}
	if v.X.Sign() != 0 || r.X.Sign() != 0 || s.X.Sign() != 0 {
		t.Signature = &Signature{
			V: v.X,
			R: r.X,
			S: s.X,
		}
	}
	return len(data), nil
}

// Hash returns the hash of the transaction (transaction ID).
func (t Transaction) Hash(h HashFunc) (Hash, error) {
	raw, err := t.Raw()
	if err != nil {
		return Hash{}, err
	}
	return h(raw), nil
}

type jsonTransaction struct {
	From                 *Address   `json:"from,omitempty,omitempty"`
	To                   *Address   `json:"to,omitempty,omitempty"`
	GasLimit             *Number    `json:"gas,omitempty"`
	GasPrice             *Number    `json:"gasPrice,omitempty"`
	MaxFeePerGas         *Number    `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *Number    `json:"maxPriorityFeePerGas,omitempty"`
	Input                Bytes      `json:"input,omitempty"`
	Nonce                *Number    `json:"nonce,omitempty"`
	Value                *Number    `json:"value,omitempty"`
	AccessList           AccessList `json:"accessList,omitempty"`
	V                    *Number    `json:"v,omitempty"`
	R                    *Number    `json:"r,omitempty"`
	S                    *Number    `json:"s,omitempty"`
}

// OnChainTransaction represents a transaction that is included in a block.
type OnChainTransaction struct {
	Transaction

	// On-chain fields, only available when the transaction is included in a block:
	Hash             *Hash    // Hash of the transaction.
	BlockHash        *Hash    // BlockHash is the hash of the block where this transaction was in.
	BlockNumber      *big.Int // BlockNumber is the block number where this transaction was in.
	TransactionIndex *uint64  // TransactionIndex is the index of the transaction in the block.
}

type jsonOnChainTransaction struct {
	jsonTransaction
	Hash             *Hash   `json:"hash,omitempty"`
	BlockHash        *Hash   `json:"blockHash,omitempty"`
	BlockNumber      *Number `json:"blockNumber,omitempty"`
	TransactionIndex *Number `json:"transactionIndex,omitempty"`
}

func (t OnChainTransaction) MarshalJSON() ([]byte, error) {
	transaction := &jsonOnChainTransaction{}
	transaction.To = t.To
	transaction.From = t.From
	if t.GasLimit != nil {
		transaction.GasLimit = NumberFromUint64Ptr(*t.GasLimit)
	}
	if t.GasPrice != nil {
		transaction.GasPrice = NumberFromBigIntPtr(t.GasPrice)
	}
	if t.MaxFeePerGas != nil {
		transaction.MaxFeePerGas = NumberFromBigIntPtr(t.MaxFeePerGas)
	}
	if t.MaxPriorityFeePerGas != nil {
		transaction.MaxPriorityFeePerGas = NumberFromBigIntPtr(t.MaxPriorityFeePerGas)
	}
	transaction.Input = t.Input
	if t.Nonce != nil {
		transaction.Nonce = NumberFromUint64Ptr(*t.Nonce)
	}
	if t.Value != nil {
		transaction.Value = NumberFromBigIntPtr(t.Value)
	}
	transaction.AccessList = t.AccessList
	if t.Signature != nil {
		transaction.V = NumberFromBigIntPtr(t.Signature.V)
		transaction.R = NumberFromBigIntPtr(t.Signature.R)
		transaction.S = NumberFromBigIntPtr(t.Signature.S)
	}
	transaction.Hash = t.Hash
	transaction.BlockHash = t.BlockHash
	if t.BlockNumber != nil {
		transaction.BlockNumber = NumberFromBigIntPtr(t.BlockNumber)
	}
	if t.TransactionIndex != nil {
		transaction.TransactionIndex = NumberFromUint64Ptr(*t.TransactionIndex)
	}
	return json.Marshal(transaction)
}

func (t *OnChainTransaction) UnmarshalJSON(data []byte) error {
	transaction := &jsonOnChainTransaction{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return err
	}
	t.To = transaction.To
	t.From = transaction.From
	if transaction.GasLimit != nil {
		gas := transaction.GasLimit.Big().Uint64()
		t.GasLimit = &gas
	}
	if transaction.GasPrice != nil {
		t.GasPrice = transaction.GasPrice.Big()
	}
	if transaction.MaxFeePerGas != nil {
		t.MaxFeePerGas = transaction.MaxFeePerGas.Big()
	}
	if transaction.MaxPriorityFeePerGas != nil {
		t.MaxPriorityFeePerGas = transaction.MaxPriorityFeePerGas.Big()
	}
	t.Input = transaction.Input
	if transaction.Nonce != nil {
		nonce := transaction.Nonce.Big().Uint64()
		t.Nonce = &nonce
	}
	if transaction.Value != nil {
		t.Value = transaction.Value.Big()
	}
	t.AccessList = transaction.AccessList
	if transaction.V != nil && transaction.R != nil && transaction.S != nil {
		t.Signature = SignatureFromVRSPtr(transaction.V.Big(), transaction.R.Big(), transaction.S.Big())
	}
	t.Hash = transaction.Hash
	t.BlockHash = transaction.BlockHash
	if transaction.BlockNumber != nil {
		t.BlockNumber = transaction.BlockNumber.Big()
	}
	if transaction.TransactionIndex != nil {
		index := transaction.TransactionIndex.Big().Uint64()
		t.TransactionIndex = &index
	}
	return nil
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
		tuple := tuple // Copy value because of loop variable reuse.
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
		if err := tuple.DecodeTo(&t); err != nil {
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
	if err := l[0].DecodeTo(&a.Address); err != nil {
		return n, err
	}
	h, err := l[1].GetList()
	if err != nil {
		return n, err
	}
	for _, item := range h {
		var hash Hash
		if err := item.DecodeTo(&hash); err != nil {
			return n, err
		}
		a.StorageKeys = append(a.StorageKeys, hash)
	}
	return n, nil
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
	Transactions      []OnChainTransaction // Transactions is the list of transactions in the block.
	TransactionHashes []Hash               // TransactionHashes is the list of transaction hashes in the block.
	ExtraData         []byte               // ExtraData is the "extra data" field of this block.
}

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
	Objects []OnChainTransaction
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
	OldestBlock   uint64       // OldestBlock is the oldest block number for which the base fee and gas used are returned.
	Reward        [][]*big.Int // Reward is the reward for each block in the range [OldestBlock, LatestBlock].
	BaseFeePerGas []*big.Int   // BaseFeePerGas is the base fee per gas for each block in the range [OldestBlock, LatestBlock].
	GasUsedRatio  []float64    // GasUsedRatio is the gas used ratio for each block in the range [OldestBlock, LatestBlock].
}

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

func (q *FilterLogsQuery) SetAddresses(addresses []Address) *FilterLogsQuery {
	q.Address = addresses
	return q
}

func (q *FilterLogsQuery) AddAddresses(addresses ...Address) *FilterLogsQuery {
	q.Address = append(q.Address, addresses...)
	return q
}

func (q *FilterLogsQuery) SetFromBlock(fromBlock *BlockNumber) *FilterLogsQuery {
	q.FromBlock = fromBlock
	return q
}

func (q *FilterLogsQuery) SetToBlock(toBlock *BlockNumber) *FilterLogsQuery {
	q.ToBlock = toBlock
	return q
}

func (q *FilterLogsQuery) SetTopics(topics [][]Hash) *FilterLogsQuery {
	q.Topics = topics
	return q
}

func (q *FilterLogsQuery) AddTopics(topics ...[]Hash) *FilterLogsQuery {
	q.Topics = append(q.Topics, topics...)
	return q
}

func (q *FilterLogsQuery) SetBlockHash(blockHash *Hash) *FilterLogsQuery {
	q.BlockHash = blockHash
	return q
}

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
		logsQuery.Topics = make([]hashList, len(q.Topics))
		for i, t := range q.Topics {
			logsQuery.Topics[i] = make([]Hash, len(t))
			copy(logsQuery.Topics[i], t)
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
			copy(q.Topics[i], t)
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
