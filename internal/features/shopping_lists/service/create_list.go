package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreateListRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (s *ShoppingListsService) CreateList(
	ctx context.Context,
	householdID uuid.UUID,
	callerID uuid.UUID,
	request CreateListRequest,
) (core_domain.ShoppingList, error) {
	if err := s.requireMember(ctx, householdID, callerID); err != nil {
		return core_domain.ShoppingList{}, err
	}

	listDomain := core_domain.NewShoppingListUninitialized(householdID, request.Name, callerID)

	if err := listDomain.Validate(); err != nil {
		return core_domain.ShoppingList{}, fmt.Errorf("validate shopping list domain: %w", err)
	}

	created, err := s.shoppingListsRepository.CreateList(ctx, listDomain)
	if err != nil {
		return core_domain.ShoppingList{}, fmt.Errorf("create shopping list: %w", err)
	}

	return created, nil
}
