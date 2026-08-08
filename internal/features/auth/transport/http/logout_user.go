package auth_transport_http

import (
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *LogoutRequest) Validate() error {
	if r.RefreshToken == "" {
		return fmt.Errorf("`RefreshToken` can't be NULL")
	}

	return nil
}

func (h *AuthHTTPHandler) LogoutUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	var request LogoutRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	if err := h.authService.LogoutUser(ctx, request.RefreshToken); err != nil {
		responseHandler.ErrorResponse(
			err,
			"internal server errors",
		)

		return
	}

	responseHandler.NoContentResponse()
}
