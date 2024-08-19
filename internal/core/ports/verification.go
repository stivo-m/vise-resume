package ports

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

type VerificationPort interface {
	CreateCode(ctx context.Context, user dto.VerificationDto) error
	FindCode(ctx context.Context, payload dto.VerificationDto) (*domain.Verifications, error)
	DeleteCode(ctx context.Context, id string) error
}

type VerificationService interface {
	VerifyCode(ctx context.Context, payload dto.VerificationDto) error
	GenerateCode(ctx context.Context, payload dto.VerificationDto) error
}
