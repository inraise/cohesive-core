package shoppinglists_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	shoppinglists_service "cohesive-core/internal/features/shopping_lists/service"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ShoppingListsHTTPHandler struct {
	shoppingListsService ShoppingListsService
	tokenManager         *core_jwt.TokenManager
}

type ShoppingListsService interface {
	CreateList(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		request shoppinglists_service.CreateListRequest,
	) (core_domain.ShoppingList, error)

	ListLists(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) ([]core_domain.ShoppingList, error)

	GetList(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		callerID uuid.UUID,
	) (shoppinglists_service.ListWithItems, error)

	RenameList(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		callerID uuid.UUID,
		request shoppinglists_service.RenameListRequest,
	) (core_domain.ShoppingList, error)

	DeleteList(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		callerID uuid.UUID,
	) error

	CreateItem(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		callerID uuid.UUID,
		request shoppinglists_service.CreateItemRequest,
	) (core_domain.ShoppingListItem, error)

	PatchItem(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		itemID uuid.UUID,
		callerID uuid.UUID,
		patch core_domain.ShoppingListItemPatch,
	) (core_domain.ShoppingListItem, error)

	DeleteItem(
		ctx context.Context,
		householdID uuid.UUID,
		listID uuid.UUID,
		itemID uuid.UUID,
		callerID uuid.UUID,
	) error
}

func NewShoppingListsHTTPHandler(
	shoppingListsService ShoppingListsService,
	tokenManager *core_jwt.TokenManager,
) *ShoppingListsHTTPHandler {
	return &ShoppingListsHTTPHandler{
		shoppingListsService: shoppingListsService,
		tokenManager:         tokenManager,
	}
}

func (h *ShoppingListsHTTPHandler) Routes() []core_transport_http_server.Route {
	authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)
	withAuth := []core_transport_http_middleware.Middleware{authenticate}

	return []core_transport_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/households/{id}/shopping-lists",
			Handler:    h.CreateList,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/shopping-lists",
			Handler:    h.ListLists,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/shopping-lists/{list_id}",
			Handler:    h.GetList,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/households/{id}/shopping-lists/{list_id}",
			Handler:    h.RenameList,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/households/{id}/shopping-lists/{list_id}",
			Handler:    h.DeleteList,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/households/{id}/shopping-lists/{list_id}/items",
			Handler:    h.CreateItem,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/households/{id}/shopping-lists/{list_id}/items/{item_id}",
			Handler:    h.PatchItem,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/households/{id}/shopping-lists/{list_id}/items/{item_id}",
			Handler:    h.DeleteItem,
			Middleware: withAuth,
		},
	}
}
