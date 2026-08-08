package auth_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	auth_service "cohesive-core/internal/features/auth/service"
	"context"
	"net/http"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	CreateUser(
		ctx context.Context,
		user auth_service.CreateUserRequest,
	) (core_domain.User, error)

	LoginUser(
		ctx context.Context,
		email, password string,
	) (*auth_service.LoginDTOResponse, error)

	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (*auth_service.LoginDTOResponse, error)

	LogoutUser(
		ctx context.Context,
		refreshToken string,
	) error
}

func NewAuthHTTPHandler(
	authService AuthService,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
	}
}

func (h *AuthHTTPHandler) Routes() []core_transport_http_server.Route {
	return []core_transport_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.LogoutUser,
		},
	}
}
