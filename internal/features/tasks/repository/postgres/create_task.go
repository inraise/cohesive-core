package tasks_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task core_domain.Task,
) (core_domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO tasks (household_id, title, description, status, assigned_to, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, household_id, version, title, description, status, assigned_to, created_by, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		task.HouseholdID,
		task.Title,
		task.Description,
		string(task.Status),
		task.AssignedTo,
		task.CreatedBy,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return scanTask(row)
}
