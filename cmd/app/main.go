package main

import (
	"WalletRieltaTestTask/config"
	"WalletRieltaTestTask/internal/app"
	"WalletRieltaTestTask/pkg/logger"
	"fmt"
)

func main() {
	// Init configuration
	cfg := config.MustLoad()

	// Init logger
	log := logger.SetupLogger(cfg.Log.Level)

	application := app.New(log, cfg)

	fmt.Println(application)

}
