package shoppinglists_service

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ShoppingListsService struct {
	shoppingListsRepository ShoppingListsRepository
}

type ShoppingListsRepository interface {
	IsHouseholdMember(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (bool, error)

	CreateList(
		ctx context.Context,
		list core_domain.ShoppingList,
	) (core_domain.ShoppingList, error)

	GetListByID(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
	) (core_domain.ShoppingList, error)

	ListLists(
		ctx context.Context,
		householdID uuid.UUID,
	) ([]core_domain.ShoppingList, error)

	RenameList(
		ctx context.Context,
		listID uuid.UUID,
		name string,
		version int,
	) (core_domain.ShoppingList, error)

	DeleteList(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
	) error

	ListItems(
		ctx context.Context,
		shoppingListID uuid.UUID,
	) ([]core_domain.ShoppingListItem, error)

	CreateItem(
		ctx context.Context,
		item core_domain.ShoppingListItem,
	) (core_domain.ShoppingListItem, error)

	GetItemByID(
		ctx context.Context,
		shoppingListID uuid.UUID,
		itemID uuid.UUID,
	) (core_domain.ShoppingListItem, error)

	PatchItem(
		ctx context.Context,
		item core_domain.ShoppingListItem,
	) (core_domain.ShoppingListItem, error)

	DeleteItem(
		ctx context.Context,
		shoppingListID uuid.UUID,
		itemID uuid.UUID,
	) error
}

func NewShoppingListsService(
	shoppingListsRepository ShoppingListsRepository,
) *ShoppingListsService {
	return &ShoppingListsService{
		shoppingListsRepository: shoppingListsRepository,
	}
}

func (s *ShoppingListsService) requireMember(ctx context.Context, householdID, userID uuid.UUID) error {
	isMember, err := s.shoppingListsRepository.IsHouseholdMember(ctx, householdID, userID)
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
