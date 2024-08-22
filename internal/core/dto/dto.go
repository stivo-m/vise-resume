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
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordDto struct {
	Code     string `json:"code"  validate:"required,min=6,max=6"`
	Password string `json:"password"  validate:"required,min=5,max=255"`
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
	ID              string    `json:"id"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	EmailVerifiedAt time.Time `json:"email_verified_at,omitempty"`
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
	Email  string `json:"email" validate:"email,required"`
	Code   string `json:"code" validate:"required,min=6,max=6"`
	Type   string `json:"type" validate:"required,oneof=email-verification password-reset"`
}

// PostmanCollection represents the structure of a Postman collection.
type PostmanCollection struct {
	Info     PostmanInfo       `json:"info"`
	Item     []PostmanItem     `json:"item"`
	Variable []PostmanVariable `json:"variable,omitempty"`
}

type PostmanInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type PostmanItem struct {
	Name    string          `json:"name"`
	Item    []PostmanItem   `json:"item,omitempty"`    // For grouping (folders)
	Request *PostmanRequest `json:"request,omitempty"` // Set as a pointer to omit if nil
}

type PostmanRequest struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header,omitempty"`
	Url    PostmanUrl      `json:"url"`
	Body   *PostmanBody    `json:"body,omitempty"`
}

type PostmanUrl struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}
