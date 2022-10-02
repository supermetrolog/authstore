package mysql_test

import (
	user "authstore/internal/domain/user/entity"
	"authstore/internal/domain/user/repository/mysql"
	"authstore/pkg/logging"
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	return db, mock
}

func TestCreate(t *testing.T) {
	var (
		email    = "test@mail.ru"
		username = "admin"
		password = "password"
	)

	var createUserDTO = &user.CreateUserDTO{
		Email:    &email,
		Username: &username,
		Password: &password,
	}
	db, mock := NewMock()
	repo := mysql.NewRepository(logging.GetLogger(), db)
	defer repo.Close()

	query := "INSERT INTO user \\(email, username, password_hash\\) VALUES \\(\\?, \\?, \\?\\)"
	//Ожидается такой запрос
	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(createUserDTO.Email, createUserDTO.Username, createUserDTO.Password).WillReturnResult(sqlmock.NewResult(0, 1))
	lastInsertId, err := repo.Create(context.Background(), createUserDTO)

	t.Log(lastInsertId, err)
	assert.Equal(t, user.UserID(0), lastInsertId)
	assert.NoError(t, err)
}

func TestFindByID(t *testing.T) {
	var (
		id           user.UserID = 45234
		email                    = "fuck@gmail.com"
		username                 = "sukablyat"
		passwordHash             = "sdmaksmdwnaldmnawda"
	)
	var findID user.UserID = 45234
	var model = &user.User{
		ID:           &id,
		Email:        &email,
		Username:     &username,
		PasswordHash: &passwordHash,
	}

	db, mock := NewMock()
	repo := mysql.NewRepository(logging.GetLogger(), db)
	defer repo.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"email",
		"username",
		"password_hash",
		"refresh_token",
	}).AddRow(
		model.ID,
		model.Email,
		model.Username,
		model.PasswordHash,
		model.RefreshToken,
	)
	query := "SELECT id, username, email, password_hash, refresh_token FROM user WHERE id = \\? LIMIT 1"
	//Ожидается такой запрос
	mock.ExpectQuery(query).WithArgs(findID).WillReturnRows(rows)

	u, err := repo.FindById(context.Background(), findID)
	assert.NotNil(t, u)
	assert.Equal(t, findID, *u.ID)

	assert.NoError(t, err)
}
