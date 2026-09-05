package tasks_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateTaskRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=200"`
	Description *string    `json:"description" validate:"omitempty,max=2000"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
}

func (s *TasksService) CreateTask(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	request CreateTaskRequest,
) (core_domain.Task, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.Task{}, err
	}

	if request.AssignedTo != nil {
		isMember, err := s.tasksRepository.IsHouseholdMember(ctx, householdID, *request.AssignedTo)
		if err != nil {
			return core_domain.Task{}, fmt.Errorf("check assignee membership: %w", err)
		}

		if !isMember {
			return core_domain.Task{}, fmt.Errorf(
				"assignee %q is not a member of household %q: %w",
				*request.AssignedTo, householdID, core_errors.ErrInvalidArgument,
			)
		}
	}

	taskDomain := core_domain.NewTaskUninitialized(
		householdID, request.Title, request.Description, request.AssignedTo, callerID,
	)

	if err := taskDomain.Validate(); err != nil {
		return core_domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}

	created, err := s.tasksRepository.CreateTask(ctx, taskDomain)
	if err != nil {
		return core_domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return created, nil
}
