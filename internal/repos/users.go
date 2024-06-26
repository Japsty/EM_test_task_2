package repos

import (
	"EMTask/internal/models"
	"EMTask/internal/repos/queries"
	"context"
	"database/sql"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (ur *UsersRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	rows, err := ur.db.QueryContext(ctx, queries.GetAllUsers)
	if err != nil {
		return nil, err
	}

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

	}
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
func (ur *UsersRepository) FindUserByPassportHash(ctx context.Context, hash string) (models.User, error) {

}
func (ur *UsersRepository) FindUserByID(ctx context.Context, usrID int) (models.User, error) {

}
func (ur *UsersRepository) UpdateUser(ctx context.Context) (models.User, error) {

}
func (ur *UsersRepository) DeleteUser(ctx context.Context, usrID int) error {

}
