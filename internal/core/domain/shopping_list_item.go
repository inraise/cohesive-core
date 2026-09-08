package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ShoppingListItem struct {
	ID             uuid.UUID
	ShoppingListID uuid.UUID
	Version        int

	Name        string
	Quantity    *string
	IsPurchased bool
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewShoppingListItem(
	id uuid.UUID,
	shoppingListID uuid.UUID,
	version int,
	name string,
	quantity *string,
	isPurchased bool,
	createdBy uuid.UUID,
	createdAt, updatedAt time.Time,
) ShoppingListItem {
	return ShoppingListItem{
		ID:             id,
		ShoppingListID: shoppingListID,
		Version:        version,
		Name:           name,
		Quantity:       quantity,
		IsPurchased:    isPurchased,
		CreatedBy:      createdBy,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

func NewShoppingListItemUninitialized(
	shoppingListID uuid.UUID,
	name string,
	quantity *string,
	createdBy uuid.UUID,
) ShoppingListItem {
	now := time.Now()

	return ShoppingListItem{
		ID:             UninitializedID,
		ShoppingListID: shoppingListID,
		Version:        UninitializedVersion,
		Name:           name,
		Quantity:       quantity,
		IsPurchased:    false,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (i *ShoppingListItem) Validate() error {
	nameLen := len([]rune(i.Name))
	if nameLen < 1 || nameLen > 200 {
		return fmt.Errorf("invalid `Name` len: %d: %w", nameLen, core_errors.ErrInvalidArgument)
	}

	if i.Quantity != nil {
		quantityLen := len([]rune(*i.Quantity))
		if quantityLen > 50 {
			return fmt.Errorf("invalid `Quantity` len: %d: %w", quantityLen, core_errors.ErrInvalidArgument)
		}
	}

	if i.UpdatedAt.Before(i.CreatedAt) {
		return fmt.Errorf("`UpdatedAt` is before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type ShoppingListItemPatch struct {
	Name        Nullable[string]
	Quantity    Nullable[string]
	IsPurchased Nullable[bool]
}

func NewShoppingListItemPatch(
	name Nullable[string],
	quantity Nullable[string],
	isPurchased Nullable[bool],
) ShoppingListItemPatch {
	return ShoppingListItemPatch{
		Name:        name,
		Quantity:    quantity,
		IsPurchased: isPurchased,
	}
}

func (p *ShoppingListItemPatch) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf("`Name` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.IsPurchased.Set && p.IsPurchased.Value == nil {
		return fmt.Errorf("`IsPurchased` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (i *ShoppingListItem) ApplyPatch(patch ShoppingListItemPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate shopping list item patch: %w", err)
	}

	tmp := *i

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if patch.Quantity.Set {
		tmp.Quantity = patch.Quantity.Value
	}

	if patch.IsPurchased.Set {
		tmp.IsPurchased = *patch.IsPurchased.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched shopping list item: %w", err)
	}

	*i = tmp

	return nil
}
