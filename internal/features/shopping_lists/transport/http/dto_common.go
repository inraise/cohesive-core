package shoppinglists_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type ShoppingListDTOResponse struct {
	ID          uuid.UUID `json:"id"`
	HouseholdID uuid.UUID `json:"household_id"`
	Version     int       `json:"version"`

	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func shoppingListDTOFromDomain(list core_domain.ShoppingList) ShoppingListDTOResponse {
	return ShoppingListDTOResponse{
		ID:          list.ID,
		HouseholdID: list.HouseholdID,
		Version:     list.Version,
		Name:        list.Name,
		CreatedBy:   list.CreatedBy,
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
	}
}

type ShoppingListItemDTOResponse struct {
	ID             uuid.UUID `json:"id"`
	ShoppingListID uuid.UUID `json:"shopping_list_id"`
	Version        int       `json:"version"`

	Name        string    `json:"name"`
	Quantity    *string   `json:"quantity"`
	IsPurchased bool      `json:"is_purchased"`
	CreatedBy   uuid.UUID `json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func itemDTOFromDomain(item core_domain.ShoppingListItem) ShoppingListItemDTOResponse {
	return ShoppingListItemDTOResponse{
		ID:             item.ID,
		ShoppingListID: item.ShoppingListID,
		Version:        item.Version,
		Name:           item.Name,
		Quantity:       item.Quantity,
		IsPurchased:    item.IsPurchased,
		CreatedBy:      item.CreatedBy,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

type ShoppingListWithItemsDTOResponse struct {
	ShoppingListDTOResponse
	Items []ShoppingListItemDTOResponse `json:"items"`
}
