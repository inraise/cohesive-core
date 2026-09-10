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

// CreateList godoc
// @Summary Создать список покупок
// @Description Создать новый именованный список покупок в доме
// @Tags shoppinglists
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param request body shoppinglists_service.CreateListRequest true "CreateList тело запроса"
// @Success 201 {object} shoppinglists_transport_http.ShoppingListDTOResponse "Список создан"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/shopping-lists [post]
func (h *ShoppingListsHTTPHandler) CreateList(rw http.ResponseWriter, r *http.Request) {
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

	var request shoppinglists_service.CreateListRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	list, err := h.shoppingListsService.CreateList(ctx, householdID, callerID, request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create shopping list")

		return
	}

	response := shoppingListDTOFromDomain(list)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
