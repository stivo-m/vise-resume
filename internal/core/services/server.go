package services

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/adapters/database/repository"
	"github.com/stivo-m/vise-resume/internal/adapters/http/handlers"
)

// Create a slog logger, which:
//   - Logs to stdout.

type Server struct {
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{db}
}

func (s *Server) PrepareServer() (*fiber.App, error) {
	app := fiber.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.Use(slogfiber.New(logger))
	app.Use(recover.New())

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
	authHandlers := handlers.NewAuthHandler(userService, tokenService)
	authHandlers.RegisterRoutes(api)

	return app, nil
}
