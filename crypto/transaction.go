package crypto

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/types"
)

func signingHash(t *types.Transaction) (types.Hash, error) {
	var (
		chainID              = uint64(1)
		nonce                = uint64(0)
		gasPrice             = big.NewInt(0)
		gasLimit             = uint64(0)
		maxPriorityFeePerGas = big.NewInt(0)
		maxFeePerGas         = big.NewInt(0)
		to                   = ([]byte)(nil)
		value                = big.NewInt(0)
		accessList           = (types.AccessList)(nil)
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
	switch t.Type {
	case types.LegacyTxType:
		list := rlp.NewList(
			rlp.NewUint(nonce),
			rlp.NewBigInt(gasPrice),
			rlp.NewUint(gasLimit),
			rlp.NewBytes(to),
			rlp.NewBigInt(value),
			rlp.NewBytes(t.Input),
		)
		if t.ChainID != nil && *t.ChainID != 0 {
			list.Append(
				rlp.NewUint(chainID),
				rlp.NewUint(0),
				rlp.NewUint(0),
			)
		}
		bin, err := list.EncodeRLP()
		if err != nil {
			return types.Hash{}, err
		}
		return Keccak256(bin), nil
	case types.AccessListTxType:
		bin, err := rlp.NewList(
			rlp.NewUint(chainID),
			rlp.NewUint(nonce),
			rlp.NewBigInt(gasPrice),
			rlp.NewUint(gasLimit),
			rlp.NewBytes(to),
			rlp.NewBigInt(value),
			rlp.NewBytes(t.Input),
			&t.AccessList,
		).EncodeRLP()
		if err != nil {
			return types.Hash{}, err
		}
		bin = append([]byte{byte(t.Type)}, bin...)
		return Keccak256(bin), nil
	case types.DynamicFeeTxType:
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
		).EncodeRLP()
		if err != nil {
			return types.Hash{}, err
		}
		bin = append([]byte{byte(t.Type)}, bin...)
		return Keccak256(bin), nil
	default:
		return types.Hash{}, fmt.Errorf("invalid transaction type: %d", t.Type)
	}
}
