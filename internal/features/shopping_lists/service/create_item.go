package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateItemRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=200"`
	Quantity *string `json:"quantity" validate:"omitempty,max=50"`
}

func (s *ShoppingListsService) CreateItem(
	ctx context.Context,
	householdID uuid.UUID,
	listID uuid.UUID,
	callerID uuid.UUID,
	request CreateItemRequest,
) (core_domain.ShoppingListItem, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.ShoppingListItem{}, err
	}

	if _, err := s.shoppingListsRepository.GetListByID(ctx, householdID, listID); err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("get shopping list from repository: %w", err)
	}

	itemDomain := core_domain.NewShoppingListItemUninitialized(listID, request.Name, request.Quantity, callerID)

	if err := itemDomain.Validate(); err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("validate shopping list item domain: %w", err)
	}

	created, err := s.shoppingListsRepository.CreateItem(ctx, itemDomain)
	if err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("create item: %w", err)
	}

	return created, nil
}
