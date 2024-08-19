package ports

import (
	"context"

	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

type UserPort interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	FindUser(ctx context.Context, payload dto.FindUserDto) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, id string) error
	CreateToken(ctx context.Context, payload dto.ManageTokenDto) error
	DeleteToken(ctx context.Context, payload dto.ManageTokenDto) error
}

type UserService interface {
	RegisterUser(ctx context.Context, payload dto.RegisterDto) (*dto.ProfileResponse, error)
	LoginUser(ctx context.Context, payload dto.LoginDto) (*dto.LoginResponse, error)
	VerifyEmailAddress(ctx context.Context, payload dto.VerificationDto) error
	UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error
	ForgetPassword(ctx context.Context, payload dto.EmailDto) error
	ResetPassword(ctx context.Context, payload dto.ResetPasswordDto) error
	LogoutUser(ctx context.Context, token string) error
}
