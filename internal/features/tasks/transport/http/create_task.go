package tasks_transport_http

import (
	core_errors "cohesive-core/internal/core/errors"
	core_logger "cohesive-core/internal/core/logger"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_request "cohesive-core/internal/core/transport/http/request"
	core_transport_http_response "cohesive-core/internal/core/transport/http/response"
	tasks_service "cohesive-core/internal/features/tasks/service"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// CreateTask godoc
// @Summary Создать задачу
// @Description Создать новую задачу в доме
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "ID дома"
// @Param request body tasks_service.CreateTaskRequest true "CreateTask тело запроса"
// @Success 201 {object} tasks_transport_http.TaskDTOResponse "Задача создана"
// @Failure 400 {object} core_transport_http_response.ErrorResponse "Bad request"
// @Failure 401 {object} core_transport_http_response.ErrorResponse "Unauthorized"
// @Failure 404 {object} core_transport_http_response.ErrorResponse "Household not found"
// @Failure 500 {object} core_transport_http_response.ErrorResponse "Internal server error"
// @Router /households/{id}/tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
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

	var request tasks_service.CreateTaskRequest
	if err := core_transport_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	task, err := h.tasksService.CreateTask(ctx, householdID, callerID, request)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")

		return
	}

	response := taskDTOFromDomain(task)
	responseHandler.JSONResponse(response, http.StatusCreated)
}
