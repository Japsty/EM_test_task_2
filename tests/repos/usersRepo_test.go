package repos

import (
	"EMTask/internal/repos"
	"EMTask/internal/repos/queries"
	"EMTask/internal/services"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"regexp"
	"testing"
)

func TestAddUser(t *testing.T) {
	testCases := []struct {
		TestCaseID    int
		Name          string
		ExpectedError error
	}{
		{
			TestCaseID:    1,
			Name:          "Success",
			ExpectedError: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected ", err)
			}
			defer db.Close()

			es := services.NewEncodeService("secret")
			repo := repos.NewUsersRepository(db, es)

			mock.ExpectQuery(regexp.QuoteMeta(queries.CreateUser)).WithArgs(
				tc.InputComment.PostID).WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "content", "user_id", "comments_allowed", "created_at"}).
					AddRow(
						tc.InputComment.ID,
						tc.InputComment.Content,
						tc.InputComment.AuthorID,
						tc.InputComment.PostID,
						tc.InputComment.ParentID,
						createdAtTime,
					).AddRow(
					mockChildComment.ID,
					mockChildComment.Content,
					mockChildComment.AuthorID,
					mockChildComment.PostID,
					mockChildComment.ParentID,
					createdAtTime,
				).AddRow(
					mockChildChildComment.ID,
					mockChildChildComment.Content,
					mockChildChildComment.AuthorID,
					mockChildChildComment.PostID,
					mockChildChildComment.ParentID,
					createdAtTime,
				),
			)

			comments, err := repo.AddUser(context.Background())
			if err != nil {
				t.Fatalf("TestGetCommentByPostIDPaginated Error: %s", err)
			}

			for idx, comm := range comments {
			}

			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
