package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/crypto/txsign"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackSimulate hijacks "eth_send*Transaction" methods and simulates the
// transaction execution before sending it.
type hijackSimulate struct {
	decoder types.RLPTransactionDecoder
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackSimulate) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}
			switch method {
			case "eth_sendTransaction":
				tx, ok := args[0].(types.Transaction)
				if !ok {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("no transaction found")}
				}
				if err := h.simulate(ctx, t, tx); err != nil {
					return &ErrHijackFailed{name: "simulate", err: err}
				}
			case "eth_sendRawTransaction", "eth_sendPrivateTransaction":
				raw, ok := args[0].(types.Bytes)
				if !ok {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("no raw transaction found")}
				}
				tx, err := h.decoder.DecodeRLP(raw)
				if err != nil {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("failed to decode transaction: %w", err)}
				}
				if err := h.simulate(ctx, t, tx); err != nil {
					return &ErrHijackFailed{name: "simulate", err: err}
				}
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackSimulate) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackSimulate) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

func (h *hijackSimulate) simulate(ctx context.Context, t transport.Transport, tx types.Transaction) error {
	// If the transaction has not a corrseponding call, we cannot simulate it.
	// This can happen for custom transactions types.
	call := tx.Call()
	if call == nil {
		return nil
	}

	// Recover the transaction sender if it is not present.
	// This can happen if the transaction is encoded using RLP, as this format
	// contains only the signature and not the sender address.
	if td := tx.GetTransactionData(); td.Signature != nil && call.GetCallData().From == nil {
		from, err := txsign.Recover(tx)
		if err != nil {
			return fmt.Errorf("unable to recover transaction sender: %w", err)
		}
		call.GetCallData().From = from
		if cd := getCallData(tx); cd != nil {
			cd.From = from
		}
	}

	// Simulate the transaction.
	if _, err := (&MethodsCommon{Transport: t}).Call(ctx, call, types.LatestBlockNumber); err != nil {
		return err
	}
	return nil
}
