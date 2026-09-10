package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ListWithItems struct {
	List  core_domain.ShoppingList
	Items []core_domain.ShoppingListItem
}

func (s *ShoppingListsService) GetList(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	callerID uuid.UUID,
) (ListWithItems, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return ListWithItems{}, err
	}

	list, err := s.shoppingListsRepository.GetListByID(ctx, householdID, listID)
	if err != nil {
		return ListWithItems{}, fmt.Errorf("get shopping list from repository: %w", err)
	}

	items, err := s.shoppingListsRepository.ListItems(ctx, listID)
	if err != nil {
		return ListWithItems{}, fmt.Errorf("list items from repository: %w", err)
	}

	return ListWithItems{List: list, Items: items}, nil
}
