package repos

import (
	"EMTask/internal/models"
	"EMTask/internal/repos/queries"
	"EMTask/internal/services"
	"context"
	"database/sql"
	"errors"
)

var errNotFound = errors.New("user not found: ")

type UsersRepository struct {
	db         *sql.DB
	encService services.EncodeService
}

func NewUsersRepository(db *sql.DB, es *services.EncodeService) *UsersRepository {
	return &UsersRepository{db: db, encService: *es}
}

func (ur *UsersRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := ur.db.QueryContext(ctx, queries.GetAllUsers)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		var passHash string

		err := rows.Scan(
			user.ID,
			passHash,
			user.Surname,
			user.Name,
			user.Patronymic,
			user.Address,
		)
		if err != nil {
			return nil, err
		}
		user.PassportNumber, err = ur.encService.Decrypt(passHash)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UsersRepository) AddUser(ctx context.Context, user models.ServiceUser) (int, error) {
	var userID int
	err := ur.db.QueryRowContext(
		ctx,
		queries.CreateUser,
		user.PassportHash,
		user.Surname,
		user.Name,
		user.Patronymic,
		user.Address,
	).Scan(userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
func (ur *UsersRepository) FindUserByID(ctx context.Context, usrID int) (models.User, error) {
	var user models.User
	var passHash string

	err := ur.db.QueryRowContext(ctx, queries.FindUserByID, usrID).Scan(
		user.ID,
		passHash,
		user.Surname,
		user.Name,
		user.Patronymic,
		user.Address,
	)
	if err != nil {
		return models.User{}, err
	}

	passportData, err := ur.encService.Decrypt(passHash)
	if err != nil {
		return models.User{}, err
	}

	user.PassportNumber = passportData
	return user, nil
}
func (ur *UsersRepository) UpdateUser(ctx context.Context, newUser models.ServiceUser, usrID int) (models.User, error) {
	var user models.User
	var passHash string

	err := ur.db.QueryRowContext(
		ctx,
		queries.UpdateUser,
		usrID,
		newUser.PassportHash,
		newUser.Surname,
		newUser.Name,
		newUser.Patronymic,
		newUser.Address,
	).Scan(user.ID, passHash, user.Surname, user.Name, user.Patronymic, user.Address)
	if err != nil {
		return models.User{}, err
	}

	passportData, err := ur.encService.Decrypt(passHash)
	if err != nil {
		return models.User{}, err
	}

	user.PassportNumber = passportData

	return user, nil
}
func (ur *UsersRepository) DeleteUser(ctx context.Context, usrID int) error {
	result, err := ur.db.ExecContext(ctx, queries.DeleteUser, usrID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}
