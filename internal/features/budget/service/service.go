package budget_service

import (
	"context"
	"fmt"

	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"

	"github.com/google/uuid"
)

type BudgetService struct {
	budgetRepository BudgetRepository
}

type BudgetRepository interface {
	IsHouseholdMember(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (bool, error)

	GetMemberRole(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (core_domain.HouseholdRole, error)

	CreateTransaction(
		ctx context.Context,
		tx core_domain.HouseholdTransaction,
	) (core_domain.HouseholdTransaction, error)

	GetTransactionByID(
		ctx context.Context,
		householdID uuid.UUID,
		txID uuid.UUID,
	) (core_domain.HouseholdTransaction, error)

	ListTransactions(
		ctx context.Context,
		householdID uuid.UUID,
	) ([]core_domain.HouseholdTransaction, error)

	PatchTransaction(
		ctx context.Context,
		tx core_domain.HouseholdTransaction,
	) (core_domain.HouseholdTransaction, error)

	DeleteTransaction(
		ctx context.Context,
		householdID uuid.UUID,
		txID uuid.UUID,
	) error

	GetTotals(
		ctx context.Context,
		householdID uuid.UUID,
	) (totalDeposits int64, totalExpenses int64, err error)

	GetMemberContributions(
		ctx context.Context,
		householdID uuid.UUID,
	) ([]core_domain.MemberContribution, error)
}

func NewBudgetService(
	budgetRepository BudgetRepository,
) *BudgetService {
	return &BudgetService{
		budgetRepository: budgetRepository,
	}
}

func (s *BudgetService) requireMember(ctx context.Context, householdID, userID uuid.UUID) error {
	isMember, err := s.budgetRepository.IsHouseholdMember(ctx, householdID, userID)
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

func (s *BudgetService) requireCanModify(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	tx core_domain.HouseholdTransaction,
) error {
	if tx.CreatedBy == callerID {
		return nil
	}

	role, err := s.budgetRepository.GetMemberRole(ctx, householdID, callerID)
	if err != nil {
		return fmt.Errorf("get caller role: %w", err)
	}

	if role != core_domain.HouseholdRoleOwner && role != core_domain.HouseholdRoleAdmin {
		return fmt.Errorf(
			"only the author or an owner/admin can modify this transaction: %w",
			core_errors.ErrForbidden,
		)
	}

	return nil
}
