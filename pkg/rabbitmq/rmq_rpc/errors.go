package rmq_rpc

import "errors"

var (
	ErrTimeout = errors.New("timeout")

	ErrBadHandler = errors.New("unregistered handler")

	ErrContextCanceled = errors.New("context canceled")

	ErrNotFound = errors.New("not found")

	ErrCallStatus = errors.New("call status")
)
