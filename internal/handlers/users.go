package handlers

import (
	"EMTask/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var TimeoutTime = 500 * time.Millisecond
var passportNumberPattern = `^\d{4} \d{6}$`

type UserHandler struct {
	UserService models.UserService
	ZapLogger   *zap.SugaredLogger
	Client      *http.Client
}

func NewUserHandler(us models.UserService, logger *zap.SugaredLogger, client *http.Client) *UserHandler {
	return &UserHandler{us, logger, client}
}

func (uh *UserHandler) getPeopleInfo(passportNumber string) (models.APIResponse, error) {
	apiURL := fmt.Sprintf(
		"%s/info?passportSerie=%s&passportNumber=%s",
		os.Getenv("API_URL"),
		passportNumber[:4],
		passportNumber[5:],
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, nil)
	if err != nil {
		return models.APIResponse{}, err
	}

	resp, err := uh.Client.Do(req)
	if err != nil {
		return models.APIResponse{}, err
	}

	defer resp.Body.Close()

	var apiResponse models.APIResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return models.APIResponse{}, err
	}

	return apiResponse, nil
}

// @Summary Get Users
// @Description Получить юзеров с пагинацией и фильтрацией
// @Tags users
// @Produce json
// @Param passport query string false "1234 567890"
// @Param surname query string false "Иванов"
// @Param name query string false "Иван"
// @Param patronymic query string false "Иванович"
// @Param address query string false "г. Москва, ул. Ленина, д. 5, кв. 1"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Limit per page"
// @Success 200 {array} models.User
// @Failure 400 {string} string "Invalid Page or Limit param"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
func (uh *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	queryParams := r.URL.Query()
	filter := models.UserFilter{
		PassportNum: queryParams.Get("passport"),
		Surname:     queryParams.Get("surname"),
		Name:        queryParams.Get("name"),
		Patronymic:  queryParams.Get("patronymic"),
		Address:     queryParams.Get("address"),
	}

	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil || page < 1 {
		uh.ZapLogger.Infof(reqIDString+"GetUsers Invalid Page param: ", r.URL.Query())
		http.Error(w, " Invalid Page param", http.StatusBadRequest)

		return
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit < 1 {
		uh.ZapLogger.Infof(reqIDString+"GetUsers Invalid Limit param: ", r.URL.Query())
		http.Error(w, "Invalid Limit param", http.StatusBadRequest)

		return
	}

	users, err := uh.UserService.GetAllUsers(ctxWthTimeout, filter, page, limit)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"GetUsers GetAllUsers Error: ", err)
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

// @Summary Delete User by ID
// @Description Удалить юзера по ID
// @Tags users
// @Param user_id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid user_id"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{user_id} [delete]
func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		uh.ZapLogger.Infof(reqIDString+"DeleteUser Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid user_id", http.StatusBadRequest)

		return
	}

	err = uh.UserService.DeleteUser(ctxWthTimeout, userID)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"DeleteUser Service Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update User by ID
// @Description Обновить юзера по ID
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /user/{user_id} [patch]
func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	userID, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		uh.ZapLogger.Infof(reqIDString+"UpdateUser Atoi Error: ", r.URL.Query())
		http.Error(w, "Invalid user_id", http.StatusBadRequest)

		return
	}

	var user models.APIResponse

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	updatedUser, err := uh.UserService.UpdateUser(ctxWthTimeout, user, userID)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"UpdateUser Service Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"GetUsers Encode Error: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}
}

// @Summary Add a new user
// @Description Добавить пользователя по его паспортным данным
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.NewUserRequest true "New User"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /user [post]
func (uh *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctxWthTimeout, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	reqIDString := fmt.Sprintf("requestID: %s ", r.Context().Value("requestID"))

	var usersPassportData models.NewUserRequest

	err := json.NewDecoder(r.Body).Decode(&usersPassportData)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser Decode Error, caused by: ", r.Body)
		http.Error(w, "Invalid input", http.StatusBadRequest)

		return
	}

	match, _ := regexp.MatchString(passportNumberPattern, usersPassportData.PassportNumber)
	if !match {
		uh.ZapLogger.Infof(reqIDString + "AddUser Invalid Passport Number Format")
		http.Error(w, "Invalid passport number format. Expected format: '1234 567890'", http.StatusBadRequest)

		return
	}

	apiResponse, err := uh.getPeopleInfo(usersPassportData.PassportNumber)
	if err != nil {
		uh.ZapLogger.Error(reqIDString+"AddUser getPeopleInfo Error: ", err)
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
