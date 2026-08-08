package auth_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	"errors"
	"fmt"
	"net/http"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() error {
	if r.RefreshToken == "" {
		return fmt.Errorf("`RefreshToken` can't be NULL")
	}

	return nil
}

func (h *AuthHTTPHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	var request RefreshTokenRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	tokens, err := h.authService.RefreshToken(ctx, request.RefreshToken)
	if err != nil {
		if errors.Is(err, core_errors.ErrUnauthorized) {
			responseHandler.ErrorResponse(
				err,
				"invalid or expired refresh token",
			)

			return
		}

		responseHandler.ErrorResponse(
			err,
			"internal server errors",
		)

		return
	}

	responseHandler.JSONResponse(tokens, http.StatusOK)
}
