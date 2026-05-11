package main

import (
	"log/slog"
	"os"

	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/app"
	"github.com/Romasmi/s-shop-microservices/delivery-service/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	application := app.NewApp(cfg)
	if err := application.Run(); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}
