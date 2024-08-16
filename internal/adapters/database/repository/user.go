package repository

import (
	"context"
	"fmt"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	result := repo.db.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	fmt.Printf("user: %v\n", user)

	return &user, nil
}

func (repo UserRepository) FindUser(ctx context.Context, payload dto.FindUserDto) (*domain.User, error) {

	var user domain.User
	var result *gorm.DB

	if payload.Email != "" {
		result = repo.db.Db.Where("email = ?", payload.Email).First(&user)
	} else if payload.ID != "" {
		result = repo.db.Db.Where("id = ?", payload.ID).First(&user)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo UserRepository) UpdateUser(ctx context.Context, payload domain.User) error {
	result := repo.db.Db.Model(&payload).Where("id = ?", payload.Base.ID).Updates(payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo UserRepository) DeleteUser(ctx context.Context, payload dto.FindUserDto) error {
	result := repo.db.Db.Delete(payload)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo UserRepository) CreateToken(ctx context.Context, payload dto.ManageTokenDto) error {

	tokenData := domain.Token{
		UserId:      payload.ID,
		AccessToken: payload.AccessToken,
	}
	result := repo.db.Db.Create(tokenData)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo UserRepository) DeleteToken(ctx context.Context, payload dto.ManageTokenDto) error {

	tokenData := domain.Token{
		AccessToken: payload.AccessToken,
	}
	result := repo.db.Db.Where("access_token = ?", payload.AccessToken).Delete(tokenData)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
