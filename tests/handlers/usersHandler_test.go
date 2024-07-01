package handlers

import (
	"EMTask/internal/handlers"
	"EMTask/internal/models"
	"EMTask/internal/services"
	"EMTask/tests/mocks/reposmocks"
	"bytes"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorResponseWriter struct {
	httptest.ResponseRecorder
}

type mockRequest struct {
	mockRequestMethod        string
	mockRequestURL           string
	mockRequestBody          *strings.Reader
	mockIncorrectRequestBody *bytes.Buffer
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
		callRepo       bool
		breakWrite     bool
		expectedStatus int
	}{
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: "POST",
				mockRequestURL:    "/api/posts/",
				mockRequestBody:   strings.NewReader(`{"category":"music", "text":"MockPost","title":"MockPost","postType":"text"}`),
			},
			repoResp: mockRepoResp{
				mockPost:  mockPost,
				mockError: nil,
			},
			callRepo:       true,
			expectedStatus: http.StatusOK,
		},
		{
			id:   1,
			name: "Success",
			mockReq: mockRequest{
				mockRequestMethod: "POST",
				mockRequestURL:    "/api/posts/",
				mockRequestBody:   strings.NewReader(`{"category":"music", "text":"MockPost","title":"MockPost","postType":"text"}`),
			},
			repoResp: mockRepoResp{
				mockPost:  mockPost,
				mockError: nil,
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

			mockUserRepo := new(reposmocks.MockUserRepo)

			mockUserService := services.NewUserService(mockUserRepo)

			client := &http.Client{}

			userHandler := handlers.NewUserHandler(mockUserService, logger, client)

			mockUser := models.ServiceUser{}

			mockUserRepo.On("AddUser", mock.AnythingOfType("*context.timerCtx")).Return(tc.repoResp.usrID, tc.repoResp.mockError)

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
				mockUserRepo.AssertCalled(t, "GetAllPosts", mock.Anything)
			}
		})
	}
}
