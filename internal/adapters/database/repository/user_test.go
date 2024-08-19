package repository

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stretchr/testify/assert"
)

func GenerateFakeUser() domain.User {
	return domain.User{
		Base: domain.Base{
			ID:        gofakeit.UUID(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		},
		FullName: gofakeit.Name(),
		Email:    gofakeit.Email(),
	}
}

func TestCreateUser(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")

	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, err := repo.CreateUser(ctx, user)

	assert.NoError(t, err, "Expected no error on user creation")
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")
	assert.Equal(t, user.FullName, createdUser.FullName, "Expected full name to match")
	assert.Equal(t, user.Email, createdUser.Email, "Expected email to match")
}

func TestCreateUserWithExistingEmail(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	otherUser, err := repo.CreateUser(ctx, user)
	assert.Nil(t, otherUser)
	assert.Error(t, err, "UNIQUE constraint failed: users.email")
}

func TestFindUserById(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{ID: createdUser.ID})

	assert.NotNil(t, foundUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, foundUser.ID, "Expected a non-empty user ID")
	assert.Equal(t, foundUser.ID, createdUser.ID)
}

func TestFindUserByEmail(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{Email: user.Email})

	assert.NotNil(t, foundUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, foundUser.ID, "Expected a non-empty user ID")
	assert.Equal(t, foundUser.ID, createdUser.ID)
}

func TestFindUserReturnsNilForNonExistingUsers(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{ID: gofakeit.UUID()})

	assert.Nil(t, foundUser)
	assert.NotNil(t, err)
}

func TestUpdateUser(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	updates := map[string]interface{}{
		"full_name": "JANE USED FOR TESTING",
	}

	err = repo.UpdateUser(ctx, createdUser.ID, updates)
	assert.Nil(t, err)

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{ID: createdUser.ID})
	assert.NotNil(t, foundUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, foundUser.ID, "Expected a non-empty user ID")
	assert.Equal(t, foundUser.FullName, updates["full_name"])
}

func TestUpdateUserWithInvalidFields(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	updates := map[string]interface{}{
		"random_field": "JANE USED FOR TESTING",
	}

	err = repo.UpdateUser(ctx, createdUser.ID, updates)
	assert.Error(t, err, "no such column: random_field")
}

func TestUpdateUserWithEmptyRecords(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	updates := map[string]interface{}{}

	err = repo.UpdateUser(ctx, createdUser.ID, updates)
	assert.Nil(t, err)

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{ID: createdUser.ID})
	assert.NotNil(t, foundUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, foundUser.ID, "Expected a non-empty user ID")
}

func TestUpdateNonExistingRecord(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()

	updates := map[string]interface{}{
		"full_name": "JANE USED FOR TESTING",
	}

	id := gofakeit.UUID()
	err = repo.UpdateUser(ctx, id, updates)
	assert.Error(t, err, "record not found")
}

func TestDeleteExistingUser(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	err = repo.DeleteUser(ctx, createdUser.ID)
	assert.Nil(t, err)

	foundUser, err := repo.FindUser(ctx, dto.FindUserDto{ID: createdUser.ID})
	assert.Nil(t, foundUser)
	assert.Error(t, err, "record not found")
}

func TestDeleteAlreadyDeletedUser(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	err = repo.DeleteUser(ctx, createdUser.ID)
	assert.Nil(t, err)

	err = repo.DeleteUser(ctx, createdUser.ID)
	assert.Error(t, err, "record not found")
}

func TestDeleteNonExistingUser(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	id := gofakeit.UUID()

	err = repo.DeleteUser(ctx, id)
	assert.Error(t, err, "record not found")
}

func TestCreateAccessToken(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()
	tokenString := gofakeit.Name()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	err = repo.CreateToken(ctx, dto.ManageTokenDto{
		ID:          createdUser.ID,
		AccessToken: tokenString,
	})
	assert.Nil(t, err)
}

func TestDeleteToken(t *testing.T) {
	db, err := database.SetupMockDB()
	assert.NoError(t, err, "Failed to setup test database")
	repo := NewUserRepository(db)
	ctx := context.Background()
	user := GenerateFakeUser()
	tokenString := gofakeit.Name()

	createdUser, _ := repo.CreateUser(ctx, user)
	assert.NotNil(t, createdUser, "Expected a valid created user")
	assert.NotEmpty(t, createdUser.ID, "Expected a non-empty user ID")

	err = repo.CreateToken(ctx, dto.ManageTokenDto{
		ID:          createdUser.ID,
		AccessToken: tokenString,
	})
	assert.Nil(t, err)

	err = repo.DeleteToken(ctx, dto.ManageTokenDto{AccessToken: tokenString})
	assert.Nil(t, err)
}
