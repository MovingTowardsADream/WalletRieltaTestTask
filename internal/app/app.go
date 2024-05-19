package app

import (
	"WalletRieltaTestTask/config"
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

	fmt.Println(rmqClient)

	return &App{
		DB: pg,
	}
}
