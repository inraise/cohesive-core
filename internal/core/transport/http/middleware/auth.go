package core_transport_http_middleware

import (
	core_errors "cohesive-core/internal/core/errors"
	core_jwt "cohesive-core/internal/core/jwt"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type userIDContextKey struct{}

var userIDKey = userIDContextKey{}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)

	return id, ok
}

func Authenticate(tokenManager *core_jwt.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, w)

			token, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				responseHandler.ErrorResponse(
					fmt.Errorf("missing or malformed Authorization header: %w", core_errors.ErrUnauthorized),
					"unauthorized",
				)

				return
			}

			claims, err := tokenManager.ValidateToken(token)
			if err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("%v: %w", err, core_errors.ErrUnauthorized),
					"unauthorized",
				)

				return
			}

			ctx = context.WithValue(ctx, userIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimPrefix(header, prefix)
	if token == "" {
		return "", false
	}

	return token, true
}
