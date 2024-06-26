package services

import (
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"context"
)

type UsersService struct {
	usersRepo repos.UsersRepository
}

func NewUserService(repo repos.UsersRepository) *UsersService {
	return &UsersService{usersRepo: repo}
}

func (us *UsersService) GetAllUsers(ctx context.Context) ([]models.User, error) {

}

func (us *UsersService) CreateUser(context.Context, string) (models.User, error) {

}
func (us *UsersService) GetUserByPassportNumber(context.Context, string) (models.User, error) {

}
func (us *UsersService) GetUserByID(context.Context, int) (models.User, error) {

}
func (us *UsersService) UpdateUser(ctx context.Context) (models.User, error) {

}
func (us *UsersService) DeleteUser(context.Context, int) error {

}
