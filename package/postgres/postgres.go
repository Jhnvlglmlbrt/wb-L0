package postgres

import (
	"context"
	"fmt"

	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.PG) (*pgxpool.Pool, error) {
	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return pool, nil
}
