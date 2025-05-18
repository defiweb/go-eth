package rpc

import (
	"context"
	"errors"
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

// Call implements the transport.Hijacker interface.
func (h *hijackSimulate) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			switch method {
			case "eth_sendTransaction":
				// Verify arguments:
				if len(args) == 0 {
					return &ErrHijackFailed{name: "simulate", err: errors.New("missing transaction data")}
				}
				tx, ok := args[0].(types.Transaction)
				if !ok {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("invalid transaction type: %T", args[0])}
				}

				// Simulate transaction:
				if err := h.simulate(ctx, t, tx); err != nil {
					return &ErrHijackFailed{name: "simulate", err: err}
				}
			case "eth_sendRawTransaction", "eth_sendPrivateTransaction":
				// Verify arguments:
				if len(args) == 0 {
					return &ErrHijackFailed{name: "simulate", err: errors.New("missing transaction data")}
				}
				raw, ok := args[0].(types.Bytes)
				if !ok {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("invalid transaction data type: %T", args[0])}
				}

				// Decode raw transaction:
				tx, err := h.decoder.DecodeRLP(raw)
				if err != nil {
					return &ErrHijackFailed{name: "simulate", err: fmt.Errorf("failed to decode transaction: %w", err)}
				}

				// Simulate transaction:
				if err := h.simulate(ctx, t, tx); err != nil {
					return &ErrHijackFailed{name: "simulate", err: err}
				}
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the transport.Hijacker interface.
func (h *hijackSimulate) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the transport.Hijacker interface.
func (h *hijackSimulate) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

func (h *hijackSimulate) simulate(ctx context.Context, t transport.Transport, tx types.Transaction) error {
	// Recover transaction sender if not present:
	txd := tx.TransactionData()
	if txc, ok := tx.(types.CallData); ok && txd.Signature != nil && txc.CallData().From == nil {
		from, err := txsign.Recover(tx)
		if err != nil {
			return fmt.Errorf("unable to recover transaction sender: %w", err)
		}
		txc.CallData().From = from
	}

	// Get call data from transaction and execute it:
	call := tx.Call()
	if call == nil {
		return errors.New("unable to create call data from transaction")
	}
	if _, err := (&MethodsCommon{Transport: t}).Call(ctx, call, types.LatestBlockNumber); err != nil {
		return err
	}

	return nil
}
