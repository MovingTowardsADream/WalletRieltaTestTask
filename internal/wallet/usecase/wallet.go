package usecase

import (
	"WalletRieltaTestTask/internal/entity"
	"context"
	"fmt"
	"time"
)

const (
	_defaultTimeout      = 5 * time.Second
	_defaultBalance uint = 100
)

// WalletUseCase -.
type WalletUseCase struct {
	gateway        WalletGateway
	timeout        time.Duration
	defaultBalance uint
}

// New -.
func NewWallet(gw WalletGateway, opts ...Option) *WalletUseCase {
	uc := &WalletUseCase{
		gateway:        gw,
		timeout:        _defaultTimeout,
		defaultBalance: _defaultBalance,
	}

	for _, opt := range opts {
		opt(uc)
	}

	return uc
}

func (uc *WalletUseCase) CreateNewWalletWithDefaultBalance(ctx context.Context) (*entity.Wallet, error) {
	// Установка timeout на операцию
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	wallet, err := uc.gateway.CreateNewWalletWithBalance(ctxTimeout, uc.defaultBalance)
	if err != nil {
		return nil,
			fmt.Errorf("WalletUseCase - CreateNewWalletWithDefaultBalance - uc.gateway.CreateNewWalletWithBalance: %w", err)
	}

	return wallet, nil
}

func (uc *WalletUseCase) SendFunds(ctx context.Context, from string, to string, amount uint) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	if amount <= 0 {
		return entity.ErrWrongAmount
	}

	if len(from) == 0 || len(to) == 0 {
		return entity.ErrEmptyWallet
	}

	if from == to {
		return entity.ErrSenderIsReceiver
	}

	err := uc.gateway.SendFunds(ctxTimeout, from, to, amount)
	if err != nil {
		return fmt.Errorf("WalletUseCase - SendFunds - uc.gateway.SendFunds: %w", err)
	}

	return nil
}

func (uc *WalletUseCase) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	transactions, err := uc.gateway.GetWalletHistoryByID(ctxTimeout, walletID)
	if err != nil {
		return nil,
			fmt.Errorf("WalletUseCase - GetWalletHistoryByID - uc.gateway.GetWalletHistoryByID: %w", err)
	}

	return transactions, nil
}

func (uc *WalletUseCase) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	wallet, err := uc.gateway.GetWalletByID(ctxTimeout, walletID)
	if err != nil {
		return nil,
			fmt.Errorf("WalletUseCase - GetWalletByID - uc.gateway.GetWalletByID: %w", err)
	}

	return wallet, nil
}
