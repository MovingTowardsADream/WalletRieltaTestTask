package worker_postgres

import (
	"WalletRieltaTestTask/internal/entity"
	"WalletRieltaTestTask/pkg/postgres"
	"context"
)

type WalletRepo struct {
	db *postgres.Postgres
}

func New(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

// CreateNewWallet - creating new wallet entry in the db.
func (r *WalletRepo) CreateNewWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	// TODO

	return wallet, nil
}

// SendFunds - decreasing the balance of the sender and an increasing the receiver.
// Adding an entry to a transaction table.
func (r *WalletRepo) SendFunds(ctx context.Context, transaction *entity.Transaction) error {
	// TODO

	return nil
}

// GetWalletHistoryByID - getting all transaction records from the user with the walletID.
func (r *WalletRepo) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	// TODO

	return transactions, nil
}

// GetWalletByID - getting wallet info by walletID.
func (r *WalletRepo) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	wallet := new(entity.Wallet)

	// TODO

	return wallet, nil
}
