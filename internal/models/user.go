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

type ServiceUser struct {
	PassportHash string `json:"passportHash"`
	Surname      string `json:"surname"`
	Name         string `json:"name"`
	Patronymic   string `json:"patronymic"`
	Address      string `json:"address"`
}

type APIResponse struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

type UserRepo interface {
	GetAllUsers(context.Context) ([]User, error)
	AddUser(context.Context, ServiceUser) (int, error)
	FindUserByID(context.Context, int) (User, error)
	UpdateUser(context.Context, ServiceUser, int) (User, error)
	DeleteUser(context.Context, int) error
}

type UserService interface {
	GetAllUsers(context.Context) ([]User, error)
	CreateUser(context.Context, APIResponse, string) (User, error)
	GetUserByID(context.Context, int) (User, error)
	UpdateUser(context.Context, APIResponse, int) (User, error)
	DeleteUser(context.Context, int) error
}
