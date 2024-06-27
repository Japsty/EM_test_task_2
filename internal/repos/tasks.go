package repos

import (
	"EMTask/internal/models"
	"EMTask/internal/repos/queries"
	"context"
	"database/sql"
	"time"
)

type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (tr *TasksRepository) AddTask(ctx context.Context, name string, usrID int) (models.Task, error) {
	var task models.Task

	err := tr.db.QueryRowContext(ctx, queries.CreateTask, name, usrID).Scan(
		&task.ID,
		&task.Name,
		&task.UserID,
		&task.StartTime,
	)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
func (tr *TasksRepository) FindTaskByID(ctx context.Context, id int) (models.Task, error) {
	var task models.Task

	err := tr.db.QueryRowContext(ctx, queries.FindTaskByID, id).Scan(
		&task.ID,
		&task.Name,
		&task.UserID,
		&task.StartTime,
		&task.EndTime,
	)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (tr *TasksRepository) FindTasksByUserID(ctx context.Context, usrID int) ([]models.Task, error) {
	rows, err := tr.db.QueryContext(ctx, queries.FindTasksByUserID, usrID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Name, &task.UserID, &task.StartTime, &task.EndTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (tr *TasksRepository) DeleteTaskByID(ctx context.Context, id int) error {
	_, err := tr.db.ExecContext(ctx, queries.DeleteTask, id)
	return err
}

func (tr *TasksRepository) StartTimeTracker(ctx context.Context, id, usrID int) error {
	_, err := tr.db.ExecContext(ctx, queries.StartTimeTracker, time.Now(), id, usrID)
	return err
}

func (tr *TasksRepository) StopTimeTracker(ctx context.Context, id, usrID int) error {
	_, err := tr.db.ExecContext(ctx, queries.StopTimeTracker, time.Now(), id, usrID)
	return err
}
