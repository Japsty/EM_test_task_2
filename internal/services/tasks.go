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

func (tr *TaskService) CreateTask(ctx context.Context, name string, usrID int) (models.Task, error) {
	task, err := tr.tasksRepo.AddTask(ctx, name, usrID)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
func (tr *TaskService) GetTaskByID(ctx context.Context, id int) (models.Task, error) {
	task, err := tr.tasksRepo.FindTaskByID(ctx, id)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (tr *TaskService) GetTasksByUserID(ctx context.Context, usrID int) ([]models.Task, error) {
	tasks, err := tr.tasksRepo.FindTasksByUserID(ctx, usrID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tr *TaskService) DeleteTaskByID(ctx context.Context, id int) error {
	err := tr.tasksRepo.DeleteTaskByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TaskService) StartTimeTracker(ctx context.Context, id int, usrID int) error {
	err := tr.tasksRepo.StartTimeTracker(ctx, id, usrID)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TaskService) StopTimeTracker(ctx context.Context, id int, usrID int) error {
	err := tr.tasksRepo.StopTimeTracker(ctx, id, usrID)
	if err != nil {
		return err
	}

	return nil
}
