package budget_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	core_http_types "cohesive-core/internal/core/transport/http/types"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type PatchTransactionRequest struct {
	Amount      core_http_types.Nullable[int64]  `json:"amount"`
	Description core_http_types.Nullable[string] `json:"description"`
}

func (r *PatchTransactionRequest) Validate() error {
	if r.Amount.Set {
		if r.Amount.Value == nil {
			return fmt.Errorf("`Amount` can't be NULL")
		}

		if *r.Amount.Value <= 0 {
			return fmt.Errorf("`Amount` must be positive")
		}
	}

	if r.Description.Set && r.Description.Value != nil {
		descriptionLen := len([]rune(*r.Description.Value))
		if descriptionLen > 300 {
			return fmt.Errorf("`Description` must be at most 300 symbols")
		}
	}

	return nil
}

func (h *BudgetHTTPHandler) PatchTransaction(rw http.ResponseWriter, r *http.Request) {
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

	var request PatchTransactionRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	patch := core_domain.NewHouseholdTransactionPatch(
		request.Amount.ToDomain(),
		request.Description.ToDomain(),
	)

	tx, err := h.budgetService.PatchTransaction(ctx, householdID, txID, callerID, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch transaction")

		return
	}

	response := transactionDTOFromDomain(tx)
	responseHandler.JSONResponse(response, http.StatusOK)
}
