package budget_transport_http

import (
	"fmt"
	"net/http"

	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	budget_service "cohesive-core/internal/features/budget/service"

	"github.com/google/uuid"
)

// CreateTransaction godoc
// @Summary Записать приход/расход
// @Description Добавить транзакцию (deposit - пополнение общего бюджета, expense - трата из него)
// @Tags budget
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param request body budget_service.CreateTransactionRequest true "CreateTransaction тело запроса"
// @Success 201 {object} budget_transport_http.TransactionDTOResponse "Транзакция создана"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/transactions [post]
func (h *BudgetHTTPHandler) CreateTransaction(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	callerID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(fmt.Errorf("user id not found in request context"), "internal server errors")

		return
	}

	householdID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse household id %q: %v: %w", r.PathValue("id"), err, core_errors.ErrInvalidArgument),
			"invalid household id",
		)

		return
	}

	var request budget_service.CreateTransactionRequest
	if err = core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	tx, err := h.budgetService.CreateTransaction(ctx, householdID, callerID, request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create transaction")

		return
	}

	response := transactionDTOFromDomain(tx)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
