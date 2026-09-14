package core_domain

import (
	core_errors "cohesive-core/internal/core/errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeDeposit TransactionType = "deposit"
	TransactionTypeExpense TransactionType = "expense"
)

type HouseholdTransaction struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Type        TransactionType
	Amount      int64
	Description *string
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewHouseholdTransaction(
	id uuid.UUID,
	householdID uuid.UUID,
	version int,
	txType TransactionType,
	amount int64,
	description *string,
	createdBy uuid.UUID,
	createdAt, updatedAt time.Time,
) HouseholdTransaction {
	return HouseholdTransaction{
		ID:          id,
		HouseholdID: householdID,
		Version:     version,
		Type:        txType,
		Amount:      amount,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewHouseholdTransactionUninitialized(
	householdID uuid.UUID,
	txType TransactionType,
	amount int64,
	description *string,
	createdBy uuid.UUID,
) HouseholdTransaction {
	now := time.Now()

	return HouseholdTransaction{
		ID:          UninitializedID,
		HouseholdID: householdID,
		Version:     UninitializedVersion,
		Type:        txType,
		Amount:      amount,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *HouseholdTransaction) Validate() error {
	switch t.Type {
	case TransactionTypeDeposit, TransactionTypeExpense:
	default:
		return fmt.Errorf("invalid `Type` %q: %w", t.Type, core_errors.ErrInvalidArgument)
	}

	if t.Amount <= 0 {
		return fmt.Errorf("`Amount` must be positive: %w", core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descriptionLen := len([]rune(*t.Description))
		if descriptionLen > 300 {
			return fmt.Errorf("invalid `Description` len: %d: %w", descriptionLen, core_errors.ErrInvalidArgument)
		}
	}

	if t.UpdatedAt.Before(t.CreatedAt) {
		return fmt.Errorf("`UpdatedAt` is before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type HouseholdTransactionPatch struct {
	Amount      Nullable[int64]
	Description Nullable[string]
}

func NewHouseholdTransactionPatch(
	amount Nullable[int64],
	description Nullable[string],
) HouseholdTransactionPatch {
	return HouseholdTransactionPatch{
		Amount:      amount,
		Description: description,
	}
}

func (p *HouseholdTransactionPatch) Validate() error {
	if p.Amount.Set && p.Amount.Value == nil {
		return fmt.Errorf("`Amount` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (t *HouseholdTransaction) ApplyPatch(patch HouseholdTransactionPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate transaction patch: %w", err)
	}

	tmp := *t

	if patch.Amount.Set {
		tmp.Amount = *patch.Amount.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched transaction: %w", err)
	}

	*t = tmp

	return nil
}
