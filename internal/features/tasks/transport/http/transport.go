package tasks_transport_http

import (
	core_domain "cohesive-core/internal/core/domain"
	core_jwt "cohesive-core/internal/core/jwt"
	core_transport_http_middleware "cohesive-core/internal/core/transport/http/middleware"
	core_transport_http_server "cohesive-core/internal/core/transport/http/server"
	tasks_service "cohesive-core/internal/features/tasks/service"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type TasksHTTPHandler struct {
	tasksService TasksService
	tokenManager *core_jwt.TokenManager
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
		request tasks_service.CreateTaskRequest,
	) (core_domain.Task, error)

	DeleteTask(
		ctx context.Context,
		householdID uuid.UUID,
		taskID uuid.UUID,
		callerID uuid.UUID,
	) error

	ListTasks(
		ctx context.Context,
		householdID uuid.UUID,
		callerID uuid.UUID,
	) ([]core_domain.Task, error)

	GetTask(
		ctx context.Context,
		householdID uuid.UUID,
		taskID uuid.UUID,
		callerID uuid.UUID,
	) (core_domain.Task, error)

	PatchTask(
		ctx context.Context,
		householdID uuid.UUID,
		taskID uuid.UUID,
		callerID uuid.UUID,
		patch core_domain.TaskPatch,
	) (core_domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
	tokenManager *core_jwt.TokenManager,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
		tokenManager: tokenManager,
	}
}

func (h *TasksHTTPHandler) Routes() []core_transport_http_server.Route {
	authenticate := core_transport_http_middleware.Authenticate(h.tokenManager)
	withAuth := []core_transport_http_middleware.Middleware{authenticate}

	return []core_transport_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/households/{id}/tasks",
			Handler:    h.CreateTask,
			Middleware: withAuth,
		},
	}
}
