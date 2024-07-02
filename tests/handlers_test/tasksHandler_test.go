package handlers_test

import (
	"EMTask/internal/handlers"
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"EMTask/internal/services"
	"EMTask/tests/mocks/reposmocks"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var mockTask = models.Task{
	ID:        1,
	Name:      "написать тестовое",
	UserID:    1,
	StartTime: nil,
	EndTime:   nil,
}

func TestCreateTask(t *testing.T) {
	type mockRepoResp struct {
		task      models.Task
		mockError error
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		taskReq        models.NewTaskRequest
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/tasks",
				mockRequestBody: strings.NewReader(`{
					"name":"написать тестовое",
					"user_id":1
				}`),
			},
			taskReq: models.NewTaskRequest{Name: "написать тестовое", UserID: 1},
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Decode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/tasks",
				mockRequestBody:   strings.NewReader(`{Это я сломал decode}`),
			},
			taskReq:        models.NewTaskRequest{Name: "написать тестовое", UserID: 1},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/tasks",
				mockRequestBody: strings.NewReader(`{
					"name":"написать тестовое",
					"user_id":1
				}`),
			},
			taskReq: models.NewTaskRequest{Name: "написать тестовое", UserID: 1},
			repoResp: mockRepoResp{
				task:      models.Task{},
				mockError: errors.New("эта ошибка ломает service"),
			},
			callRepo:       false,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   4,
			name: "Service UserNotExists Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/tasks",
				mockRequestBody: strings.NewReader(`{
					"name":"написать тестовое",
					"user_id":1
				}`),
			},
			taskReq: models.NewTaskRequest{Name: "написать тестовое", UserID: 1},
			repoResp: mockRepoResp{
				task:      models.Task{},
				mockError: repos.ErrUsrNotExists,
			},
			callRepo:       true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   5,
			name: "Encode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/tasks",
				mockRequestBody: strings.NewReader(`{
					"name":"написать тестовое",
					"user_id":1
				}`),
			},
			taskReq: models.NewTaskRequest{Name: "написать тестовое", UserID: 1},
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: nil,
			},
			breakWrite:     true,
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On("AddTask", mock.AnythingOfType("*context.timerCtx"), tc.taskReq.Name, tc.taskReq.UserID).Return(tc.repoResp.task, tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			if tc.breakWrite {
				taskHandler.CreateTask(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				taskHandler.CreateTask(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(t, "AddTask", mock.Anything, tc.taskReq.Name, tc.taskReq.UserID)
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	type mockRepoResp struct {
		task      models.Task
		mockError error
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		reqUserID      int
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/asfsaf",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Service NotFound error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: sql.ErrNoRows,
			},
			callRepo:       true,
			expectedStatus: http.StatusNotFound,
		},
		{
			id:   4,
			name: "Service error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: errors.New("эта ошибка ломает сервис"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   5,
			name: "Encode error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				task:      mockTask,
				mockError: nil,
			},
			breakWrite:     true,
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On("FindTaskByID", mock.AnythingOfType("*context.timerCtx"), tc.reqUserID).Return(tc.repoResp.task, tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/tasks/{task_id}", taskHandler.GetTaskByID).Methods(http.MethodGet)

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(t, "FindTaskByID", mock.Anything, tc.reqUserID)
			}
		})
	}
}

func TestDeleteTaskByID(t *testing.T) {
	type mockRepoResp struct {
		mockError error
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		reqUserID      int
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusNoContent,
		},
		{
			id:   2,
			name: "Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/asfsaf",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Service error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks/1",
				mockRequestBody:   strings.NewReader(``),
			},
			reqUserID: 1,
			repoResp: mockRepoResp{
				mockError: errors.New("эта ошибка ломает сервис"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On("DeleteTaskByID", mock.AnythingOfType("*context.timerCtx"), tc.reqUserID).Return(tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/tasks/{task_id}", taskHandler.DeleteTaskByID).Methods(http.MethodGet)

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(t, "DeleteTaskByID", mock.Anything, tc.reqUserID)
			}
		})
	}
}

func TestGetUsersTasks(t *testing.T) {
	type mockRepoResp struct {
		tasks     []models.Task
		mockError error
	}

	type findTasksReq struct {
		usrID     int
		startTime string
		endTime   string
	}

	mockTask1 := models.Task{
		ID:        2,
		Name:      "mockTask1",
		UserID:    1,
		StartTime: nil,
		EndTime:   nil,
	}
	mockTask2 := models.Task{
		ID:        3,
		Name:      "mockTask2",
		UserID:    2,
		StartTime: nil,
		EndTime:   nil,
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		mockFindReq    findTasksReq
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?user_id=1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockFindReq: findTasksReq{
				usrID:     1,
				startTime: "",
				endTime:   "",
			},
			repoResp: mockRepoResp{
				tasks:     []models.Task{mockTask, mockTask1},
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?user_id=safasf",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Invalid start_time Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?start_time=dsfdsf&user_id=1",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   4,
			name: "Invalid end_time Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?end_time=fgfgfg&user_id=1",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   5,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?user_id=1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockFindReq: findTasksReq{
				usrID:     1,
				startTime: "",
				endTime:   "",
			},
			repoResp: mockRepoResp{
				tasks:     nil,
				mockError: errors.New("эта ошибка ломает service"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   6,
			name: "Encode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/user/tasks?user_id=1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockFindReq: findTasksReq{
				usrID:     1,
				startTime: "",
				endTime:   "",
			},
			repoResp: mockRepoResp{
				tasks:     []models.Task{mockTask, mockTask1, mockTask2},
				mockError: nil,
			},
			breakWrite:     true,
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On(
				"FindTasksByUserID",
				mock.AnythingOfType("*context.timerCtx"),
				tc.mockFindReq.usrID,
				tc.mockFindReq.startTime,
				tc.mockFindReq.endTime,
			).Return(tc.repoResp.tasks, tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			if tc.breakWrite {
				taskHandler.GetUsersTasks(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				taskHandler.GetUsersTasks(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(
					t,
					"FindTasksByUserID",
					mock.Anything,
					tc.mockFindReq.usrID,
					tc.mockFindReq.startTime,
					tc.mockFindReq.endTime,
				)
			}
		})
	}
}

func TestStartTracker(t *testing.T) {
	type mockRepoResp struct {
		mockError error
	}

	type mockReqParam struct {
		taskID int
		usrID  int
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		mockReqParams  mockReqParam
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusNoContent,
		},
		{
			id:   2,
			name: "Atoi UserID Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/safafs/1",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Atoi TaskID Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/sadasd",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   4,
			name: "Service NotFound Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: repos.ErrTaskNotFound,
			},
			callRepo:       true,
			expectedStatus: http.StatusNotFound,
		},
		{
			id:   5,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: errors.New("эта ошибка ломает service"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On(
				"StartTimeTracker",
				mock.AnythingOfType("*context.timerCtx"),
				tc.mockReqParams.taskID,
				tc.mockReqParams.usrID,
			).Return(tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/user/task/track/{user_id}/{task_id}", taskHandler.StartTracker).Methods(http.MethodPost)

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(
					t,
					"StartTimeTracker",
					mock.Anything,
					tc.mockReqParams.taskID,
					tc.mockReqParams.usrID,
				)
			}
		})
	}
}

func TestStopTracker(t *testing.T) {
	type mockRepoResp struct {
		mockError error
	}

	type mockReqParam struct {
		taskID int
		usrID  int
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		mockReqParams  mockReqParam
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusNoContent,
		},
		{
			id:   2,
			name: "Atoi UserID Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/safafs/1",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Atoi TaskID Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/sadasd",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   4,
			name: "Service NotFound Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: repos.ErrTaskNotFound,
			},
			callRepo:       true,
			expectedStatus: http.StatusNotFound,
		},
		{
			id:   5,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user/task/track/1/1",
				mockRequestBody:   strings.NewReader(``),
			},
			mockReqParams: mockReqParam{
				taskID: 1,
				usrID:  1,
			},
			repoResp: mockRepoResp{
				mockError: errors.New("эта ошибка ломает service"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On(
				"StopTimeTracker",
				mock.AnythingOfType("*context.timerCtx"),
				tc.mockReqParams.taskID,
				tc.mockReqParams.usrID,
			).Return(tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/user/task/track/{user_id}/{task_id}", taskHandler.StopTracker).Methods(http.MethodPost)

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(
					t,
					"StopTimeTracker",
					mock.Anything,
					tc.mockReqParams.taskID,
					tc.mockReqParams.usrID,
				)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	type mockRepoResp struct {
		tasks     []models.Task
		mockError error
	}

	mockTask1 := models.Task{
		ID:        2,
		Name:      "mockTask1",
		UserID:    1,
		StartTime: nil,
		EndTime:   nil,
	}
	mockTask2 := models.Task{
		ID:        3,
		Name:      "mockTask2",
		UserID:    2,
		StartTime: nil,
		EndTime:   nil,
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		repoResp       mockRepoResp
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				tasks:     []models.Task{mockTask, mockTask1},
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				tasks:     nil,
				mockError: errors.New("эта ошибка ломает service"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   3,
			name: "Encode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/tasks",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				tasks:     []models.Task{mockTask, mockTask1, mockTask2},
				mockError: nil,
			},
			breakWrite:     true,
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}

			logger := zapLogger.Sugar()

			mockTasksRepo := new(reposmocks.MockTasksRepo)

			mockTaskService := services.NewTaskService(mockTasksRepo)

			taskHandler := handlers.NewTaskHandler(mockTaskService, logger)

			mockTasksRepo.On(
				"GetAllTasks",
				mock.AnythingOfType("*context.timerCtx"),
			).Return(tc.repoResp.tasks, tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			if tc.breakWrite {
				taskHandler.GetAllTasks(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				taskHandler.GetAllTasks(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockTasksRepo.AssertCalled(
					t,
					"GetAllTasks",
					mock.Anything,
				)
			}
		})
	}
}
