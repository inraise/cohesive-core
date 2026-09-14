package budget_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	budget_service "cohesive-core/internal/features/budget/service"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type BudgetHTTPHandler struct {
	budgetService BudgetService
	tokenManager  *core_jwt.TokenManager
}

type BudgetService interface {
	CreateTransaction(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		request budget_service.CreateTransactionRequest,
	) (core_domain.HouseholdTransaction, error)

	GetTransaction(
		ctx context.Context,
		householdID uuid.UUID,
		txID uuid.UUID,
		callerID uuid.UUID,
	) (core_domain.HouseholdTransaction, error)

	ListTransactions(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) ([]core_domain.HouseholdTransaction, error)

	PatchTransaction(
		ctx context.Context,
		householdID uuid.UUID,
		txID uuid.UUID,
		callerID uuid.UUID,
		patch core_domain.HouseholdTransactionPatch,
	) (core_domain.HouseholdTransaction, error)

	DeleteTransaction(
		ctx context.Context,
		householdID uuid.UUID,
		txID uuid.UUID,
		callerID uuid.UUID,
	) error

	GetSummary(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) (budget_service.BudgetSummary, error)
}

func NewBudgetHTTPHandler(
	budgetService BudgetService,
	tokenManager *core_jwt.TokenManager,
) *BudgetHTTPHandler {
	return &BudgetHTTPHandler{
		budgetService: budgetService,
		tokenManager:  tokenManager,
	}
}

func (h *BudgetHTTPHandler) Routes() []core_transport_http_server.Route {
	authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)
	withAuth := []core_transport_http_middleware.Middleware{authenticate}

	return []core_transport_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/households/{id}/transactions",
			Handler:    h.CreateTransaction,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/transactions",
			Handler:    h.ListTransactions,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/transactions/{tx_id}",
			Handler:    h.GetTransaction,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodPatch,
			Path:       "/households/{id}/transactions/{tx_id}",
			Handler:    h.PatchTransaction,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodDelete,
			Path:       "/households/{id}/transactions/{tx_id}",
			Handler:    h.DeleteTransaction,
			Middleware: withAuth,
		},
		{
			Method:     http.MethodGet,
			Path:       "/households/{id}/budget",
			Handler:    h.GetSummary,
			Middleware: withAuth,
		},
	}
}
