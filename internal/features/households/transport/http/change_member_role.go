package households_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	households_service "cohesive-core/internal/features/households/service"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// ChangeMemberRole godoc
// @Summary Сменить роль участника / передать владение
// @Description Меняет роль участника на admin/member, либо (role=owner) атомарно передаёт владение домом. Доступно только owner, нельзя применить к себе
// @Tags households
// @Accept json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Param user_id path string true "ID пользователя"
// @Param request body households_service.ChangeMemberRoleRequest true "ChangeMemberRole тело запроса"
// @Success 204 "Роль изменена"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request - can't target self"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 403 {object} core_transport_http_response.ErrorResponse "Forbidden - only owner can change roles"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 409 {object} core_transport_http_response.ErrorResponse "Ownership transfer conflict"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/members/{user_id} [patch]
func (h *HouseholdsHTTPHandler) ChangeMemberRole(rw http.ResponseWriter, r *http.Request) {
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

	targetID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		responseHandler.ErrorResponse(
			fmt.Errorf("parse user id %q: %v: %w", r.PathValue("user_id"), err, core_errors.ErrInvalidArgument),
			"invalid user id",
		)

		return
	}

	var request households_service.ChangeMemberRoleRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	if err := h.householdsService.ChangeMemberRole(ctx, householdID, callerID, targetID, request); err != nil {
		responseHandler.ErrorResponse(err, "failed to change member role")

		return
	}

	responseHandler.NoContentResponse()
}
