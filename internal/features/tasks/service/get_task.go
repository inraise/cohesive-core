package tasks_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	householdID uuid.UUID,
	taskID uuid.UUID,
	callerID uuid.UUID,
) (core_domain.Task, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.Task{}, err
	}

	task, err := s.tasksRepository.GetTaskByID(ctx, householdID, taskID)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	return task, nil
}
