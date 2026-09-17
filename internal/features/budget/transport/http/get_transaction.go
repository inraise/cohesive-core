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

// GetTransaction godoc
// @Summary Получить транзакцию
// @Description Получить одну запись о приходе/расходе по ID
// @Tags budget
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param tx_id path string true "ID транзакции"
// @Success 200 {object} budget_transport_http.TransactionDTOResponse "Транзакция"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household or transaction not found"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/transactions/{tx_id} [get]
func (h *BudgetHTTPHandler) GetTransaction(rw http.ResponseWriter, r *http.Request) {
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

	txID, err := uuid.Parse(r.PathValue("tx_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse transaction id %q: %v: %w", r.PathValue("tx_id"), err, core_errors.ErrInvalidArgument),
			"invalid transaction id",
		)

		return
	}

	tx, err := h.budgetService.GetTransaction(ctx, householdID, txID, callerID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get transaction")

		return
	}

	response := transactionDTOFromDomain(tx)
	responseHandler.JSONResponse(response, http.StatusOK)
}
