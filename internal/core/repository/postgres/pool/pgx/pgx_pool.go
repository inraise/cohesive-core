package core_pool_pgx

import (
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(ctx context.Context, config Config) (*Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgx pool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxRows{rows}, nil
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_pool.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)

	return pgxRow{row}
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_pool.CommandTag, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}

	return pgxCommandTag{tag}, nil
}
