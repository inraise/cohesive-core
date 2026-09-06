package tasks_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_pool "cohesive-core/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	task core_domain.Task,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, assigned_to = $4, updated_at = $5, version = version + 1
		WHERE id = $6 AND household_id = $7 AND version = $8
		RETURNING id, household_id, version, title, description, status, assigned_to, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		string(task.Status),
		task.AssignedTo,
		time.Now(),
		task.ID,
		task.HouseholdID,
		task.Version,
	)

	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, core_pool.ErrNoRows) {
			return core_domain.Task{}, fmt.Errorf(
				"task with id %q concurrently accessed: %w", task.ID, core_errors.ErrConflict,
			)
		}

		return core_domain.Task{}, err
	}

	return updated, nil
}
