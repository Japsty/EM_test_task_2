package repos_test

import (
	"EMTask/internal/models"
	"EMTask/internal/repos"
	"EMTask/internal/repos/queries"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
)

var mockUser = models.User{
	ID:             1,
	PassportNumber: "1234 567890",
	Surname:        "Иванов",
	Name:           "Иван",
	Patronymic:     "Иванович",
	Address:        "г. Москва, ул. Ленина, д. 5, кв. 1",
}

var mockServiceUser = models.ServiceUser{
	PassportNum: "1234 567890",
	Surname:     "Иванов",
	Name:        "Иван",
	Patronymic:  "Иванович",
	Address:     "г. Москва, ул. Ленина, д. 5, кв. 1",
}

var mockAPIUser = models.APIResponse{
	Surname:    "Иванов",
	Name:       "Иван",
	Patronymic: "Иванович",
	Address:    "г. Москва, ул. Ленина, д. 5, кв. 1",
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewUsersRepository(db)

	mock.ExpectQuery("SELECT id, passport_number, surname, name, patronymic, address FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "passport_number", "surname", "name", "patronymic", "address"}).
			AddRow(1, "1234 567890", "Иванов", "Иван", "Иванович", "г. Москва, ул. Ленина, д. 5, кв. 1").
			AddRow(2, "2234 567890", "Иванов", "Виктор", "Иванович", "г. Москва, ул. Ленина, д. 5, кв. 1"))

	users, err := repo.GetAllUsers(context.Background(), models.UserFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("error fetching users: %s", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if users[0].ID != 1 || users[1].ID != 2 {
		t.Errorf("unexpected user IDs: got %d and %d, want 1 and 2", users[0].ID, users[1].ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewUsersRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(queries.CreateUser)).
		WithArgs(
			mockUser.PassportNumber,
			mockUser.Surname,
			mockUser.Name,
			mockUser.Patronymic,
			mockUser.Address,
		).
		WillReturnRows(sqlmock.NewRows([]string{"userID"}).AddRow(1))

	userID, err := repo.AddUser(
		context.Background(),
		mockServiceUser,
	)
	if err != nil {
		t.Fatalf("AddUser Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if userID != mockUser.ID {
		t.Errorf("unexpected ID: got %v, want %v", userID, mockUser.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewUsersRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(queries.UpdateUser)).
		WithArgs(
			mockUser.ID,
			mockUser.Surname,
			mockUser.Name,
			mockUser.Patronymic,
			mockUser.Address,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"passport_number",
				"surname",
				"name",
				"patronymic",
				"address",
			}).AddRow(
				mockUser.ID,
				mockUser.PassportNumber,
				mockUser.Surname,
				mockUser.Name,
				mockUser.Patronymic,
				mockUser.Address,
			))

	user, err := repo.UpdateUser(
		context.Background(),
		mockAPIUser,
		mockUser.ID,
	)
	if err != nil {
		t.Fatalf("AddUser Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if user.ID != mockUser.ID {
		t.Errorf("unexpected ID: got %v, want %v", user.ID, mockUser.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	defer db.Close()

	repo := repos.NewUsersRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(queries.DeleteUser)).
		WithArgs(
			mockUser.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteUser(
		context.Background(),
		mockUser.ID,
	)
	if err != nil {
		t.Fatalf("AddUser Error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
