package repos

import (
	"EMTask/internal/models"
	"EMTask/internal/repos/queries"
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrUsrNotExists = errors.New("user not exists")

type TasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{db: db}
}

func (tr *TasksRepository) AddTask(ctx context.Context, name string, usrID int) (models.Task, error) {
	var task models.Task

	var exists bool
	err := tr.db.QueryRowContext(ctx, queries.ExistCheck, usrID).Scan(&exists)

	if err != nil {
		return models.Task{}, err
	}

	if !exists {
		return models.Task{}, ErrUsrNotExists
	}

	err = tr.db.QueryRowContext(ctx, queries.CreateTask, name, usrID).Scan(
		&task.ID,
		&task.Name,
		&task.UserID,
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

func (tr *TasksRepository) FindTasksByUserID(ctx context.Context, usrID int, startTime, endTime string) ([]models.Task, error) {
	query := squirrel.Select("id", "name", "user_id", "start_time", "end_time").
		From("tasks").
		Where(squirrel.Eq{"user_id": usrID})

	if startTime != "" {
		query = query.Where(squirrel.GtOrEq{"start_time": startTime})
	}

	if endTime != "" {
		query = query.Where(squirrel.LtOrEq{"end_time": endTime})
	}

	if startTime != "" && endTime != "" {
		query = query.OrderBy("AGE(end_time, start_time) DESC")
	}

	query = query.PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tr.db.QueryContext(ctx, sqlQuery, args...)
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
	res, err := tr.db.ExecContext(ctx, queries.StartTimeTracker, time.Now(), id, usrID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (tr *TasksRepository) StopTimeTracker(ctx context.Context, id, usrID int) error {
	res, err := tr.db.ExecContext(ctx, queries.StopTimeTracker, time.Now(), id, usrID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (tr *TasksRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	rows, err := tr.db.QueryContext(ctx, queries.GetAllTasks)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.UserID,
			&task.StartTime,
			&task.EndTime,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
