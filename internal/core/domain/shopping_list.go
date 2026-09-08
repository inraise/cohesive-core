package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ShoppingList struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Name      string
	CreatedBy uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewShoppingList(
	id uuid.UUID,
	householdID uuid.UUID,
	version int,
	name string,
	createdBy uuid.UUID,
	createdAt, updatedAt time.Time,
) ShoppingList {
	return ShoppingList{
		ID:          id,
		HouseholdID: householdID,
		Version:     version,
		Name:        name,
		CreatedBy:   createdBy,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewShoppingListUninitialized(
	householdID uuid.UUID,
	name string,
	createdBy uuid.UUID,
) ShoppingList {
	now := time.Now()

	return ShoppingList{
		ID:          UninitializedID,
		HouseholdID: householdID,
		Version:     UninitializedVersion,
		Name:        name,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (l *ShoppingList) Validate() error {
	nameLen := len([]rune(l.Name))
	if nameLen < 1 || nameLen > 100 {
		return fmt.Errorf("invalid `Name` len: %d: %w", nameLen, core_errors.ErrInvalidArgument)
	}

	if l.UpdatedAt.Before(l.CreatedAt) {
		return fmt.Errorf("`UpdatedAt` is before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
