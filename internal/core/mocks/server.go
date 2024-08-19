package mocks

import (
	"github.com/gofiber/fiber/v2"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/services"
)

func SetupTestServer() (*fiber.App, error) {
	db, err := database.SetupMockDB()
	if err != nil {
		return nil, err
	}

	server := services.NewServer(db)
	app, err := server.PrepareServer()
	if err != nil {
		return nil, err
	}

	return app, nil
}
