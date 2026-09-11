package shoppinglists_repository_postgres

import (
	core_domain "cohesive-core/internal/core/domain"
	"fmt"
)

type scanner interface {
	Scan(dest ...any) error
}

func scanShoppingList(row scanner) (core_domain.ShoppingList, error) {
	var model ShoppingListModel

	err := row.Scan(
		&model.ID,
		&model.HouseholdID,
		&model.Version,
		&model.Name,
		&model.CreatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return core_domain.ShoppingList{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewShoppingList(
		model.ID,
		model.HouseholdID,
		model.Version,
		model.Name,
		model.CreatedBy,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}

func scanShoppingListItem(row scanner) (core_domain.ShoppingListItem, error) {
	var model ShoppingListItemModel

	err := row.Scan(
		&model.ID,
		&model.ShoppingListID,
		&model.Version,
		&model.Name,
		&model.Quantity,
		&model.IsPurchased,
		&model.CreatedBy,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return core_domain.ShoppingListItem{}, fmt.Errorf("scan error: %w", err)
	}

	return core_domain.NewShoppingListItem(
		model.ID,
		model.ShoppingListID,
		model.Version,
		model.Name,
		model.Quantity,
		model.IsPurchased,
		model.CreatedBy,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
