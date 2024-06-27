package main

import (
	"EMTask/internal/handlers"
	"EMTask/internal/middleware"
	"EMTask/internal/repos"
	"EMTask/internal/services"
	"EMTask/pkg/storage/connect"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
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

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	postgreConn, err := connect.NewPostgresConnection(os.Getenv("DSN"))
	if err != nil {
		logger.Error("Connecting to SQL database error: ", err)
		return
	}
	defer postgreConn.Close()

	es := services.NewEncodeService(os.Getenv("ENCRYPTION_KEY"))

	userRepo := repos.NewUsersRepository(postgreConn, es)
	taskRepo := repos.NewTasksRepository(postgreConn)

	us := services.NewUserService(userRepo)
	ts := services.NewTaskService(taskRepo)

	uh := handlers.NewUserHandler(us, logger)
	th := handlers.NewTaskHandler(ts, logger)

	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return middleware.AccessLog(logger, next)
	})

	addr := ":8081"
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
