package models

import "context"

type User struct {
	ID             int    `json:"id"`
	PassportNumber string `json:"passportNumber"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type DBUser struct {
	ID           int    `json:"id"`
	PassportHash string `json:"passportHash"`
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	Patronymic   string `json:"patronymic"`
	Address      string `json:"address"`
}

type ServiceUser struct {
	PassportHash string `json:"passportHash"`
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	Patronymic   string `json:"patronymic"`
	Address      string `json:"address"`
}

type UserRepo interface {
	GetAllUsers(context.Context) ([]User, error)
	AddUser(context.Context, ServiceUser) (int, error)
	FindUserByPassportHash(context.Context, string) (User, error)
	FindUserByID(context.Context, int) (User, error)
	UpdateUser(ctx context.Context) (User, error)
	DeleteUser(context.Context, int) error
}

type UserService interface {
	GetAllUsers(context.Context) ([]User, error)
	CreateUser(context.Context, string) (User, error)
	GetUserByPassportNumber(context.Context, string) (User, error)
	GetUserByID(context.Context, int) (User, error)
	UpdateUser(ctx context.Context) (User, error)
	DeleteUser(context.Context, int) error
}
