package tasks_service

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
}

func NewTasksService(
	tasksRepository TasksRepository,
) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
