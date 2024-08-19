package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/adapters/database/repository"
	"github.com/stivo-m/vise-resume/internal/adapters/http/handlers"
)

type Server struct {
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{db}
}

func (s *Server) PrepareServer() (*fiber.App, error) {
	app := fiber.New()

	s.db.AutoMigrate()

	// Repository
	userRepo := repository.NewUserRepository(s.db)
	verificationRepo := repository.NewVerificationRepository(s.db)

	// Services
	tokenService := NewTokenService()
	passwordService := NewPasswordService()
	userService := NewUserService(
		userRepo,
		tokenService,
		passwordService,
		verificationRepo,
	)

	// handlers
	api := app.Group("/api/v1")
	authHandlers := handlers.NewAuthHandler(userService)
	authHandlers.RegisterRoutes(api)

	return app, nil
}
