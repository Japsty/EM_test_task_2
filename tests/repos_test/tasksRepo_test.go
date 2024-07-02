package repos_test

import (
	"EMTask/internal/repos"
	"EMTask/internal/repos/queries"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAddTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(queries.ExistCheck)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(regexp.QuoteMeta(queries.CreateTask)).
		WithArgs("task name", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id"}).AddRow(1, "task name", 1))

	task, err := repo.AddTask(context.Background(), "task name", 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "task name", task.Name)
	assert.Equal(t, 1, task.UserID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(queries.FindTaskByID)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"user_id",
			"start_time",
			"end_time"}).
			AddRow(1, "task name", 1, nil, nil))

	task, err := repo.FindTaskByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindTaskByID Error: %s", err)
	}

	assert.Equal(t, 1, task.ID)
	assert.Equal(t, "task name", task.Name)
	assert.Equal(t, 1, task.UserID)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindTasksByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT id, name, user_id, start_time, end_time FROM tasks WHERE user_id = $1")).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"user_id",
			"start_time",
			"end_time"}).
			AddRow(1, "task name", 1, nil, nil))

	tasks, err := repo.FindTasksByUserID(context.Background(), 1, "", "")
	if err != nil {
		t.Fatalf("FindTaskByID Error: %s", err)
	}

	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "task name", tasks[0].Name)
	assert.Equal(t, 1, tasks[0].UserID)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(queries.DeleteTask)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteTaskByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteTaskByID Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStartTimeTracker(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectExec(
		regexp.QuoteMeta(queries.StartTimeTracker)).
		WithArgs(sqlmock.AnyArg(), 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.StartTimeTracker(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("DeleteTaskByID Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStopTimeTracker(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectExec(
		regexp.QuoteMeta(queries.StopTimeTracker)).
		WithArgs(sqlmock.AnyArg(), 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.StopTimeTracker(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("DeleteTaskByID Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAllTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewTasksRepository(db)

	mock.ExpectQuery(
		regexp.QuoteMeta(queries.GetAllTasks)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"name",
			"user_id",
			"start_time",
			"end_time"}).
			AddRow(1, "task name", 1, nil, nil))

	tasks, err := repo.GetAllTasks(context.Background())
	if err != nil {
		t.Fatalf("DeleteTaskByID Error: %s", err)
	}

	assert.Len(t, tasks, 1)
	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "task name", tasks[0].Name)
	assert.Equal(t, 1, tasks[0].UserID)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
