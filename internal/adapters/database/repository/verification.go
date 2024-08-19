package repository

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *database.DB
}

func NewVerificationRepository(db *database.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r VerificationRepository) CreateCode(ctx context.Context, payload dto.VerificationDto) error {
	current, _ := r.FindCode(ctx, payload)

	if current != nil && current.ID != "" {
		_ = r.DeleteCode(ctx, current.ID)
	}

	data := domain.Verifications{
		UserId: payload.UserID,
		Code:   payload.Code,
		Type:   payload.Type,
	}

	result := r.db.Db.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r VerificationRepository) FindCode(ctx context.Context, payload dto.VerificationDto) (*domain.Verifications, error) {
	var verification domain.Verifications
	var result *gorm.DB

	if payload.UserID != "" {
		result = r.db.Db.Where("user_id = ?", payload.UserID).Where("type = ?", payload.Type).First(&verification)
	} else {
		result = r.db.Db.Where("code = ?", payload.Code).Where("type = ?", payload.Type).First(&verification)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &verification, nil
}

func (r VerificationRepository) DeleteCode(ctx context.Context, id string) error {
	result := r.db.Db.Delete(&domain.Verifications{Base: domain.Base{ID: id}})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
