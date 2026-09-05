package tasks_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) ListTasks(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) ([]core_domain.Task, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return nil, err
	}

	tasks, err := s.tasksRepository.ListTasks(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("list tasks from repository: %w", err)
	}

	return tasks, nil
}
