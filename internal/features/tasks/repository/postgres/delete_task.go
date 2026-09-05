package tasks_repository_postgres

import (
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *TasksRepository) DeleteTask(
	ctx context.Context,
	householdID uuid.UUID,
	taskID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM tasks WHERE id = $1 AND household_id = $2;`

	tag, err := r.pool.Exec(ctx, query, taskID, householdID)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("task with id %q: %w", taskID, core_errors.ErrNotFound)
	}

	return nil
}
