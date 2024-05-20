package usecase

import (
	"WalletRieltaTestTask/internal/entity"
	"context"
	"fmt"
)

type WalletWorkerUseCase struct {
	repo WalletWorkerRepo
}

func NewWalletWorker(r WalletWorkerRepo) *WalletWorkerUseCase {
	return &WalletWorkerUseCase{
		repo: r,
	}
}

// Creating a new wallet with balance in repository.
func (uc *WalletWorkerUseCase) CreateNewWalletWithBalance(ctx context.Context, balance uint) (*entity.Wallet, error) {
	// Create a new instance of the wallet with default balance
	defaultWallet := &entity.Wallet{
		Balance: balance,
	}

	wallet, err := uc.repo.CreateNewWallet(ctx, defaultWallet)
	if err != nil {
		return nil, fmt.Errorf("WalletWorkerUseCase - CreateNewWalletWithBalance - w.repo.CreateNewWallet: %w", err)
	}

	return wallet, nil
}

// Sending funds through wallets in repository.
func (uc *WalletWorkerUseCase) SendFunds(ctx context.Context, from string, to string, amount uint) error {
	transaction := &entity.Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}

	err := uc.repo.SendFunds(ctx, transaction)
	if err != nil {
		return fmt.Errorf("WalletWorkerUseCase - SendFunds - w.repo.SendFunds: %w", err)
	}

	return nil
}

// Getting wallet history by id from repository.
func (uc *WalletWorkerUseCase) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	transactions, err := uc.repo.GetWalletHistoryByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletWorkerUseCase - GetWalletHistoryByID - w.repo.GetWalletHistoryByID: %w", err)
	}

	return transactions, nil
}

// Getting wallet info by id from repository.
func (uc *WalletWorkerUseCase) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	wallet, err := uc.repo.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletWorkerUseCase - GetWalletByID - w.repo.GetWalletByID: %w", err)
	}

	return wallet, nil
}
