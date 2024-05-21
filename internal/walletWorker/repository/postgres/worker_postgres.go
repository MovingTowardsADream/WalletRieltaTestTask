package worker_postgres

import (
	"WalletRieltaTestTask/internal/entity"
	"WalletRieltaTestTask/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	tableWallets      = "wallets"
	tableTransactions = "transactions"
)

type WalletRepo struct {
	db *postgres.Postgres
}

func New(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

// CreateNewWallet - creating new wallet entry in the db.
func (r *WalletRepo) CreateNewWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	sql, args, _ := r.db.Builder.
		Insert(tableWallets).
		Columns("balance").
		Values(wallet.Balance).
		Suffix("RETURNING id").
		ToSql()

	err := r.db.Pool.QueryRow(ctx, sql, args...).Scan(&wallet.ID)
	if err != nil {
		return wallet, fmt.Errorf("CreateNewWallet - r.Pool.QueryRow: %v", err)
	}

	return wallet, nil
}

// SendFunds - decreasing the balance of the sender and an increasing the receiver.
// Adding an entry to a transaction table.
func (r *WalletRepo) SendFunds(ctx context.Context, transaction *entity.Transaction) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo.SendFunds - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.db.Builder.
		Update(tableWallets).
		Set("balance", squirrel.Expr("balance - ?", transaction.Amount)).
		Where("id = ?", transaction.From).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo.SendFunds - tx.Exec: %v", err)
	}

	sql, args, _ = r.db.Builder.
		Update(tableWallets).
		Set("balance", squirrel.Expr("balance + ?", transaction.Amount)).
		Where("id = ?", transaction.To).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("WalletRepo.SendFunds - tx.Exec: %v", err)
	}

	sql, args, _ = r.db.Builder.
		Insert(tableTransactions).
		Columns("from_wallet_id", "to_wallet_id", "amount").
		Values(transaction.From, transaction.To, transaction.Amount).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo.SendFunds - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo.SendFunds - tx.Commit: %v", err)
	}

	return nil
}

// GetWalletHistoryByID - getting all transaction records from the user with the walletID.
func (r *WalletRepo) GetWalletHistoryByID(ctx context.Context, walletID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction

	sqlQuery, args, _ := r.db.Builder.
		Select("time, from_wallet_id, to_wallet_id, amount").
		From(tableTransactions).
		Where("from_wallet_id = ? OR to_wallet_id = ?", walletID, walletID).
		ToSql()

	rows, err := r.db.Pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo.GetWalletHistoryByID - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.Time, &transaction.From, &transaction.To, &transaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("OperationRepo.paginationOperationsByDate - rows.Scan: %v", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// GetWalletByID - getting wallet info by walletID.
func (r *WalletRepo) GetWalletByID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	sql, args, _ := r.db.Builder.
		Select("id, balance").
		From(tableWallets).
		Where("id = ?", walletID).
		ToSql()

	wallet := new(entity.Wallet)
	err := r.db.Pool.QueryRow(ctx, sql, args...).Scan(
		&wallet.ID,
		&wallet.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Not found")
		}
		return wallet, fmt.Errorf("WalletRepo.GetWalletByID - r.Pool.QueryRow: %v", err)
	}

	return wallet, nil
}
