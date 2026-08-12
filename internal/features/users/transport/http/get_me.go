package users_transport_http

import (
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type GetMeResponse UserDTOResponse

func (h *UsersHTTPHandler) GetMe(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	userID, ok := core_transport_http_middleware.UserIDFromContext(ctx)
	if !ok {
		responseHandler.ErrorResponse(
			fmt.Errorf("user id not found in request context"),
			"internal server errors",
		)

		return
	}

	userDomain, err := h.usersService.GetMe(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")

		return
	}

	response := userDTOFromDomain(userDomain)
	responseHandler.JSONResponse(response, http.StatusOK)
}