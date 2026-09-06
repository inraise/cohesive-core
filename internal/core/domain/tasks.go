package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Title       string
	Description *string
	Status      TaskStatus
	AssignedTo  *uuid.UUID
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTask(
	id uuid.UUID,
	householdID uuid.UUID,
	version int,
	title string,
	description *string,
	status TaskStatus,
	assignedTo *uuid.UUID,
	createdBy uuid.UUID,
	createdAt, updatedAt time.Time,
) Task {
	return Task{
		ID:          id,
		HouseholdID: householdID,
		Version:     version,
		Title:       title,
		Description: description,
		Status:      status,
		AssignedTo:  assignedTo,
		CreatedBy:   createdBy,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewTaskUninitialized(
	householdID uuid.UUID,
	title string,
	description *string,
	assignedTo *uuid.UUID,
	createdBy uuid.UUID,
) Task {
	now := time.Now()

	return Task{
		ID:          UninitializedID,
		HouseholdID: householdID,
		Version:     UninitializedVersion,
		Title:       title,
		Description: description,
		Status:      TaskStatusPending,
		AssignedTo:  assignedTo,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Task) Validate() error {
	titleLen := len([]rune(t.Title))
	if titleLen < 1 || titleLen > 200 {
		return fmt.Errorf("invalid `Title` len: %d: %w", titleLen, core_errors.ErrInvalidArgument)
	}

	switch t.Status {
	case TaskStatusPending, TaskStatusInProgress, TaskStatusDone:
	default:
		return fmt.Errorf("invalid `Status` %q: %w", t.Status, core_errors.ErrInvalidArgument)
	}

	if t.UpdatedAt.Before(t.CreatedAt) {
		return fmt.Errorf("`UpdatedAt` is before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Status      Nullable[string]
	AssignedTo  Nullable[uuid.UUID]
}

func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	status Nullable[string],
	assignedTo Nullable[uuid.UUID],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Status:      status,
		AssignedTo:  assignedTo,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Status.Set && p.Status.Value == nil {
		return fmt.Errorf("`Status` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Status.Set {
		tmp.Status = TaskStatus(*patch.Status.Value)
	}

	if patch.AssignedTo.Set {
		tmp.AssignedTo = patch.AssignedTo.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp

	return nil
}
