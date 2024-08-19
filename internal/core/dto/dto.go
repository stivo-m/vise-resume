package dto

import "time"

type RegisterDto struct {
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailDto struct {
	Email string `json:"email"`
}

type ResetPasswordDto struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type FindUserDto struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	WithPassword bool
}

type ManageTokenDto struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}

type UpdateUserDto struct {
	FullName string `json:"full_name"`
}

type UserResponseDto struct {
	ID              string     `json:"id"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"type"`
}

type LoginResponse struct {
	User  UserResponseDto `json:"user"`
	Token TokenResponse   `json:"token"`
}

type ProfileResponse struct {
	User UserResponseDto `json:"user"`
}

type ApiResponse[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

type VerificationDto struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
	Type   string `json:"type"`
}
