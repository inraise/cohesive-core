package households_repository_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *HouseholdsRepository) DeleteHousehold(
	ctx context.Context,
	householdID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM households WHERE id = $1;`

	if _, err := r.pool.Exec(ctx, query, householdID); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
