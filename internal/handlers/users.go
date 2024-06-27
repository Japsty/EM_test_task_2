package handlers

import (
	"EMTask/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var TimeoutTime = 500 * time.Millisecond

type UserHandler struct {
	UserService models.UserService
	ZapLogger   *zap.SugaredLogger
}

func NewUserHandler(us models.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{us, logger}
}

func getPeopleInfo(apiURL, passportNumber string) (models.APIResponse, error) {
	url := fmt.Sprintf("%s/info?passportSerie=%s&passportNumber=%s", apiURL, passportNumber[:4], passportNumber[5:])
	resp, err := http.Get(url)
	if err != nil {
		return models.APIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.APIResponse{}, fmt.Errorf("error fetching people info: %s", body)
	}

	var apiResp models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return models.APIResponse{}, err
	}

	return apiResp, nil
}

func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	users, err := uh.UserService.GetAllUsers(ctxWthTimeout)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"GetUsers Atoi Error, caused by: ", r.Body)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"GetUsers Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])
	err := uh.UserService.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	userID, _ := strconv.Atoi(mux.Vars(r)["user_id"])

	var user models.APIResponse
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedUser, err := uh.UserService.UpdateUser(ctxWthTimeout, user, userID)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"GetUsers Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	var usersPassportData models.NewUserRequest

	err := json.NewDecoder(r.Body).Decode(&usersPassportData)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Decode Error, caused by: ", r.Body)
		http.Error(w, "Invalid input", http.StatusInternalServerError)
		return
	}

	passportData := strings.Split(usersPassportData.PassportNumber, " ")

	series, err := strconv.Atoi(passportData[0])
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Atoi Error, caused by: ", r.Body)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	number, err := strconv.Atoi(passportData[1])
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Atoi Error, caused by: ", r.Body)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if series/1000 == 0 && number/100000 == 0 {

	}

	apiURL := fmt.Sprintf(
		"%s/info?passportSerie=%s&passportNumber=%s",
		os.Getenv("API_URL"),
		usersPassportData.PassportNumber[:5],
		usersPassportData.PassportNumber[5:],
	)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Request Error, caused by: ", r.Body)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var apiResponse models.APIResponse

	err = json.NewDecoder(req.Body).Decode(&apiResponse)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Decode Error, caused by: ", r.Body)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user, err := uh.UserService.CreateUser(ctxWthTimeout, apiResponse, usersPassportData.PassportNumber)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Service Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
