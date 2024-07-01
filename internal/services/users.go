package services

import (
	"EMTask/internal/models"
	"context"
)

type UsersService struct {
	usersRepo models.UserRepo
	apiURL    string
}

func NewUserService(repo models.UserRepo) *UsersService {
	return &UsersService{usersRepo: repo}
}

func (us *UsersService) GetAllUsers(ctx context.Context, filter models.UserFilter, pg, lim int) ([]models.User, error) {
	users, err := us.usersRepo.GetAllUsers(ctx, filter, pg, lim)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UsersService) CreateUser(ctx context.Context, resp models.APIResponse, passport string) (models.User, error) {
	user := models.ServiceUser{
		PassportNum: passport,
		Surname:     resp.Surname,
		Name:        resp.Name,
		Patronymic:  resp.Patronymic,
		Address:     resp.Address,
	}

	userID, err := us.usersRepo.AddUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:             userID,
		PassportNumber: passport,
		Surname:        resp.Surname,
		Name:           resp.Name,
		Patronymic:     resp.Patronymic,
		Address:        resp.Address,
	}, nil
}

func (us *UsersService) GetUserByID(ctx context.Context, usrID int) (models.User, error) {
	users, err := us.usersRepo.FindUserByID(ctx, usrID)
	if err != nil {
		return models.User{}, err
	}

	return users, nil
}
func (us *UsersService) UpdateUser(ctx context.Context, response models.APIResponse, usrID int) (models.User, error) {
	users, err := us.usersRepo.UpdateUser(ctx, response, usrID)
	if err != nil {
		return models.User{}, err
	}

	return users, nil
}
func (us *UsersService) DeleteUser(ctx context.Context, usrID int) error {
	err := us.usersRepo.DeleteUser(ctx, usrID)
	if err != nil {
		return err
	}

	return nil
}
