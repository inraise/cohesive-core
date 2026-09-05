package core_transport_http_middleware

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_pool_redis "cohesive-core/internal/core/repository/redis/pool"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RateLimit(
	redisClient core_pool_redis.Client,
	keyPrefix string,
	limit int64,
	window time.Duration,
) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, w)

			key := fmt.Sprintf("ratelimit:%s:%s", keyPrefix, clientIP(r))

			count, err := redisClient.Incr(ctx, key)
			if err != nil {
				log.Warn("rate limit check failed, allowing request through", zap.Error(err))
				next.ServeHTTP(w, r)

				return
			}

			if count == 1 {
				if _, err := redisClient.Expire(ctx, key, window); err != nil {
					log.Warn("failed to set rate limit window", zap.Error(err))
				}
			}

			if count > limit {
				responseHandler.ErrorResponse(
					fmt.Errorf("rate limit exceeded for %q: %w", keyPrefix, core_errors.ErrTooManyRequests),
					"too many requests, try again later",
				)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
