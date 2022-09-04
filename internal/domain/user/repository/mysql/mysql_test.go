package mysql_test

import (
	"authstore/internal/domain/user/entity/user"
	"authstore/internal/domain/user/repository/mysql"
	"authstore/pkg/logging"
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

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

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	return db, mock
}

func TestCreate(t *testing.T) {
	db, mock := NewMock()
	repo := mysql.NewRepository(logging.GetLogger(), db)
	defer repo.Close()

	query := "INSERT INTO user \\(email, username, password_hash\\) VALUES \\(\\?, \\?, \\?\\)"
	//Ожидается такой запрос
	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(createUserDTO.Email, createUserDTO.Username, createUserDTO.Password).WillReturnResult(sqlmock.NewResult(0, 1))
	lastInsertId, err := repo.Create(context.Background(), createUserDTO)

	assert.Equal(t, user.UserID(0), lastInsertId)
	assert.NoError(t, err)
}
