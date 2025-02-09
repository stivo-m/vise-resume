package repository

import (
	"context"

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

	return &user, nil
}

func (repo UserRepository) FindUser(ctx context.Context, payload dto.FindUserDto) (*domain.User, error) {

	var user domain.User
	var result *gorm.DB

	if payload.Email != "" {
		result = repo.db.Db.Preload("Password").Where("email = ?", payload.Email).First(&user)
	} else if payload.ID != "" {
		result = repo.db.Db.Preload("Password").Where("id = ?", payload.ID).First(&user)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo UserRepository) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	result := repo.db.Db.Model(&domain.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo UserRepository) UpdateUserPassword(ctx context.Context, id string, password domain.Password) error {
	user, err := repo.FindUser(ctx, dto.FindUserDto{
		ID: id,
	})

	if err != nil {
		return err
	}

	_ = repo.db.Db.Delete(&domain.Password{Base: domain.Base{ID: user.Password.ID}})
	user.Password = password
	if err := repo.db.Db.Save(&user).Error; err != nil {
		return err
	}

	return nil

}

func (repo UserRepository) DeleteUser(ctx context.Context, id string) error {
	result := repo.db.Db.Delete(&domain.User{Base: domain.Base{ID: id}})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repo UserRepository) FindToken(ctx context.Context, payload dto.ManageTokenDto) (*domain.Token, error) {
	var token *domain.Token
	result := repo.db.Db.Where("user_id = ?", payload.ID).Where("access_token = ?", payload.AccessToken).First(&token)

	if result.Error != nil {
		return nil, result.Error
	}

	return token, nil
}

func (repo UserRepository) CreateToken(ctx context.Context, payload dto.ManageTokenDto) error {

	tokenData := domain.Token{
		UserId:      payload.ID,
		AccessToken: payload.AccessToken,
	}
	result := repo.db.Db.Create(&tokenData)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo UserRepository) DeleteToken(ctx context.Context, payload dto.ManageTokenDto) error {

	tokenData := domain.Token{
		AccessToken: payload.AccessToken,
	}
	result := repo.db.Db.Where("access_token = ?", payload.AccessToken).Delete(&tokenData)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
