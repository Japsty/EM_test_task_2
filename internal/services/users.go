package services

import (
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"context"
)

type UsersService struct {
	usersRepo  *repos.UsersRepository
	encService EncodeService
	apiURL     string
}

func NewUserService(repo *repos.UsersRepository) *UsersService {
	return &UsersService{usersRepo: repo}
}

func (us *UsersService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := us.usersRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UsersService) CreateUser(ctx context.Context, resp models.APIResponse, passport string) (models.User, error) {
	passportEncoded, err := us.encService.Encrypt(passport)
	if err != nil {
		return models.User{}, err
	}

	user := models.ServiceUser{
		PassportHash: passportEncoded,
		Surname:      resp.Surname,
		Name:         resp.Name,
		Patronymic:   resp.Patronymic,
		Address:      resp.Address,
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
	serviceUser := models.ServiceUser{}

	users, err := us.usersRepo.UpdateUser(ctx, serviceUser, usrID)
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
