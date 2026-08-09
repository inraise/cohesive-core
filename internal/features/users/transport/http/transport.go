package users_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	GetMe(
		ctx context.Context,
		id uuid.UUID,
	) (core_domain.User, error)
}

func (h *UsersHTTPHandler) NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes(tokenManager *core_jwt.TokenManager) []core_transport_http_server.Route {
	var middlewares []core_transport_http_middleware.Middleware
	middlewares = append(middlewares, core_transport_http_middleware.Authenticate(tokenManager))

	return []core_transport_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/users/me",
			Handler:    h.GetMe,
			Middleware: middlewares,
		},
	}
}
