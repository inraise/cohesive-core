package shoppinglists_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	shoppinglists_service "cohesive-core/internal/features/shopping_lists/service"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// RenameList godoc
// @Summary Переименовать список покупок
// @Description Изменить название списка покупок
// @Tags shoppinglists
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param list_id path string true "ID списка"
// @Param request body shoppinglists_service.RenameListRequest true "RenameList тело запроса"
// @Success 200 {object} shoppinglists_transport_http.ShoppingListDTOResponse "Список переименован"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household or list not found"
// @Failure 409 {object} core_transport_http_response.ErrorResponse "Version conflict"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/shopping-lists/{list_id} [patch]
func (h *ShoppingListsHTTPHandler) RenameList(rw http.ResponseWriter, r *http.Request) {
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

	var request shoppinglists_service.RenameListRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	list, err := h.shoppingListsService.RenameList(ctx, householdID, listID, callerID, request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to rename shopping list")

		return
	}

	response := shoppingListDTOFromDomain(list)
	responseHandler.JSONResponse(response, http.StatusOK)
}
