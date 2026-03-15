package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/Romasmi/s-shop-microservices/internal/infra/database"
)

type App struct {
	DbConn *database.Connection
	Config *config.Config
	server *http.Server
}

func (a *App) GetDB() *database.Connection {
	return a.DbConn
}

func NewApp(configPath string) (*App, error) {
	app := &App{}
	return app, app.init(configPath)
}

func (a *App) init(configPath string) error {
	envConfig, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading Config: %v\n", err)
	}
	a.Config = envConfig

	dbConn := &database.Connection{Config: &envConfig.Database}
	if err = dbConn.Connect(); err != nil {
		return fmt.Errorf("error connecting to DB: %v\n", err)
	}
	a.DbConn = dbConn

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	var shutdownErr error

	if a.server != nil {
		fmt.Println("Shutting down HTTP server...")
		if err := a.server.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("server shutdown error: %w", err)
			fmt.Printf("HTTP server shutdown error: %v\n", err)
		}
	}

	if a.DbConn != nil && a.DbConn.DB != nil {
		fmt.Println("Closing database connections...")
		select {
		case <-ctx.Done():
			fmt.Println("Shutdown timeout reached, forcing database close")
		default:
			a.DbConn.DB.Close()
		}
	}

	fmt.Println("Cleanup completed")
	return shutdownErr
}
