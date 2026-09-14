package budget_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

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
