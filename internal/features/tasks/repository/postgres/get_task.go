package tasks_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *TasksRepository) GetTaskByID(
	ctx context.Context,
	householdID uuid.UUID,
	taskID uuid.UUID,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, household_id, version, title, description, status, assigned_to, created_by, created_at, updated_at
		FROM tasks
		WHERE id = $1 AND household_id = $2;
	`

	row := r.pool.QueryRow(ctx, query, taskID, householdID)

	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.Task{}, fmt.Errorf(
				"task with id %q in household %q: %w", taskID, householdID, core_errors.ErrNotFound,
			)
		}

		return core_domain.Task{}, err
	}

	return task, nil
}
