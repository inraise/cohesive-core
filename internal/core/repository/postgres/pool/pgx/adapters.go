package core_pool_pgx

import (
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {

	err := r.Row.Scan(dest...)
	if err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		foreignKeyViolationCode = "23503"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == foreignKeyViolationCode {
			return fmt.Errorf("%v: %w", err, core_pool.ErrViolatesForeignKey)
		}
	}

	return fmt.Errorf("%v: %w", err, core_pool.ErrUnknown)
}
