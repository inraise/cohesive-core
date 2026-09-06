package tasks_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	householdID uuid.UUID,
	taskID uuid.UUID,
	callerID uuid.UUID,
) error {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return err
	}

	if err := s.tasksRepository.DeleteTask(ctx, householdID, taskID); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
