package handlers

import (
	"EMTask/internal/models"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type TaskHandler struct {
	TaskService models.TaskService
	ZapLogger   *zap.SugaredLogger
}

func NewTaskHandler(ts models.TaskService, logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{ts, logger}
}

func (th *TaskHandler) GetUsersTasks(w http.ResponseWriter, r *http.Request) {

}

func (th *TaskHandler) StartTracker(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = th.TaskService.StartTimeTracker(r.Context(), task.ID, task.UserID)
	if err != nil {
		http.Error(w, "Error starting tracker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) StopTracker(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = th.TaskService.StopTimeTracker(r.Context(), task.ID, task.UserID)
	if err != nil {
		http.Error(w, "Error stopping tracker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
