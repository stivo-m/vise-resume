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

type Server struct {
	db *database.DB
}

func NewServer(db *database.DB) *Server {
	return &Server{db}
}

func (s *Server) PrepareServer() (*fiber.App, error) {
	app := fiber.New()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.Use(slogfiber.New(logger))
	app.Use(recover.New())

	// Repository
	userRepo := repository.NewUserRepository(s.db)
	verificationRepo := repository.NewVerificationRepository(s.db)
	resumeRepo := repository.NewResumeRepository(s.db)

	// Services
	tokenService := NewTokenService()
	passwordService := NewPasswordService()
	userService := NewUserService(
		userRepo,
		tokenService,
		passwordService,
		verificationRepo,
	)
	resumeService := NewResumeService(resumeRepo)

	// handlers
	api := app.Group("/api/v1")
	authHandlers := handlers.NewAuthHandler(userService, userRepo, tokenService)
	authHandlers.RegisterAuthRoutes(api)

	resumeHandler := handlers.NewResumeHandler(resumeService, userRepo, tokenService)
	resumeHandler.RegisterResumeRoutes(api)

	return app, nil
}
