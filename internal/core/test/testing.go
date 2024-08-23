package test

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/services"
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

func GetAuthenticatedTestUser(db *database.DB) (*domain.User, *domain.Token, error) {
	passwordService := services.NewPasswordService()
	tokenService := services.NewTokenService()
	userData := GenerateFakeUser()
	password, _ := passwordService.HashPassword("password")
	now := time.Now()
	payload := domain.User{
		Base:            userData.Base,
		FullName:        userData.FullName,
		Email:           userData.Email,
		EmailVerifiedAt: &now,
		Password: domain.Password{
			Value: password,
		},
	}
	result := db.Db.Create(&payload)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)
	tokenString, err := tokenService.CreateToken(payload.ID, expiry)
	if err != nil {
		return nil, nil, err
	}

	token := domain.Token{
		UserId:      payload.ID,
		AccessToken: tokenString,
	}
	result = db.Db.Create(&token)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &payload, &token, nil
}
