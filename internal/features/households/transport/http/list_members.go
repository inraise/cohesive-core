package households_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type ListMembersResponse []MemberDTOResponse

// ListMembers godoc
// @Summary Список участников дома
// @Description Получить список участников дома с ролями. Доступно любому участнику
// @Tags households
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID дома"
// @Success 200 {array} households_transport_http.MemberDTOResponse "Список участников"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found or not a member"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/members [get]
func (h *HouseholdsHTTPHandler) ListMembers(rw http.ResponseWriter, r *http.Request) {
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

	members, err := h.householdsService.ListMembers(ctx, householdID, callerID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to list members")

		return
	}

	response := make(ListMembersResponse, 0, len(members))
	for _, member := range members {
		response = append(response, memberDTOFromDomain(member))
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
