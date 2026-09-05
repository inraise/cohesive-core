package auth_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_pool_redis "cohesive-core/internal/core/repository/redis/pool"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	auth_service "cohesive-core/internal/features/auth/service"
	"context"
	"net/http"
	"time"
)

const (
	loginRateLimit  = 5
	loginRateWindow = time.Minute

	refreshRateLimit  = 20
	refreshRateWindow = time.Minute
)

type AuthHTTPHandler struct {
	authService AuthService
	redisClient core_pool_redis.Client
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
	redisClient core_pool_redis.Client,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService: authService,
		redisClient: redisClient,
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
			Middleware: []core_transport_http_middleware.Middleware{
				core_transport_http_middleware.RateLimit(h.redisClient, "login", loginRateLimit, loginRateWindow),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
			Middleware: []core_transport_http_middleware.Middleware{
				core_transport_http_middleware.RateLimit(h.redisClient, "refresh", refreshRateLimit, refreshRateWindow),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/logout",
			Handler: h.LogoutUser,
		},
	}
}
