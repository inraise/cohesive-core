package shoppinglists_repository_postgres

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingListModel struct {
	ID          uuid.UUID
	HouseholdID uuid.UUID
	Version     int

	Name      string
	CreatedBy uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ShoppingListItemModel struct {
	ID             uuid.UUID
	ShoppingListID uuid.UUID
	Version        int

	Name        string
	Quantity    *string
	IsPurchased bool
	CreatedBy   uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
}
