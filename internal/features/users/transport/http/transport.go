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
	tokenManager *core_jwt.TokenManager
}

type UsersService interface {
	GetMe(
		ctx context.Context,
		id uuid.UUID,
	) (core_domain.User, error)
}

func (h *UsersHTTPHandler) NewUsersHTTPHandler(
	usersService UsersService,
	tokenManager *core_jwt.TokenManager,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
		tokenManager: tokenManager,
	}
}

func (h *UsersHTTPHandler) Routes() []core_transport_http_server.Route {
	authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)

	return []core_transport_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: h.GetMe,
			Middleware: []core_transport_http_middleware.Middleware{
				authenticate,
			},
		},
	}
}
