package shoppinglists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ShoppingListsService) DeleteItem(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	itemID uuid.UUID,
	callerID uuid.UUID,
) error {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return err
	}

	if _, err := s.shoppingListsRepository.GetListByID(ctx, householdID, listID); err != nil {
		return fmt.Errorf("get shopping list from repository: %w", err)
	}

	if err := s.shoppingListsRepository.DeleteItem(ctx, listID, itemID); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}
