package amqp_rpc

import (
	"WalletRieltaTestTask/internal/entity"
	"WalletRieltaTestTask/internal/walletWorker/usecase"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/server"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

type walletWorkerRoutes struct {
	w usecase.WalletWorker
}

// Вeclaring routes for rmq rpc.
func newWalletWorkerRoutes(routes map[string]server.CallHandler, w usecase.WalletWorker) {
	r := &walletWorkerRoutes{w}
	{
		routes["createNewWallet"] = r.createNewWalletWithBalance()
		routes["sendFunds"] = r.sendFunds()
		routes["getWalletHistoryByID"] = r.getWalletHistoryByID()
		routes["getWalletByID"] = r.getWalletByID()
	}
}

// Handles a remote "createNewWallet" call.
func (r *walletWorkerRoutes) createNewWalletWithBalance() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.CreateNewWalletWithBalanceRequest

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - createNewWalletWithBalance - json.Unmarshal: %w", err)
		}

		wallet, err := r.w.CreateNewWalletWithBalance(context.Background(), request.Balance)
		if err != nil {
			return nil,
				fmt.Errorf("amqp_rpc - walletWorkerRoutes - createNewWalletWithBalance - r.w.CreateNewWalletWithBalance: %w", err)
		}

		return wallet, nil
	}
}

// Handles a remote "sendFunds" call.
func (r *walletWorkerRoutes) sendFunds() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.SendFundsRequest

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - sendFunds - json.Unmarshal: %w", err)
		}

		err := r.w.SendFunds(context.Background(), request.From, request.To, request.Amount)
		if err != nil {
			if errors.Is(err, entity.ErrWalletNotFound) {
				return nil, entity.ErrNotFound
			}

			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - sendFunds - r.w.SendFunds: %w", err)
		}

		return nil, nil
	}
}

// Handles a remote "getWalletHistoryByID" call.
func (r *walletWorkerRoutes) getWalletHistoryByID() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.GetWalletHistoryByIDRequest

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - GetWalletHistoryByID - json.Unmarshal: %w", err)
		}

		transactions, err := r.w.GetWalletHistoryByID(context.Background(), request.WalletID)
		if err != nil {
			if errors.Is(err, entity.ErrWalletNotFound) {
				return nil, entity.ErrNotFound
			}

			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - GetWalletHistoryByID - r.w.GetWalletHistoryByID: %w", err)
		}

		return transactions, nil
	}
}

// Handles a remote "getWalletByID" call.
func (r *walletWorkerRoutes) getWalletByID() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		var request entity.GetWalletByIDRequest

		if err := json.Unmarshal(d.Body, &request); err != nil {
			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - GetWalletByID - json.Unmarshal: %w", err)
		}

		wallet, err := r.w.GetWalletByID(context.Background(), request.WalletID)
		if err != nil {
			if errors.Is(err, entity.ErrWalletNotFound) {
				return nil, entity.ErrNotFound
			}

			return nil, fmt.Errorf("amqp_rpc - walletWorkerRoutes - GetWalletByID - r.w.GetWalletByID: %w", err)
		}

		return wallet, nil
	}
}
