package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// getTransactionData extracts the [types.TransactionData] from the given
// value.
func getTransactionData(v any) *types.TransactionData {
	if td, ok := v.(types.HasTransactionData); ok {
		return td.GetTransactionData()
	}
	return nil
}

// getCallData extracts the [types.CallData] from the given value.
func getCallData(v any) *types.CallData {
	if cd, ok := v.(types.HasCallData); ok {
		return cd.GetCallData()
	}
	return nil
}

// getLegacyPriceData extracts the [types.LegacyPriceData] from the given
// value.
func getLegacyPriceData(v any) *types.LegacyPriceData {
	if lpc, ok := v.(types.HasLegacyPriceData); ok {
		return lpc.GetLegacyPriceData()
	}
	return nil
}

// getAccessListData extracts the [types.AccessListData] from the given value.
func getAccessListData(v any) *types.AccessListData {
	if ald, ok := v.(types.HasAccessListData); ok {
		return ald.GetAccessListData()
	}
	return nil
}

// getDynamicFeeData extracts the [types.DynamicFeeData] from the given value.
func getDynamicFeeData(v any) *types.DynamicFeeData {
	if dfd, ok := v.(types.HasDynamicFeeData); ok {
		return dfd.GetDynamicFeeData()
	}
	return nil
}

// convertTXToLegacyPrice converts a transaction to one that has legacy
// price data.
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

// convertTXToAccessList converts a transaction to one that has access list
// data.
func convertTXToDynamicFee(tx types.Transaction) types.Transaction {
	if getDynamicFeeData(tx) != nil {
		return tx
	}
	return convertTX(tx, types.DynamicFeeTxType)
}

// convertTX converts a transaction to the specified type.
func convertTX(tx types.Transaction, typ types.TransactionType) types.Transaction {
	if tx.Type() == typ {
		return tx
	}
	switch typ {
	case types.LegacyTxType:
		ltx := types.NewTransactionLegacy()
		ltx.SetTransactionData(*tx.GetTransactionData())
		if tx, ok := tx.(types.HasCallData); ok {
			ltx.SetCallData(*tx.GetCallData())
		}
		if tx, ok := tx.(types.HasLegacyPriceData); ok {
			ltx.SetLegacyPriceData(*tx.GetLegacyPriceData())
		}
		return ltx
	case types.AccessListTxType:
		altx := types.NewTransactionAccessList()
		altx.SetTransactionData(*tx.GetTransactionData())
		if tx, ok := tx.(types.HasCallData); ok {
			altx.SetCallData(*tx.GetCallData())
		}
		if tx, ok := tx.(types.HasLegacyPriceData); ok {
			altx.SetLegacyPriceData(*tx.GetLegacyPriceData())
		}
		if tx, ok := tx.(types.HasAccessListData); ok {
			altx.SetAccessListData(*tx.GetAccessListData())
		}
		return altx
	case types.DynamicFeeTxType:
		dftx := types.NewTransactionDynamicFee()
		dftx.SetTransactionData(*tx.GetTransactionData())
		if tx, ok := tx.(types.HasCallData); ok {
			dftx.SetCallData(*tx.GetCallData())
		}
		if tx, ok := tx.(types.HasAccessListData); ok {
			dftx.SetAccessListData(*tx.GetAccessListData())
		}
		if tx, ok := tx.(types.HasDynamicFeeData); ok {
			dftx.SetDynamicFeeData(*tx.GetDynamicFeeData())
		}
		return dftx
	case types.BlobTxType:
		btx := types.NewTransactionBlob()
		btx.SetTransactionData(*tx.GetTransactionData())
		if tx, ok := tx.(types.HasCallData); ok {
			btx.SetCallData(*tx.GetCallData())
		}
		if tx, ok := tx.(types.HasAccessListData); ok {
			btx.SetAccessListData(*tx.GetAccessListData())
		}
		if tx, ok := tx.(types.HasDynamicFeeData); ok {
			btx.SetDynamicFeeData(*tx.GetDynamicFeeData())
		}
		if tx, ok := tx.(types.HasBlobData); ok {
			btx.SetBlobData(*tx.GetBlobData())
		}
		return btx
	default:
		return nil
	}
}

// subscribe creates a subscription to the given method and returns a channel
// that will receive the subscription messages. The messages are unmarshalled
// to the T type. The subscription is unsubscribed and channel closed when the
// context is cancelled.
func subscribe[T any](ctx context.Context, t transport.Transport, method string, params ...any) (chan T, error) {
	st, ok := t.(transport.SubscriptionTransport)
	if !ok {
		return nil, errors.New("transport does not support subscriptions")
	}
	rawCh, subID, err := st.Subscribe(ctx, method, params...)
	if err != nil {
		return nil, err
	}
	msgCh := make(chan T)
	go func() {
		defer close(msgCh)
		defer st.Unsubscribe(ctx, subID)
		for {
			select {
			case <-ctx.Done():
				return
			case raw, ok := <-rawCh:
				if !ok {
					return
				}
				var msg T
				if err := json.Unmarshal(raw, &msg); err != nil {
					continue
				}
				msgCh <- msg
			}
		}
	}()
	return msgCh, nil
}

// signTransactionResult is the result of an eth_signTransaction request.
// Some backends return only RLP encoded data, others return a JSON object,
// this type can handle both.
type signTransactionResult struct {
	Raw types.Bytes `json:"raw"`
}

func (s *signTransactionResult) UnmarshalJSON(input []byte) error {
	type alias signTransactionResult
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		return json.Unmarshal(input, &s.Raw)
	}
	var dec alias
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	s.Raw = dec.Raw
	return nil
}
