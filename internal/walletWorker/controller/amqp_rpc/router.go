package amqp_rpc

import (
	"WalletRieltaTestTask/internal/walletWorker/usecase"
	"WalletRieltaTestTask/pkg/rabbitmq/rmq_rpc/server"
)

func NewRouter(r usecase.WalletWorker) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newWalletWorkerRoutes(routes, r)
	}

	return routes
}
