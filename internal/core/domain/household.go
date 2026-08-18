package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Household struct {
	ID      uuid.UUID
	Version int

	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewHousehold(
	id uuid.UUID,
	version int,
	name string,
	createdAt, updatedAt time.Time,
) Household {
	return Household{
		ID:        id,
		Version:   version,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewHouseholdUninitialized(name string) Household {
	return Household{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (h *Household) Validate() error {
	nameLen := len([]rune(h.Name))
	if nameLen < 1 || nameLen > 100 {
		return fmt.Errorf(
			"invalid `Name` len: %d: %w",
			nameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if h.UpdatedAt.Before(h.CreatedAt) {
		return fmt.Errorf(
			"`UpdatedAt` is before `CreatedAt`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
