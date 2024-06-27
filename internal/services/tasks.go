package services

import (
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"context"
)

type TaskService struct {
	tasksRepo repos.TasksRepository
}

func NewTaskService(repo *repos.TasksRepository) *TaskService {
	return &TaskService{tasksRepo: *repo}
}

func (tr *TaskService) CreateTask(context.Context, string, int) (Task, error) {

}
func (tr *TaskService) GetTaskByID(context.Context, int) (models.Task, error) {

}

func (tr *TaskService) GetTasksByUserID(context.Context, int) ([]models.Task, error) {

}

func (tr *TaskService) DeleteTaskById(context.Context, int) error {

}

func (tr *TaskService) StartTimeTracker(context.Context, int, int) error {

}

func (tr *TaskService) StopTimeTracker(context.Context, int, int) error {

}
