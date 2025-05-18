package rpc

import "github.com/defiweb/go-eth/types"

func getCallData(tx types.Transaction) *types.CallFields {
	if tx, ok := tx.(types.CallData); ok {
		return tx.CallData()
	}
	return nil
}

func getLegacyPriceData(tx types.Transaction) *types.LegacyPriceField {
	if tx, ok := tx.(types.LegacyPriceData); ok {
		return tx.LegacyPriceData()
	}
	return nil
}

func getAccessListData(tx types.Transaction) *types.AccessListField {
	if tx, ok := tx.(types.AccessListData); ok {
		return tx.AccessListData()
	}
	return nil
}

func getDynamicFeeData(tx types.Transaction) *types.DynamicFeeFields {
	if tx, ok := tx.(types.DynamicFeeData); ok {
		return tx.DynamicFeeData()
	}
	return nil
}

func convertTXToLegacyPrice(tx types.Transaction) types.Transaction {
	if getLegacyPriceData(tx) != nil {
		return tx
	}
	typ := types.LegacyTxType
	if getAccessListData(tx) != nil {
		typ = types.AccessListTxType
	}
	return convertTX(tx, typ)
}

func convertTXToDynamicFee(tx types.Transaction) types.Transaction {
	if getDynamicFeeData(tx) != nil {
		return tx
	}
	return convertTX(tx, types.DynamicFeeTxType)
}

func convertTX(tx types.Transaction, typ types.TransactionType) types.Transaction {
	switch typ {
	case types.LegacyTxType:
		ltx := types.NewTransactionLegacy()
		ltx.SetTransactionData(*tx.TransactionData())
		if tx, ok := tx.(types.CallData); ok {
			ltx.SetCallData(*tx.CallData())
		}
		if tx, ok := tx.(types.LegacyPriceData); ok {
			ltx.SetLegacyPriceData(*tx.LegacyPriceData())
		}
		return ltx
	case types.AccessListTxType:
		altx := types.NewTransactionAccessList()
		altx.SetTransactionData(*tx.TransactionData())
		if tx, ok := tx.(types.CallData); ok {
			altx.SetCallData(*tx.CallData())
		}
		if tx, ok := tx.(types.LegacyPriceData); ok {
			altx.SetLegacyPriceData(*tx.LegacyPriceData())
		}
		if tx, ok := tx.(types.AccessListData); ok {
			altx.SetAccessListData(*tx.AccessListData())
		}
		return altx
	case types.DynamicFeeTxType:
		dftx := types.NewTransactionDynamicFee()
		dftx.SetTransactionData(*tx.TransactionData())
		if tx, ok := tx.(types.CallData); ok {
			dftx.SetCallData(*tx.CallData())
		}
		if tx, ok := tx.(types.AccessListData); ok {
			dftx.SetAccessListData(*tx.AccessListData())
		}
		if tx, ok := tx.(types.DynamicFeeData); ok {
			dftx.SetDynamicFeeData(*tx.DynamicFeeData())
		}
		return dftx
	case types.BlobTxType:
		btx := types.NewTransactionBlob()
		btx.SetTransactionData(*tx.TransactionData())
		if tx, ok := tx.(types.CallData); ok {
			btx.SetCallData(*tx.CallData())
		}
		if tx, ok := tx.(types.AccessListData); ok {
			btx.SetAccessListData(*tx.AccessListData())
		}
		if tx, ok := tx.(types.DynamicFeeData); ok {
			btx.SetDynamicFeeData(*tx.DynamicFeeData())
		}
		if tx, ok := tx.(types.BlobData); ok {
			btx.SetBlobData(*tx.BlobData())
		}
		return btx
	default:
		return nil
	}
}
