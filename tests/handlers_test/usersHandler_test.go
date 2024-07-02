package handlers_test

import (
	"EMTask/internal/handlers"
	"EMTask/internal/models"
	"EMTask/internal/services"
	"EMTask/tests/mocks/reposmocks"
	"bytes"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type errorResponseWriter struct {
	httptest.ResponseRecorder
}

func (e *errorResponseWriter) Write([]byte) (int, error) {
	return 500, errors.New("Mock error")
}

type mockRequest struct {
	mockRequestMethod        string
	mockRequestURL           string
	mockRequestBody          *strings.Reader
	mockIncorrectRequestBody *bytes.Buffer
}

var mockServiceUser = models.ServiceUser{
	PassportNum: "1234 567890",
	Surname:     "Иванов",
	Name:        "Иван",
	Patronymic:  "Иванович",
	Address:     "г. Москва, ул. Ленина, д. 5, кв. 1",
}

var mockUser = models.User{
	ID:             1,
	PassportNumber: "1234 567890",
	Surname:        "Иванов",
	Name:           "Иван",
	Patronymic:     "Иванович",
	Address:        "г. Москва, ул. Ленина, д. 5, кв. 1",
}

var mockAPIUser = models.APIResponse{
	Surname:    "Викторов",
	Name:       "Виктор",
	Patronymic: "Викторович",
	Address:    "г.Санкт-Петербург",
}

func TestAddUser(t *testing.T) {
	type mockRepoResp struct {
		usrID     int
		mockError error
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		repoResp       mockRepoResp
		apiURL         string
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{"passportNumber": "1234 567890"}`),
			},
			repoResp: mockRepoResp{
				usrID:     1,
				mockError: nil,
			},
			apiURL:         "http://0.0.0.0:4010",
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Bad Input Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{"passportNumber": "bad input"}`),
			},
			apiURL:         "http://0.0.0.0:4010",
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "GetInfoError",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{"passportNumber": "1234 567890"}`),
			},
			apiURL:         "http://0.0.0.0:4011",
			callRepo:       false,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   4,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{"passportNumber": "1234 567890"}`),
			},
			repoResp: mockRepoResp{
				mockError: errors.New("Эта ошибка ломает сервис"),
			},
			apiURL:         "http://0.0.0.0:4010",
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   5,
			name: "Decode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{Вот эти слова сломают DECODE}`),
			},
			apiURL:         "http://0.0.0.0:4010",
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   6,
			name: "Encode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPost,
				mockRequestURL:    "/user",
				mockRequestBody:   strings.NewReader(`{"passportNumber": "1234 567890"}`),
			},
			apiURL:         "http://0.0.0.0:4010",
			callRepo:       false,
			breakWrite:     true,
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

			os.Setenv("API_URL", tc.apiURL)

			mockUserRepo := new(reposmocks.MockUserRepo)

			mockUserService := services.NewUserService(mockUserRepo)

			client := &http.Client{}

			userHandler := handlers.NewUserHandler(mockUserService, logger, client)

			mockUserRepo.On("AddUser", mock.AnythingOfType("*context.timerCtx"), mockServiceUser).Return(tc.repoResp.usrID, tc.repoResp.mockError)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			if tc.breakWrite {
				userHandler.AddUser(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				userHandler.AddUser(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockUserRepo.AssertCalled(t, "AddUser", mock.Anything, mockServiceUser)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	type mockRepoResp struct {
		users []models.User
		err   error
	}
	mockUser1 := mockUser
	mockUser2 := mockUser
	mockUser3 := mockUser

	mockUser2.PassportNumber = "2234 567890"
	mockUser2.ID = 2
	mockUser2.Name = "Владимир"
	mockUser2.Surname = "Владимиров"
	mockUser2.Address = "г.Жуковский"

	mockUser3.PassportNumber = "3234 567890"
	mockUser3.ID = 3
	mockUser3.Name = "Илья"
	mockUser3.Surname = "Ильин"
	mockUser3.Patronymic = "Ильич"
	mockUser3.Address = "г.Санкт-Петербург"

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		mockPageNum    int
		mockLimitNum   int
		mockFilter     models.UserFilter
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
				mockRequestURL:    "/users?page=1&limit=10",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{},
			repoResp: mockRepoResp{
				users: []models.User{mockUser1, mockUser2, mockUser3},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Page Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?page=0&limit=10",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:    0,
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Limit Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?page=1&limit=0",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:    1,
			mockLimitNum:   0,
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   4,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?page=1&limit=10",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{},
			repoResp: mockRepoResp{
				users: nil,
				err:   errors.New("Эта ошибка ломает сервис"),
			},
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   5,
			name: "Encode Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?page=1&limit=10",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{},
			repoResp: mockRepoResp{
				users: []models.User{mockUser1, mockUser2, mockUser3},
				err:   nil,
			},
			callRepo:       true,
			breakWrite:     true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   6,
			name: "PassportNum Filter Check",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?page=1&limit=10&passport=1234 567890",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{PassportNum: "1234 567890"},
			repoResp: mockRepoResp{
				users: []models.User{mockUser1},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   7,
			name: "Surname Filter Check",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?surname=Владимиров&page=1&limit=10&",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{Surname: "Владимиров"},
			repoResp: mockRepoResp{
				users: []models.User{mockUser2},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   8,
			name: "Name Filter Check",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?name=Илья&page=1&limit=10&",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{Name: "Илья"},
			repoResp: mockRepoResp{
				users: []models.User{mockUser3},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   9,
			name: "Patronymic Filter Check",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?patronymic=Иванович&page=1&limit=10&",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{Patronymic: "Иванович"},
			repoResp: mockRepoResp{
				users: []models.User{mockUser1, mockUser2},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   10,
			name: "Address Filter Check",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodGet,
				mockRequestURL:    "/users?address=г. Москва, ул. Ленина, д. 5, кв. 1&page=1&limit=10&",
				mockRequestBody:   strings.NewReader(""),
			},
			mockPageNum:  1,
			mockLimitNum: 10,
			mockFilter:   models.UserFilter{Address: "г. Москва, ул. Ленина, д. 5, кв. 1"},
			repoResp: mockRepoResp{
				users: []models.User{mockUser1},
				err:   nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				t.Fatal(err)
			}
			logger := zapLogger.Sugar()

			mockUserRepo := new(reposmocks.MockUserRepo)

			mockUserService := services.NewUserService(mockUserRepo)

			client := &http.Client{}

			userHandler := handlers.NewUserHandler(mockUserService, logger, client)

			mockUserRepo.On("GetAllUsers", mock.AnythingOfType("*context.timerCtx"), tc.mockFilter, tc.mockPageNum, tc.mockLimitNum).Return(tc.repoResp.users, tc.repoResp.err)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			if tc.breakWrite {
				userHandler.GetUsers(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				userHandler.GetUsers(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockUserRepo.AssertCalled(t, "GetAllUsers", mock.Anything, tc.mockFilter, tc.mockPageNum, tc.mockLimitNum)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type mockRepoResp struct {
		err error
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		repoResp       mockRepoResp
		mockUsrID      int
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodDelete,
				mockRequestURL:    "/user/1",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				err: nil,
			},
			mockUsrID:      1,
			callRepo:       true,
			expectedStatus: http.StatusNoContent,
		},
		{
			id:   2,
			name: "Atoi Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodDelete,
				mockRequestURL:    "/user/dgsg",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				err: nil,
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Service Error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodDelete,
				mockRequestURL:    "/user/1",
				mockRequestBody:   strings.NewReader(``),
			},
			repoResp: mockRepoResp{
				err: errors.New("Эта ошибка ломает сервис"),
			},
			mockUsrID:      1,
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

			mockUserRepo := new(reposmocks.MockUserRepo)

			mockUserService := services.NewUserService(mockUserRepo)

			client := &http.Client{}

			userHandler := handlers.NewUserHandler(mockUserService, logger, client)

			mockUserRepo.On("DeleteUser", mock.AnythingOfType("*context.timerCtx"), tc.mockUsrID).Return(tc.repoResp.err)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/user/{user_id}", userHandler.DeleteUser).Methods("DELETE")

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockUserRepo.AssertCalled(t, "DeleteUser", mock.Anything, tc.mockUsrID)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type mockRepoResp struct {
		user models.User
		err  error
	}

	var updatedUser = models.User{
		ID:             1,
		PassportNumber: "1234 567890",
		Surname:        "Викторов",
		Name:           "Виктор",
		Patronymic:     "Викторович",
		Address:        "г.Санкт-Петербург",
	}

	testCases := []struct {
		id             int
		name           string
		mockReq        mockRequest
		repoResp       mockRepoResp
		userID         int
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPatch,
				mockRequestURL:    "/user/1",
				mockRequestBody: strings.NewReader(`{
				  "surname": "Викторов",
				  "name": "Виктор",
				  "patronymic": "Викторович",
				  "address": "г.Санкт-Петербург"
				}`),
			},
			repoResp: mockRepoResp{
				user: updatedUser,
				err:  nil,
			},
			userID:         1,
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   2,
			name: "Atoi error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPatch,
				mockRequestURL:    "/user/safsaf",
				mockRequestBody:   strings.NewReader(``),
			},
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   3,
			name: "Decode error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPatch,
				mockRequestURL:    "/user/1",
				mockRequestBody:   strings.NewReader(`{Это я сломал decode}`),
			},
			userID:         1,
			callRepo:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			id:   4,
			name: "Service error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPatch,
				mockRequestURL:    "/user/1",
				mockRequestBody: strings.NewReader(`{
				  "surname": "Викторов",
				  "name": "Виктор",
				  "patronymic": "Викторович",
				  "address": "г.Санкт-Петербург"
				}`),
			},
			repoResp: mockRepoResp{
				user: models.User{},
				err:  errors.New("Эта ошибка ломает сервис"),
			},
			userID:         1,
			callRepo:       true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			id:   5,
			name: "Encode error",
			mockReq: mockRequest{
				mockRequestMethod: http.MethodPatch,
				mockRequestURL:    "/user/1",
				mockRequestBody: strings.NewReader(`{
				  "surname": "Викторов",
				  "name": "Виктор",
				  "patronymic": "Викторович",
				  "address": "г.Санкт-Петербург"
				}`),
			},
			repoResp: mockRepoResp{
				user: updatedUser,
				err:  nil,
			},
			userID:         1,
			callRepo:       true,
			breakWrite:     true,
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

			mockUserRepo := new(reposmocks.MockUserRepo)

			mockUserService := services.NewUserService(mockUserRepo)

			client := &http.Client{}

			userHandler := handlers.NewUserHandler(mockUserService, logger, client)

			mockUserRepo.On("UpdateUser", mock.AnythingOfType("*context.timerCtx"), mockAPIUser, tc.userID).Return(tc.repoResp.user, tc.repoResp.err)

			req, err := http.NewRequest(tc.mockReq.mockRequestMethod, tc.mockReq.mockRequestURL, tc.mockReq.mockRequestBody)
			if err != nil {
				t.Fatal(err)
			}

			mockWriter := &errorResponseWriter{}

			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/user/{user_id}", userHandler.UpdateUser).Methods(http.MethodPatch)

			if tc.breakWrite {
				router.ServeHTTP(mockWriter, req)
				assert.Equal(t, tc.expectedStatus, mockWriter.Code)
			} else {
				router.ServeHTTP(rr, req)
				assert.Equal(t, tc.expectedStatus, rr.Code)
			}

			if tc.callRepo {
				mockUserRepo.AssertCalled(t, "UpdateUser", mock.Anything, mockAPIUser, tc.userID)
			}
		})
	}
}
