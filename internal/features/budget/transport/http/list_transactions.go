package budget_transport_http

import (
	"fmt"
	"net/http"

	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"

	"github.com/google/uuid"
)

type ListTransactionsResponse []TransactionDTOResponse

// ListTransactions godoc
// @Summary Список транзакций
// @Description Получить все записи о приходах/расходах дома
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Success 200 {array} budget_transport_http.TransactionDTOResponse "Список транзакций"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/transactions [get]
func (h *BudgetHTTPHandler) ListTransactions(rw http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.budgetService.ListTransactions(ctx, householdID, callerID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list transactions")

		return
	}

	response := make(ListTransactionsResponse, 0, len(transactions))
	for _, tx := range transactions {
		response = append(response, transactionDTOFromDomain(tx))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
