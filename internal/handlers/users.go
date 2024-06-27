package handlers

import (
	"EMTask/internal/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

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
