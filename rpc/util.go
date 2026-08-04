package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// convertTXToLegacyPrice converts a transaction to one that has legacy
// price data.
func convertTXToLegacyPrice(tx types.Transaction) types.Transaction {
	if types.GetLegacyFeeData(tx) != nil {
		return tx
	}
	typ := types.LegacyTxType
	if types.GetAccessListData(tx) != nil {
		typ = types.AccessListTxType
	}
	return convertTX(tx, typ)
}

// convertTXToAccessList converts a transaction to one that has access list
// data.
func convertTXToDynamicFee(tx types.Transaction) types.Transaction {
	if types.GetDynamicFeeData(tx) != nil {
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
		if tx, ok := tx.(types.HasSigningData); ok {
			ltx.SetSigningData(*tx.GetSigningData())
		}
		if tx, ok := tx.(types.HasExecutionData); ok {
			ltx.SetExecutionData(*tx.GetExecutionData())
		}
		if tx, ok := tx.(types.HasLegacyFeeData); ok {
			ltx.SetLegacyFeeData(*tx.GetLegacyFeeData())
		}
		return ltx
	case types.AccessListTxType:
		altx := types.NewTransactionAccessList()
		if tx, ok := tx.(types.HasSigningData); ok {
			altx.SetSigningData(*tx.GetSigningData())
		}
		if tx, ok := tx.(types.HasExecutionData); ok {
			altx.SetExecutionData(*tx.GetExecutionData())
		}
		if tx, ok := tx.(types.HasLegacyFeeData); ok {
			altx.SetLegacyFeeData(*tx.GetLegacyFeeData())
		}
		if tx, ok := tx.(types.HasAccessListData); ok {
			altx.SetAccessListData(*tx.GetAccessListData())
		}
		return altx
	case types.DynamicFeeTxType:
		dftx := types.NewTransactionDynamicFee()
		if tx, ok := tx.(types.HasSigningData); ok {
			dftx.SetSigningData(*tx.GetSigningData())
		}
		if tx, ok := tx.(types.HasExecutionData); ok {
			dftx.SetExecutionData(*tx.GetExecutionData())
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
		if tx, ok := tx.(types.HasSigningData); ok {
			btx.SetSigningData(*tx.GetSigningData())
		}
		if tx, ok := tx.(types.HasExecutionData); ok {
			btx.SetExecutionData(*tx.GetExecutionData())
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
		defer st.Unsubscribe(ctx, subID) //nolint:errcheck
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
				select {
				case msgCh <- msg:
				case <-ctx.Done():
					return
				}
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
