package main

import (
	"EMTask/internal/handlers"
	"EMTask/internal/middleware"
	"EMTask/internal/repos"
	"EMTask/internal/services"
	"EMTask/pkg/storage/connect"
	"EMTask/pkg/storage/migrate"
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

// @title Time Tracker
// @version 1.0
// @description RESTful Time Tracker for EM
func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return
	}

	defer func(zapLogger *zap.Logger) {
		err = zapLogger.Sync()
		if err != nil {
			return
		}
	}(zapLogger)

	logger := zapLogger.Sugar()

	postgreConn, err := connect.NewPostgresConnection(os.Getenv("DSN"))
	if err != nil {
		logger.Error("Connecting to SQL database error: ", err)
		return
	}
	defer postgreConn.Close()

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = migrate.UpMigration(ctxWithTimeout, postgreConn)
	if err != nil {
		logger.Fatal("Failed to up migration: ", err)
	}

	client := http.Client{
		Timeout: time.Second,
	}

	userRepo := repos.NewUsersRepository(postgreConn)
	taskRepo := repos.NewTasksRepository(postgreConn)

	us := services.NewUserService(userRepo)
	ts := services.NewTaskService(taskRepo)

	uh := handlers.NewUserHandler(us, logger, &client)
	th := handlers.NewTaskHandler(ts, logger)

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return middleware.AccessLog(logger, next)
	})

	r.HandleFunc("/users", uh.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/user/{user_id}", uh.DeleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/user/{user_id}", uh.UpdateUser).Methods(http.MethodPatch)
	r.HandleFunc("/user", uh.AddUser).Methods(http.MethodPost)

	r.HandleFunc("/tasks", th.CreateTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{task_id}", th.GetTaskByID).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{task_id}", th.DeleteTaskByID).Methods(http.MethodDelete)
	r.HandleFunc("/user/tasks", th.GetUsersTasks).Methods(http.MethodGet)
	r.HandleFunc("/user/task/track/{user_id}/{task_id}", th.StartTracker).Methods(http.MethodPost)
	r.HandleFunc("/user/task/stop/{user_id}/{task_id}", th.StopTracker).Methods(http.MethodPost)
	r.HandleFunc("/tasks", th.GetAllTasks).Methods(http.MethodGet)

	addr := ":" + os.Getenv("PORT")
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Error("main ListenAndServe error: ", err)
		return
	}
}
