package app

import (
	"WalletRieltaTestTask/config"
	gateway "WalletRieltaTestTask/internal/wallet/gateway/rabbitmq"
	"WalletRieltaTestTask/internal/wallet/usecase"
	"WalletRieltaTestTask/pkg/postgres"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/client"
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
	walletUseCase := usecase.NewWallet(
		gateway.New(rmqClient),
		usecase.Timeout(cfg.App.Timeout),
		usecase.DefaultBalance(cfg.App.DefaultBalance),
	)

	return &App{
		DB: pg,
	}
}
