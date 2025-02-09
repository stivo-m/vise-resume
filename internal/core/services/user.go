package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/dto"
	"github.com/stivo-m/vise-resume/internal/core/ports"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

type UserService struct {
	userPort         ports.UserPort
	tokenService     ports.TokenService
	passwordService  ports.PasswordService
	verificationPort ports.VerificationPort
}

func NewUserService(
	userPort ports.UserPort,
	tokenService ports.TokenService,
	passwordService ports.PasswordService,
	verificationPort ports.VerificationPort,

) *UserService {
	return &UserService{
		userPort:         userPort,
		tokenService:     tokenService,
		passwordService:  passwordService,
		verificationPort: verificationPort,
	}
}

// The [RegisterUser] usecase is primarily to have a user created for the system
func (s UserService) RegisterUser(ctx context.Context, payload dto.RegisterDto) (*dto.ProfileResponse, error) {

	password, err := s.passwordService.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	userData := domain.User{
		FullName: payload.FullName,
		Email:    payload.Email,
		Password: domain.Password{Value: password},
	}

	user, err := s.userPort.CreateUser(ctx, userData)
	if err != nil {
		return nil, err
	}

	// TODO: Send verification email to the user
	code := utils.EncodeToString(6)
	err = s.verificationPort.CreateCode(ctx, dto.VerificationDto{
		UserID: user.ID,
		Code:   code,
		Type:   "email-verification",
	})

	if err != nil {
		return nil, err
	}

	response := dto.ProfileResponse{
		User: dto.UserResponseDto{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
		},
	}

	return &response, nil
}

// The [LoginUser] usecase is primarily used to ensure a user can be authenticated by the system
// This function should verify credentials used and provide the user with a token for their identity.
func (s UserService) LoginUser(ctx context.Context, payload dto.LoginDto) (*dto.LoginResponse, error) {
	user, err := s.userPort.FindUser(ctx, dto.FindUserDto{
		Email:        payload.Email,
		WithPassword: true,
	})

	if err != nil {
		return nil, err
	}

	if user.ID == "" {
		utils.TextLogger.Error("User not found", fmt.Errorf("user: %v", user))
		return nil, errors.New("either user was not found or password is incorrect")
	}

	if user.EmailVerifiedAt == nil {
		return nil, errors.New("email address is not verified")
	}

	match := s.passwordService.VerifyPassword(payload.Password, user.Password.Value)
	if !match {
		utils.TextLogger.Error("password mismatch", fmt.Errorf("user: %v", user))
		return nil, errors.New("either user was not found or password is incorrect")
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)
	token, err := s.tokenService.CreateToken(user.ID, expiry)
	if err != nil {
		return nil, err
	}

	err = s.userPort.CreateToken(ctx, dto.ManageTokenDto{ID: user.ID, AccessToken: token})
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponse{
		User: dto.UserResponseDto{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			EmailVerifiedAt: *user.EmailVerifiedAt,
		},
		Token: dto.TokenResponse{
			Type:        "Bearer",
			AccessToken: token,
		},
	}

	return &response, nil
}

// The [UpdateUser] usecase allows for users to update their bio-data when needed
func (s UserService) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	err := s.userPort.UpdateUser(ctx, id, updates)

	if err != nil {
		return err
	}

	return nil
}

// The [ForgetPassword] usecase allows for a user to forget a password and
// get a verification code sent to their email to reset the password
func (s UserService) ForgetPassword(ctx context.Context, payload dto.EmailDto) error {
	user, err := s.userPort.FindUser(ctx, dto.FindUserDto{Email: payload.Email})
	if err != nil {
		return err
	}

	// TODO: Send reset password to the user
	code := utils.EncodeToString(6)
	err = s.verificationPort.CreateCode(ctx, dto.VerificationDto{
		UserID: user.ID,
		Code:   code,
		Type:   "password-reset",
	})

	if err != nil {
		return err
	}

	return nil
}

// The [ResetPassword] usecase allows for a user to reset their password
func (s UserService) ResetPassword(ctx context.Context, payload dto.ResetPasswordDto) error {

	verificationCode, err := s.verificationPort.FindCode(ctx, dto.VerificationDto{Code: payload.Code, Type: "password-reset"})
	if err != nil {
		utils.TextLogger.Error("verification code not found", err)
		return err
	}

	user, err := s.userPort.FindUser(ctx, dto.FindUserDto{ID: verificationCode.UserId})
	if err != nil {
		utils.TextLogger.Error("user not found", err)
		return err
	}

	if user == nil {
		utils.TextLogger.Error("user not found, i.e returning nil with no error")
		return errors.New("either code is invalid or user does not exist")
	}

	if payload.Code != verificationCode.Code {
		utils.TextLogger.Error(fmt.Sprintf(
			"verification code mismatch. Given %v, expecting %v", payload.Code, verificationCode.Code,
		))
		return errors.New("either code is invalid or user does not exist")
	}

	password, err := s.passwordService.HashPassword(payload.Password)
	if err != nil {
		utils.TextLogger.Error("password mismatch", err)
		return err
	}

	updates := domain.Password{Value: password}
	err = s.userPort.UpdateUserPassword(ctx, user.ID, updates)
	if err != nil {
		utils.TextLogger.Error("update user password failed", err)
		return err
	}

	_ = s.verificationPort.DeleteCode(ctx, verificationCode.ID)

	return nil
}

// The [LogoutUser] usecase should delete a users token and de-authenticate them immediately
func (s UserService) LogoutUser(ctx context.Context, token string) error {
	err := s.userPort.DeleteToken(ctx, dto.ManageTokenDto{AccessToken: token})
	if err != nil {
		utils.TextLogger.Error("unable to delete user token", err)
		return err
	}

	return nil
}

func (s UserService) ShowProfile(ctx context.Context, id string) (*dto.UserResponseDto, error) {
	user, err := s.userPort.FindUser(ctx, dto.FindUserDto{ID: id})
	if err != nil {
		return nil, err
	}

	return &dto.UserResponseDto{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}

func (s UserService) VerifyEmailAddress(ctx context.Context, payload dto.VerificationDto) error {
	verificationCode, err := s.verificationPort.FindCode(ctx, payload)
	if err != nil {
		utils.TextLogger.Error("verification code not found", err)
		return err
	}

	user, err := s.userPort.FindUser(ctx, dto.FindUserDto{ID: verificationCode.UserId})
	if err != nil {
		utils.TextLogger.Error("unable to find the user", err)
		return err
	}

	if user == nil {
		utils.TextLogger.Error("user not found, i.e returning nil with no error")
		return errors.New("either code is invalid or user does not exist")
	}

	if payload.Code != verificationCode.Code {
		utils.TextLogger.Error(fmt.Sprintf(
			"verification code mismatch. Given %v, expecting %v", payload.Code, verificationCode.Code,
		))
		return errors.New("either code is invalid or user does not exist")
	}

	updates := map[string]interface{}{
		"email_verified_at": time.Now(),
	}

	err = s.userPort.UpdateUser(ctx, user.ID, updates)
	if err != nil {
		utils.TextLogger.Error("update user failed", err)
		return err
	}
	_ = s.verificationPort.DeleteCode(ctx, verificationCode.ID)

	return nil
}
