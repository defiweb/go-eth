package rpc

import "fmt"

// ErrHijackFailed is returned when a hijacker fails to process a request.
type ErrHijackFailed struct {
	name string
	err  error
}

func (e *ErrHijackFailed) Name() string  { return e.name }
func (e *ErrHijackFailed) Error() string { return fmt.Sprintf("%s hijacker failed: %v", e.name, e.err) }
func (e *ErrHijackFailed) Unwrap() error { return e.err }
