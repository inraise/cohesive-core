package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ShoppingListsService) PatchItem(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	itemID uuid.UUID,
	callerID uuid.UUID,
	patch core_domain.ShoppingListItemPatch,
) (core_domain.ShoppingListItem, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.ShoppingListItem{}, err
	}

	if _, err := s.shoppingListsRepository.GetListByID(ctx, householdID, listID); err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("get shopping list from repository: %w", err)
	}

	item, err := s.shoppingListsRepository.GetItemByID(ctx, listID, itemID)
	if err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("get item from repository: %w", err)
	}

	if err := item.ApplyPatch(patch); err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("apply item patch: %w", err)
	}

	updated, err := s.shoppingListsRepository.PatchItem(ctx, item)
	if err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("patch item: %w", err)
	}

	return updated, nil
}
