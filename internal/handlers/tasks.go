package handlers

import (
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type TaskHandler struct {
	TaskService models.TaskService
	ZapLogger   *zap.SugaredLogger
}

func NewTaskHandler(ts models.TaskService, logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{ts, logger}
}

func (th *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	var newTaskRequest models.NewTaskRequest

	err := json.NewDecoder(r.Body).Decode(&newTaskRequest)
	if err != nil {
		th.ZapLogger.Error(reqIDString+"CreateTask Decode Error, caused by: ", r.Body)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := th.TaskService.CreateTask(ctxWthTimeout, newTaskRequest.Name, newTaskRequest.UserID)
	if err != nil {
		if errors.As(err, &repos.ErrUsrNotExists) {
			th.ZapLogger.Error(reqIDString+"CreateTask Error: ", err)
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
		th.ZapLogger.Error(reqIDString+"CreateTask Service Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		th.ZapLogger.Error(reqIDString+"CreateTask Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" GetTaskByID Invalid task_id: ", err)
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	task, err := th.TaskService.GetTaskByID(ctxWthTimeout, taskID)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			th.ZapLogger.Infof(reqIDString+" GetTaskByID Not Found: ", err)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		th.ZapLogger.Error(reqIDString+" GetTaskByID TaskService Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		th.ZapLogger.Error(reqIDString+" GetTaskByID Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" DeleteTaskByID Invalid task_id: ", err)
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	err = th.TaskService.DeleteTaskByID(ctxWthTimeout, taskID)
	if err != nil {
		th.ZapLogger.Error(reqIDString+" DeleteTaskByID TaskService Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) GetUsersTasks(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	usrIDStr := r.URL.Query().Get("user_id")
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	usrID, err := strconv.Atoi(usrIDStr)
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" GetUsersTasks Invalid user_id: ", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if startTime != "" {
		if _, err := time.Parse(time.RFC3339, startTime); err != nil {
			th.ZapLogger.Infof(reqIDString+" GetUsersTasks Invalid start_time: ", err)
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
	}
	if endTime != "" {
		if _, err := time.Parse(time.RFC3339, endTime); err != nil {
			th.ZapLogger.Infof(reqIDString+" GetUsersTasks Invalid end_time: ", err)
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	}

	tasks, err := th.TaskService.GetTasksByUserID(ctxWthTimeout, usrID, startTime, endTime)
	if err != nil {
		th.ZapLogger.Error(reqIDString+" GetUsersTasks Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		th.ZapLogger.Error(reqIDString+" GetUsersTasks Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) StartTracker(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" StartTracker Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" StartTracker Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid task_id", http.StatusBadRequest)
		return
	}

	err = th.TaskService.StartTimeTracker(ctxWthTimeout, taskID, userID)
	if err != nil {
		if errors.As(err, &repos.ErrTaskNotFound) {
			th.ZapLogger.Infof(reqIDString+" StartTimeTracker TaskNotFound: ", err)
			http.Error(w, "Task not Found", http.StatusNotFound)
			return
		}
		th.ZapLogger.Error(reqIDString+" StartTimeTracker Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) StopTracker(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" StopTracker Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
	if err != nil {
		th.ZapLogger.Infof(reqIDString+" StopTracker Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = th.TaskService.StopTimeTracker(ctxWthTimeout, taskID, userID)
	if err != nil {
		if errors.As(err, &repos.ErrTaskNotFound) {
			th.ZapLogger.Infof(reqIDString+" StopTimeTracker TaskNotFound: ", err)
			http.Error(w, "Task not Found", http.StatusNotFound)
			return
		}
		th.ZapLogger.Error(reqIDString+" StopTimeTracker Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	users, err := th.TaskService.GetAllTasks(ctxWthTimeout)
	if err != nil {
		th.ZapLogger.Error(reqIDString+"GetAllTasks Error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		th.ZapLogger.Error(reqIDString+"GetAllTasks Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
