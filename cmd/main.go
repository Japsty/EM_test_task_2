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

	es := services.NewEncodeService(os.Getenv("ENCRYPTION_KEY"))

	userRepo := repos.NewUsersRepository(postgreConn, es)
	taskRepo := repos.NewTasksRepository(postgreConn)

	us := services.NewUserService(userRepo)
	ts := services.NewTaskService(taskRepo)

	uh := handlers.NewUserHandler(us, logger, &client)
	th := handlers.NewTaskHandler(ts, logger)

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return middleware.AccessLog(logger, next)
	})

	r.HandleFunc("/users", uh.GetUsers).Methods("GET")
	r.HandleFunc("/user/{user_id}", uh.DeleteUser).Methods("DELETE")
	r.HandleFunc("/user/{user_id}", uh.UpdateUser).Methods("PATCH")
	r.HandleFunc("/user", uh.AddUser).Methods("POST")

	r.HandleFunc("/user/tasks", th.GetUsersTasks).Methods("GET")
	r.HandleFunc("/user/task/track/{user_id}/{task_id}", th.StartTracker).Methods("POST")
	r.HandleFunc("/user/task/stop/{user_id}/{task_id}", th.StopTracker).Methods("POST")

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
