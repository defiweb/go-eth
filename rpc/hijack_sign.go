package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

// hijackSign hijacks calls to methods that require account access
// and simulates their behavior using the provided keys.
type hijackSign struct {
	keys []wallet.Key
}

// Call implements the [transport.Hijacker] interface.
func (k *hijackSign) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) error {
			switch {
			case method == "eth_accounts":
				accounts, ok := result.(*[]types.Address)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid result type: %T", result)}
				}
				k.hijackAccountsCall(accounts)
			case method == "eth_sign" && len(args) == 2:
				signature, ok := result.(*types.Signature)
				if !ok {
					return fmt.Errorf("invalid result type: %T", result)
				}
				account, ok := args[0].(types.Address)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid result type: %T", args[0])}
				}
				data, ok := args[1].([]byte)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid data type: %T", args[1])}
				}
				return k.hijackSignCall(ctx, signature, account, data)
			case method == "eth_signTransaction" && len(args) == 1:
				raw, ok := result.(*[]byte)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid result type: %T", result)}
				}
				tx, ok := args[0].(types.SignableTransaction)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid transaction type: %T", args[0])}
				}
				return k.hijackSignTransactionCall(ctx, raw, tx)
			case method == "eth_sendTransaction" && len(args) == 1:
				hash, ok := result.(*types.Hash)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid result type: %T", args[0])}
				}
				tx, ok := args[0].(types.SignableTransaction)
				if !ok {
					return &ErrHijackFailed{name: "sign", err: fmt.Errorf("invalid transaction type: %T", args[0])}
				}
				return k.hijackSendTransactionCall(ctx, t, next, hash, tx)
			default:
				return next(ctx, t, result, method, args...)
			}
			return nil
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (k *hijackSign) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (k *hijackSign) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

func (k *hijackSign) hijackAccountsCall(result *[]types.Address) {
	*result = make([]types.Address, len(k.keys))
	for n, key := range k.keys {
		(*result)[n] = key.Address()
	}
}

func (k *hijackSign) hijackSignCall(ctx context.Context, result *types.Signature, account types.Address, data []byte) error {
	for _, key := range k.keys {
		if key.Address() != account {
			continue
		}
		signature, err := key.SignMessage(ctx, data)
		if err != nil {
			return err
		}
		*result = *signature
		return nil
	}
	return &ErrHijackFailed{name: "sign", err: fmt.Errorf("no key found for address %s", account)}
}

func (k *hijackSign) hijackSignTransactionCall(ctx context.Context, result *[]byte, tx types.SignableTransaction) error {
	if len(k.keys) == 0 {
		return fmt.Errorf("no keys found")
	}
	ed := types.GetExecutionData(tx)
	if ed == nil || ed.From == nil {
		return &ErrHijackFailed{name: "sign", err: fmt.Errorf("'from' field not set")}
	}
	for _, key := range k.keys {
		if key.Address() != *ed.From {
			continue
		}
		if err := key.SignTransaction(ctx, tx); err != nil {
			return &ErrHijackFailed{name: "sign", err: err}
		}
		raw, err := tx.EncodeRLP()
		if err != nil {
			return &ErrHijackFailed{name: "sign", err: err}
		}
		*result = raw
		return nil
	}
	return &ErrHijackFailed{name: "sign", err: fmt.Errorf("no key found for address %s", *ed.From)}
}

func (k *hijackSign) hijackSendTransactionCall(ctx context.Context, t transport.Transport, next transport.CallFunc, result *types.Hash, tx types.SignableTransaction) error {
	if len(k.keys) == 0 {
		return fmt.Errorf("no keys found")
	}
	ed := types.GetExecutionData(tx)
	if ed == nil || ed.From == nil {
		return &ErrHijackFailed{name: "sign", err: fmt.Errorf("'from' field not set")}
	}
	for _, key := range k.keys {
		if key.Address() != *ed.From {
			continue
		}
		if err := key.SignTransaction(ctx, tx); err != nil {
			return &ErrHijackFailed{name: "sign", err: err}
		}
		raw, err := tx.EncodeRLP()
		if err != nil {
			return &ErrHijackFailed{name: "sign", err: err}
		}
		return next(ctx, t, result, "eth_sendRawTransaction", types.Bytes(raw))
	}
	return &ErrHijackFailed{name: "sign", err: fmt.Errorf("no key found for address %s", *ed.From)}
}
