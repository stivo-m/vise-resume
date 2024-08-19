package services

import "golang.org/x/crypto/bcrypt"

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s PasswordService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s PasswordService) VerifyPassword(password string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	return err == nil
}
