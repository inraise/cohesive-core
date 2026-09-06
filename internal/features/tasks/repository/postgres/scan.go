package tasks_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"fmt"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(row scanner) (core_domain.Task, error) {
	var (
		model  TaskModel
		status string
	)

	err := row.Scan(
		&model.ID,
		&model.HouseholdID,
		&model.Version,
		&model.Title,
		&model.Description,
		&status,
		&model.AssignedTo,
		&model.CreatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewTask(
		model.ID,
		model.HouseholdID,
		model.Version,
		model.Title,
		model.Description,
		core_domain.TaskStatus(status),
		model.AssignedTo,
		model.CreatedBy,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
