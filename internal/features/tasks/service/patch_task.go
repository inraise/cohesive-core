package tasks_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	householdID uuid.UUID,
	taskID uuid.UUID,
	callerID uuid.UUID,
	patch core_domain.TaskPatch,
) (core_domain.Task, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.Task{}, err
	}

	if patch.AssignedTo.Set && patch.AssignedTo.Value != nil {
		isMember, err := s.tasksRepository.IsHouseholdMember(ctx, householdID, *patch.AssignedTo.Value)
		if err != nil {
			return core_domain.Task{}, fmt.Errorf("check assignee membership: %w", err)
		}

		if !isMember {
			return core_domain.Task{}, fmt.Errorf(
				"assignee %q is not a member of household %q: %w",
				*patch.AssignedTo.Value, householdID, core_errors.ErrInvalidArgument,
			)
		}
	}

	task, err := s.tasksRepository.GetTaskByID(ctx, householdID, taskID)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("get task from repository: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return core_domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	updated, err := s.tasksRepository.PatchTask(ctx, task)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return updated, nil
}
