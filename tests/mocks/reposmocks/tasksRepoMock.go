package reposmocks

import (
	"EMTask/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockTasksRepo struct {
	mock.Mock
}

func (tr *MockTasksRepo) AddTask(ctx context.Context, name string, usrID int) (models.Task, error) {
	args := tr.Called(ctx, name, usrID)
	return args.Get(0).(models.Task), args.Error(1)
}

func (tr *MockTasksRepo) FindTaskByID(ctx context.Context, id int) (models.Task, error) {
	args := tr.Called(ctx, id)
	return args.Get(0).(models.Task), args.Error(1)
}

func (tr *MockTasksRepo) FindTasksByUserID(ctx context.Context, usrID int, startTime, endTime string) ([]models.Task, error) {
	args := tr.Called(ctx, usrID, startTime, endTime)
	return args.Get(0).([]models.Task), args.Error(1)
}

func (tr *MockTasksRepo) DeleteTaskByID(ctx context.Context, id int) error {
	args := tr.Called(ctx, id)
	return args.Error(0)
}

func (tr *MockTasksRepo) StartTimeTracker(ctx context.Context, id, usrID int) error {
	args := tr.Called(ctx, id, usrID)
	return args.Error(0)
}

func (tr *MockTasksRepo) StopTimeTracker(ctx context.Context, id, usrID int) error {
	args := tr.Called(ctx, id, usrID)
	return args.Error(0)
}

func (tr *MockTasksRepo) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	args := tr.Called(ctx)
	return args.Get(0).([]models.Task), args.Error(1)
}
