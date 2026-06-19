package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDbPool(ctx context.Context) (*pgxpool.Pool, error) {
	connStr := getConnStr()
	if err := runMigrations(connStr); err != nil {
		return nil, fmt.Errorf("Ошибка миграции: %w", err)
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга конфига: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 2 * time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	defer dbPool.Close()
	pgxCtx, cancelPgx := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPgx()

	err = dbPool.Ping(pgxCtx)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	fmt.Println("Успешное соединение с БД")
	return dbPool, nil
}
