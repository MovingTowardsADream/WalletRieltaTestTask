package app

import (
	"WalletRieltaTestTask/config"
	gateway "WalletRieltaTestTask/internal/wallet/gateway/rabbitmq"
	walletUseCase "WalletRieltaTestTask/internal/wallet/usecase"
	worker_postgres "WalletRieltaTestTask/internal/walletWorker/repository/postgres"
	workerUC "WalletRieltaTestTask/internal/walletWorker/usecase"
	"WalletRieltaTestTask/pkg/postgres"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/client"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
)

type App struct {
	DB *postgres.Postgres
}

func New(log *slog.Logger, cfg *config.Config) *App {
	// Connect postgres db
	pg, err := postgres.NewPostgresDB(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		panic("app - Run - postgres.NewPostgresDB: " + err.Error())
	}

	rmqClient, err := client.NewRabbitMQClient(cfg.RMQ.URL, cfg.RMQ.ServerExchange, cfg.RMQ.ClientExchange)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	// Use cases
	walletUseCase := walletUseCase.NewWallet(
		gateway.New(rmqClient),
		walletUseCase.Timeout(cfg.App.Timeout),
		walletUseCase.DefaultBalance(cfg.App.DefaultBalance),
	)

	workerUseCase := workerUC.NewWalletWorker(
		worker_postgres.New(pg),
	)

	fmt.Println(walletUseCase, workerUseCase)

	return &App{
		DB: pg,
	}
}
