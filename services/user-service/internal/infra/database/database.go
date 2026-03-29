package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/Romasmi/s-shop-microservices/internal/utils/time_utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	DB     *pgxpool.Pool
	Config *config.Database
}

func (c *Connection) Connect() error {
	pgConfig, err := pgxpool.ParseConfig(c.Config.URL)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}
	pgConfig.MaxConns = int32(c.Config.MaxConnections)
	pgConfig.MinConns = int32(c.Config.MinConnections)
	pgConfig.MaxConnLifetime = time_utils.MinutesToNanoseconds(c.Config.MaxConnectionLifetime)
	pgConfig.MaxConnIdleTime = time_utils.MinutesToNanoseconds(c.Config.MaxConnectionIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c.DB, err = pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := c.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

func (c *Connection) Close() {
	if c.DB != nil {
		c.DB.Close()
		slog.Info("database connection closed")
	}
}

func (c *Connection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return c.DB.Ping(ctx)
}
