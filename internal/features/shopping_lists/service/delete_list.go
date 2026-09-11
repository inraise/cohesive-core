package shoppinglists_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *ShoppingListsService) DeleteList(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	callerID uuid.UUID,
) error {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return err
	}

	if err := s.shoppingListsRepository.DeleteList(ctx, householdID, listID); err != nil {
		return fmt.Errorf("delete shopping list: %w", err)
	}

	return nil
}
