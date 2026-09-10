package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ShoppingListsService) ListLists(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
) ([]core_domain.ShoppingList, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return nil, err
	}

	lists, err := s.shoppingListsRepository.ListLists(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("list shopping lists from repository: %w", err)
	}

	return lists, nil
}
