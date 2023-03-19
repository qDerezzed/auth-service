package postgres

import (
	"context"
	"fmt"

	"auth-service/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func New(cfg config.Config) (*pgxpool.Pool, error) {
	databaseURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	dbpool, err := pgxpool.Connect(context.Background(), databaseURL)
	return dbpool, err
}
