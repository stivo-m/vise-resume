package repository

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

type VerificationRepository struct {
	db *database.DB
}

func NewVerificationRepository(db *database.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r VerificationRepository) CreateCode(ctx context.Context, user dto.VerificationDto) error {
	return nil
}

func (r VerificationRepository) FindCode(ctx context.Context, payload dto.VerificationDto) (*domain.Verifications, error) {
	return nil, nil
}

func (r VerificationRepository) DeleteCode(ctx context.Context, payload dto.VerificationDto) error {
	return nil
}
