package handlers

import (
	"EMTask/internal/models"
	"go.uber.org/zap"
)

type TaskHandler struct {
	TaskService models.TaskService
	ZapLogger   *zap.SugaredLogger
}

func NewTaskHandler(ts models.TaskService, logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{ts, logger}
}
