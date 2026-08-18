package main

import (
	core_config "cohesive-core/internal/core/config"
	core_jwt "cohesive-core/internal/core/jwt"
	core_logger "cohesive-core/internal/core/logger"
	core_pool_pgx "cohesive-core/internal/core/repository/postgres/pool/pgx"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	auth_repository_postgres "cohesive-core/internal/features/auth/repository/postgres"
	auth_service "cohesive-core/internal/features/auth/service"
	auth_transport_http "cohesive-core/internal/features/auth/transport/http"
	households_repository_postgres "cohesive-core/internal/features/households/repository/postgres"
	households_service "cohesive-core/internal/features/households/service"
	households_transport_http "cohesive-core/internal/features/households/transport/http"
	users_repository_postgres "cohesive-core/internal/features/users/repository/postgres"
	users_service "cohesive-core/internal/features/users/service"
	users_transport_http "cohesive-core/internal/features/users/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("app time zone", zap.Any("time_zone", time.Local))
	logger.Debug("initializing postgres connection pool")
	pool, err := core_pool_pgx.NewPool(
		ctx,
		core_pool_pgx.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing JWT token manager")
	jwtConfig := core_jwt.NewConfigMust()
	tokenManager, err := core_jwt.NewTokenManager(jwtConfig.SecretKey, jwtConfig.AccessTTL)
	if err != nil {
		logger.Fatal("failed to init JWT token manager", zap.Error(err))
	}

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	authRepository := auth_repository_postgres.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository, tokenManager, jwtConfig.RefreshTTL)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_repository_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, tokenManager)

	logger.Debug("initializing feature", zap.String("feature", "households"))
	householdsRepository := households_repository_postgres.NewHouseholdsRepository(pool)
	householdsService := households_service.NewHouseholdsService(householdsRepository)
	householdsTransportHTTP := households_transport_http.NewHouseholdsHTTPHandler(householdsService, tokenManager)

	logger.Debug("initializing HTTP server")

	httpConfig := core_transport_http_server.NewConfigMust()
	httpServer := core_transport_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_transport_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_transport_http_middleware.RequestID(),
		core_transport_http_middleware.Logger(logger),
		core_transport_http_middleware.Trace(),
		core_transport_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_transport_http_server.NewAPIVersionRouter(
		core_transport_http_server.ApiVersion1,
	)

	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(householdsTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
