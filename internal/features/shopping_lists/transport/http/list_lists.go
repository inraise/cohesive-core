package shoppinglists_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ListListsResponse []ShoppingListDTOResponse

// ListLists godoc
// @Summary Список списков покупок
// @Description Получить все списки покупок дома (без вложенных пунктов)
// @Tags shoppinglists
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Success 200 {array} shoppinglists_transport_http.ShoppingListDTOResponse "Списки покупок"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/shopping-lists [get]
func (h *ShoppingListsHTTPHandler) ListLists(rw http.ResponseWriter, r *http.Request) {
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

	lists, err := h.shoppingListsService.ListLists(ctx, householdID, callerID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list shopping lists")

		return
	}

	response := make(ListListsResponse, 0, len(lists))
	for _, list := range lists {
		response = append(response, shoppingListDTOFromDomain(list))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
