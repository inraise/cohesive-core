package tasks_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	IsHouseholdMember(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (bool, error)

	CreateTask(
		ctx context.Context,
		task core_domain.Task,
	) (core_domain.Task, error)

	DeleteTask(
		ctx context.Context,
		householdID uuid.UUID,
		taskID uuid.UUID,
	) error

	ListTasks(
		ctx context.Context,
		householdID uuid.UUID,
	) ([]core_domain.Task, error)
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}

func (s *TasksService) requireMember(ctx context.Context, householdID, userID uuid.UUID) error {
	isMember, err := s.tasksRepository.IsHouseholdMember(ctx, householdID, userID)
	if err != nil {
		return fmt.Errorf("check membership: %w", err)
	}

	if !isMember {
		return fmt.Errorf(
			"user %q is not a member of household %q: %w",
			userID, householdID, core_errors.ErrNotFound,
		)
	}

	return nil
}
