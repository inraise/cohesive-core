package shoppinglists_transport_http

import (
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
)

type ShoppingListsHTTPHandler struct {
	shoppingListsService ShoppingListsService
	tokenManager         *core_jwt.TokenManager
}

type ShoppingListsService interface {
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
	_ = withAuth

	return []core_transport_http_server.Route{}
}
