package gateway

import (
	"WalletRieltaTestTask/internal/entity"
	"context"
	"errors"
	"fmt"
)

type WalletGatewayRMQ interface {
	RemoteCall(ctx context.Context, handler string, request interface{}, response interface{}) error
}

type WalletGateway struct {
	rmq WalletGatewayRMQ
}

// Init of wallet gateway, through we will making requests to rmq server.
func New(rmq WalletGatewayRMQ) *WalletGateway {
	return &WalletGateway{rmq}
}

// Creating new wallet with balance, through remote call to rmq server.
func (gw *WalletGateway) CreateNewWalletWithBalance(ctx context.Context, balance uint) (*entity.Wallet, error) {
	var wallet entity.Wallet

	request := entity.CreateNewWalletWithBalanceRequest{
		Balance: balance,
	}

	err := wrapper(ctx, func() error {
		return gw.rmq.RemoteCall(ctx, "createNewWallet", request, &wallet)
	})

	if err != nil {
		return nil, fmt.Errorf("WalletGateway - CreateNewWalletWithBalance - gw.rmq.RemoteCall: %w", err)
	}

	return &wallet, nil
}

// Sending funds, through remote call to rmq server.
func (gw *WalletGateway) SendFunds(ctx context.Context, from string, to string, amount uint) error {
	request := entity.SendFundsRequest{
		From:   from,
		To:     to,
		Amount: amount,
	}

	err := wrapper(ctx, func() error {
		return gw.rmq.RemoteCall(ctx, "sendFunds", request, nil)
	})

	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.ErrWalletNotFound
		}

		return fmt.Errorf("WalletGateway - SendFunds - gw.rmq.RemoteCall: %w", err)
	}

	return nil
}

// Getting transactions history by wallet ID, through remote call to rmq server.
func (gw *WalletGateway) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	request := entity.GetWalletHistoryByIDRequest{
		WalletID: walletID,
	}

	err := wrapper(ctx, func() error {
		return gw.rmq.RemoteCall(ctx, "getWalletHistoryByID", request, &transactions)
	})

	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrWalletNotFound
		}

		return nil, fmt.Errorf("WalletGateway - GetWalletHistoryByID - gw.rmq.RemoteCall: %w", err)
	}

	return transactions, nil
}

// Getting wallet info by ID, through remote call to rmq server.
func (gw *WalletGateway) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	var wallet entity.Wallet

	request := entity.GetWalletByIDRequest{
		WalletID: walletID,
	}

	err := wrapper(ctx, func() error {
		return gw.rmq.RemoteCall(ctx, "getWalletByID", request, &wallet)
	})

	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrWalletNotFound
		}

		return nil, fmt.Errorf("WalletGateway - GetWalletByID - gw.rmq.RemoteCall: %w", err)
	}

	return &wallet, nil
}

// Эта функция используется для выполнения функции `f` в отдельной горутине
// и ожидания ответа или истечения таймаута, заданного контекстом `ctx`
func wrapper(ctx context.Context, f func() error) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err() //nolint:wrapcheck // we need just a send ctx error
	case err := <-errCh:
		return err
	}
}
