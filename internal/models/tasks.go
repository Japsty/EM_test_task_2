package models

import (
	"context"
	"time"
)

type Task struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	UserID    int        `json:"user_id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

type NewTaskRequest struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

type TaskRepo interface {
	AddTask(context.Context, string, int) (Task, error)
	FindTaskByID(context.Context, int) (Task, error)
	FindTasksByUserID(context.Context, int, string, string) ([]Task, error)
	DeleteTaskByID(context.Context, int) error
	StartTimeTracker(context.Context, int, int) error
	StopTimeTracker(context.Context, int, int) error
	GetAllTasks(context.Context) ([]Task, error)
}

type TaskService interface {
	CreateTask(context.Context, string, int) (Task, error)
	GetTaskByID(context.Context, int) (Task, error)
	GetTasksByUserID(context.Context, int, string, string) ([]Task, error)
	DeleteTaskByID(context.Context, int) error
	StartTimeTracker(context.Context, int, int) error
	StopTimeTracker(context.Context, int, int) error
	GetAllTasks(context.Context) ([]Task, error)
}
