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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return fmt.Errorf("`Email` and `Password` can't be NULL")
	}

	return nil
}

func (h *AuthHTTPHandler) LoginUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_transport_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	tokens, err := h.authService.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, core_errors.ErrInvalidArgument) {
			responseHandler.ErrorResponse(
				err,
				"invalid email or password",
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
