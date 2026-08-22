package households_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	households_service "cohesive-core/internal/features/households/service"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type HouseholdsHTTPHandler struct {
	householdsService HouseholdsService
	tokenManager      *core_jwt.TokenManager
}

type HouseholdsService interface {
	CreateHousehold(
		ctx context.Context,
		ownerID uuid.UUID,
		request households_service.CreateHouseholdRequest,
	) (core_domain.Household, error)

	ListMyHouseholds(
		ctx context.Context,
		userID uuid.UUID,
	) ([]core_domain.HouseholdWithRole, error)

	GetHousehold(
		ctx context.Context,
		householdID uuid.UUID,
		userID uuid.UUID,
	) (core_domain.HouseholdWithRole, error)
}

func NewHouseholdsHTTPHandler(
	householdsService HouseholdsService,
	tokenManager *core_jwt.TokenManager,
) *HouseholdsHTTPHandler {
	return &HouseholdsHTTPHandler{
		householdsService: householdsService,
		tokenManager:      tokenManager,
	}
}

func (h *HouseholdsHTTPHandler) Routes() []core_transport_http_server.Route {
	authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)

	return []core_transport_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/households",
			Handler: h.CreateHousehold,
			Middleware: []core_transport_http_middleware.Middleware{
				authenticate,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/households",
			Handler: h.ListHouseholds,
			Middleware: []core_transport_http_middleware.Middleware{
				authenticate,
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/households/{id}",
			Handler: h.GetHousehold,
			Middleware: []core_transport_http_middleware.Middleware{
				authenticate,
			},
		},
	}
}
