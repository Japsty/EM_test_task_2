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

// @Summary Create a new task
// @Description Создание новой задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.NewTaskRequest true "New Task"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [post]
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
		if errors.Is(err, repos.ErrUsrNotExists) {
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

// @Summary Get task by ID
// @Description Получение задачи по ID
// @Tags tasks
// @Produce json
// @Param task_id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Invalid task_id"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{task_id} [get]
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
		if errors.Is(err, sql.ErrNoRows) {
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

// @Summary Delete task by ID
// @Description Удаление задачи по ID
// @Tags tasks
// @Param task_id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid task_id"
// @Failure 500 {string} string "Internal server error"
// @Router /tasks/{task_id} [delete]
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

// @Summary Get tasks by user
// @Description Получение задач юзера по его id с сортировкой по трудозатратам
// @Tags tasks
// @Produce json
// @Param user_id query int true "User ID"
// @Param start_time query string false "Start Time (RFC3339 format)"
// @Param end_time query string false "End Time (RFC3339 format)"
// @Success 200 {array} models.Task
// @Failure 400 {string} string "Invalid user_id"
// @Failure 400 {string} string "Invalid time format"
// @Failure 500 {string} string "Internal server error"
// @Router /user/tasks [get]
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
		if _, err = time.Parse(time.RFC3339, startTime); err != nil {
			th.ZapLogger.Infof(reqIDString+" GetUsersTasks Invalid start_time: ", err)
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)

			return
		}
	}

	if endTime != "" {
		if _, err = time.Parse(time.RFC3339, endTime); err != nil {
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

// @Summary Start task tracker
// @Description Запуск таймера на задачу юзера
// @Tags tasks
// @Param user_id path int true "User ID"
// @Param task_id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid user_id or task_id"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/task/track/{user_id}/{task_id} [post]
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
		if errors.Is(err, repos.ErrTaskNotFound) {
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

// @Summary Stop task tracker
// @Description Остановка таймера по задаче юзера
// @Tags tasks
// @Param user_id path int true "User ID"
// @Param task_id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid user_id or task_id"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {string} string "Internal server error"
// @Router /user/task/stop/{user_id}/{task_id} [post]
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
		if errors.Is(err, repos.ErrTaskNotFound) {
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

// @Summary Get all tasks
// @Description Получение списка всех задач
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Internal server error"
// @Router /tasks [get]
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
