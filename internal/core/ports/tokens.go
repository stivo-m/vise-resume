package ports

import "time"

type TokenService interface {
	CreateToken(id string, expiryDate time.Time) (string, error)
	VerifyToken(token string) (string, error)
}
