package app

import (
	"WalletRieltaTestTask/config"
	v1 "WalletRieltaTestTask/internal/wallet/controller/http/v1"
	gateway "WalletRieltaTestTask/internal/wallet/gateway/rabbitmq"
	walletUseCase "WalletRieltaTestTask/internal/wallet/usecase"
	"WalletRieltaTestTask/internal/walletWorker/controller/amqp_rpc"
	worker_postgres "WalletRieltaTestTask/internal/walletWorker/repository/postgres"
	workerUC "WalletRieltaTestTask/internal/walletWorker/usecase"
	"WalletRieltaTestTask/pkg/httpserver"
	"WalletRieltaTestTask/pkg/postgres"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/client"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/server"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log/slog"
)

type App struct {
	HTTPServer *httpserver.Server
	RMQServer  *server.Server
	DB         *postgres.Postgres
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

	// Init http server
	handler := gin.New()
	v1.NewRouter(handler, log, walletUseCase)
	httpServer := httpserver.New(log, handler, httpserver.Port(cfg.HTTP.Port), httpserver.WriteTimeout(cfg.HTTP.Timeout))

	// Init rabbitMQ RPC Server
	rmqRouter := amqp_rpc.NewRouter(workerUseCase)

	rmqServer, err := server.New(
		cfg.RMQ.URL,
		cfg.RMQ.ServerExchange,
		rmqRouter,
		log,
		server.DefaultGoroutinesCount(cfg.App.CountWorkers),
	)
	if err != nil {
		panic("app - Run - rmqServer - server.New" + err.Error())
	}

	return &App{
		HTTPServer: httpServer,
		RMQServer:  rmqServer,
		DB:         pg,
	}
}
