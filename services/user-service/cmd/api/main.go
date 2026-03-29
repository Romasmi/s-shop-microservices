package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Romasmi/s-shop-microservices/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	basePath := os.Getenv("APP_BASE_PATH")
	if basePath == "" {
		basePath = "." // current directory
	}

	appInstance, err := app.NewApp(basePath)
	if err != nil {
		fmt.Printf("error while app init: %v", err)
		os.Exit(1)
	}
	api := app.NewApi(appInstance)

	migrationsPath := filepath.Join(basePath, "migrations")
	absPath, _ := filepath.Abs(migrationsPath)
	m, err := migrate.New(
		"file://"+absPath,
		appInstance.Config.Database.URL,
	)
	if m == nil || err != nil {
		fmt.Printf("unable to create migrations driver: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
			fmt.Printf("Error closing migration driver - source: %v, db: %v\n", sourceErr, dbErr)
		}
	}()

	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Printf("error while running up migrations: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := api.Run(); err != nil {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	}

	fmt.Println("Server stopped")
}
