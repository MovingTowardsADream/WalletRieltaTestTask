package entity

import (
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc"
	"context"
	"errors"
)

var (
	// Wallet errors.
	ErrWalletNotFound   = errors.New("wallet not found")
	ErrWrongAmount      = errors.New("wrong amount")
	ErrSenderIsReceiver = errors.New("sender is receiver")
	ErrEmptyWallet      = errors.New("wallet address is empty")

	// Requset errors.
	ErrTimeout  = context.DeadlineExceeded
	ErrNotFound = rmq_rpc.ErrNotFound
)
