package reposmocks

import (
	"EMTask/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (repo *MockUserRepo) GetAllUsers(ctx context.Context, filter models.UserFilter, pg, lim int) ([]models.User, error) {
	args := repo.Called(ctx, filter, pg, lim)
	return args.Get(0).([]models.User), args.Error(1)
}

func (repo *MockUserRepo) AddUser(ctx context.Context, user models.ServiceUser) (int, error) {
	args := repo.Called(ctx, user)
	return args.Get(0).(int), args.Error(1)
}

func (repo *MockUserRepo) UpdateUser(ctx context.Context, newUser models.APIResponse, usrID int) (models.User, error) {
	args := repo.Called(ctx, newUser, usrID)
	return args.Get(0).(models.User), args.Error(1)
}

func (repo *MockUserRepo) DeleteUser(ctx context.Context, usrID int) error {
	args := repo.Called(ctx, usrID)
	return args.Error(0)
}
