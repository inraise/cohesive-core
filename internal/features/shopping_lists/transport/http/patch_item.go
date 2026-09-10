package shoppinglists_transport_http

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

type PatchItemRequest struct {
	Name        core_http_types.Nullable[string] `json:"name"`
	Quantity    core_http_types.Nullable[string] `json:"quantity"`
	IsPurchased core_http_types.Nullable[bool]   `json:"is_purchased"`
}

func (r *PatchItemRequest) Validate() error {
	if r.Name.Set {
		if r.Name.Value == nil {
			return fmt.Errorf("`Name` can't be NULL")
		}

		nameLen := len([]rune(*r.Name.Value))
		if nameLen < 1 || nameLen > 200 {
			return fmt.Errorf("`Name` must be between 1 and 200 symbols")
		}
	}

	if r.Quantity.Set && r.Quantity.Value != nil {
		quantityLen := len([]rune(*r.Quantity.Value))
		if quantityLen > 50 {
			return fmt.Errorf("`Quantity` must be at most 50 symbols")
		}
	}

	if r.IsPurchased.Set && r.IsPurchased.Value == nil {
		return fmt.Errorf("`IsPurchased` can't be NULL")
	}

	return nil
}

func (h *ShoppingListsHTTPHandler) PatchItem(rw http.ResponseWriter, r *http.Request) {
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

	listID, err := uuid.Parse(r.PathValue("list_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse list id %q: %v: %w", r.PathValue("list_id"), err, core_errors.ErrInvalidArgument),
			"invalid list id",
		)

		return
	}

	itemID, err := uuid.Parse(r.PathValue("item_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse item id %q: %v: %w", r.PathValue("item_id"), err, core_errors.ErrInvalidArgument),
			"invalid item id",
		)

		return
	}

	var request PatchItemRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	patch := core_domain.NewShoppingListItemPatch(
		request.Name.ToDomain(),
		request.Quantity.ToDomain(),
		request.IsPurchased.ToDomain(),
	)

	item, err := h.shoppingListsService.PatchItem(ctx, householdID, listID, itemID, callerID, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch item")

		return
	}

	response := itemDTOFromDomain(item)
	responseHandler.JSONResponse(response, http.StatusOK)
}
