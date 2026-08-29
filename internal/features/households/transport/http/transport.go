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

	RenameHousehold(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		request households_service.RenameHouseholdRequest,
	) (core_domain.HouseholdWithRole, error)

	DeleteHousehold(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) error

	ListMembers(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) ([]core_domain.HouseholdMember, error)

	RemoveMember(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		targetID uuid.UUID,
	) error

	ChangeMemberRole(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		targetID uuid.UUID,
		request households_service.ChangeMemberRoleRequest,
	) error

	CreateInvite(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		request households_service.CreateInviteRequest,
	) (core_domain.HouseholdInvite, error)
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
	withAuth := []core_transport_http_middleware.Middleware{authenticate}

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
		{
			Method:     http.MethodPatch,
			Path:       "/households/{id}",
			Handler:    h.RenameHousehold,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/households/{id}",
			Handler:    h.DeleteHousehold,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/members",
			Handler:    h.ListMembers,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/households/{id}/members/{user_id}",
			Handler:    h.RemoveMember,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/households/{id}/members/{user_id}",
			Handler:    h.ChangeMemberRole,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPost,
			Path:       "/households/{id}/invites",
			Handler:    h.CreateInvite,
			Middleware: withAuth,
		},
	}
}
