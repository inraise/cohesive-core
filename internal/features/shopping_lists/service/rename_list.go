package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type RenameListRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (s *ShoppingListsService) RenameList(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	callerID uuid.UUID,
	request RenameListRequest,
) (core_domain.ShoppingList, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.ShoppingList{}, err
	}

	current, err := s.shoppingListsRepository.GetListByID(ctx, householdID, listID)
	if err != nil {
		return core_domain.ShoppingList{}, fmt.Errorf("get shopping list from repository: %w", err)
	}

	updated, err := s.shoppingListsRepository.RenameList(ctx, listID, request.Name, current.Version)
	if err != nil {
		return core_domain.ShoppingList{}, fmt.Errorf("rename shopping list: %w", err)
	}

	return updated, nil
}
