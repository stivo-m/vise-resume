package services

import "github.com/stivo-m/vise-resume/internal/core/ports"

type VerificationService struct {
	verificationPort ports.VerificationPort
}

func NewVerificationService(
	verificationPort ports.VerificationPort,

) *VerificationService {
	return &VerificationService{
		verificationPort: verificationPort,
	}
}
