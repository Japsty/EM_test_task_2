package repos

import (
	"EMTask/internal/models"
	"EMTask/internal/repos/queries"
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
)

var errNotFound = errors.New("user not found: ")

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (ur *UsersRepository) GetAllUsers(ctx context.Context, filter models.UserFilter, pg, lim int) ([]models.User, error) {
	query := squirrel.Select("id", "passport_number", "surname", "name", "patronymic", "address").
		From("users")

	if filter.PassportNum != "" {
		query = query.Where(squirrel.Eq{"passport_number": filter.PassportNum})
	}
	if filter.Surname != "" {
		query = query.Where(squirrel.Eq{"surname": filter.Surname})
	}
	if filter.Name != "" {
		query = query.Where(squirrel.Eq{"name": filter.Name})
	}
	if filter.Patronymic != "" {
		query = query.Where(squirrel.Eq{"patronymic": filter.Patronymic})
	}
	if filter.Address != "" {
		query = query.Where(squirrel.Eq{"address": filter.Address})
	}

	offset := (pg - 1) * lim
	query = query.Limit(uint64(lim)).Offset(uint64(offset))

	query = query.OrderBy("id")

	query = query.PlaceholderFormat(squirrel.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.PassportNumber,
			&user.Surname,
			&user.Name,
			&user.Patronymic,
			&user.Address,
		)
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
		user.PassportNum,
		user.Surname,
		user.Name,
		user.Patronymic,
		user.Address,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur *UsersRepository) UpdateUser(ctx context.Context, newUser models.APIResponse, usrID int) (models.User, error) {
	var user models.User

	err := ur.db.QueryRowContext(
		ctx,
		queries.UpdateUser,
		usrID,
		newUser.Surname,
		newUser.Name,
		newUser.Patronymic,
		newUser.Address,
	).Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
	if err != nil {
		return models.User{}, err
	}

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
