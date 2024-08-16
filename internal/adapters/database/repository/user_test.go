package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := database.SetupMockDB()
	assert.NoError(t, err)

	repo := NewUserRepository(db)

	testName := "test name"
	testEmail := "test@example.com"
	testId := uuid.New().String()

	user := domain.User{
		Base: domain.Base{
			ID: testId,
		},
		FullName: testName,
		Email:    testEmail,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"users\"").WithArgs(user.Base.ID, user.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the method
	_, err = repo.CreateUser(context.Background(), user)

	// Validate the results
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

}
