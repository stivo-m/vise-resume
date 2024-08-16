package ports

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password string, encryptedPassword string) bool
}
