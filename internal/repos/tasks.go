package repos

import (
	"EMTask/internal/models"
	"context"
	"database/sql"
)

type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (tr *TasksRepository) AddTask(context.Context, string, int) (models.Task, error) {

}
func (tr *TasksRepository) FindTaskById(context.Context, int) (models.Task, error) {

}

func (tr *TasksRepository) FindTasksByUserID(context.Context, int) ([]models.Task, error) {

}

func (tr *TasksRepository) DeleteTaskById(context.Context, int) error {

}

func (tr *TasksRepository) StartTimeTracker(context.Context, int, int) error {

}

func (tr *TasksRepository) StopTimeTracker(context.Context, int, int) error {

}
