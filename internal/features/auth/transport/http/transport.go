package auth_transport_http

import (
	"context"
	"net/http"
	"time"

	core_domain "cohesive-core/internal/core/domain"
	core_pool_redis "cohesive-core/internal/core/repository/redis/pool"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	auth_service "cohesive-core/internal/features/auth/service"
)

const (
	loginRateLimit  = 5
	loginRateWindow = time.Minute

	refreshRateLimit  = 20
	refreshRateWindow = time.Minute

	verifyEmailRateLimit  = 10
	verifyEmailRateWindow = time.Minute

	resendVerificationRateLimit  = 3
	resendVerificationRateWindow = 10 * time.Minute
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

	VerifyEmail(
		ctx context.Context,
		token string,
	) (core_domain.User, error)

	ResendVerification(
		ctx context.Context,
		email string,
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
		{
			Method:  http.MethodGet,
			Path:    "/auth/verify-email",
			Handler: h.VerifyEmailPage,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/verify-email/confirm",
			Handler: h.ConfirmVerifyEmail,
			Middleware: []core_transport_http_middleware.Middleware{
				core_transport_http_middleware.RateLimit(
					h.redisClient, "verify-email", verifyEmailRateLimit, verifyEmailRateWindow,
				),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/verify-email",
			Handler: h.VerifyEmail,
			Middleware: []core_transport_http_middleware.Middleware{
				core_transport_http_middleware.RateLimit(
					h.redisClient, "verify-email", verifyEmailRateLimit, verifyEmailRateWindow,
				),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/resend-verification",
			Handler: h.ResendVerification,
			Middleware: []core_transport_http_middleware.Middleware{
				core_transport_http_middleware.RateLimit(
					h.redisClient, "resend-verification", resendVerificationRateLimit, resendVerificationRateWindow,
				),
			},
		},
	}
}
