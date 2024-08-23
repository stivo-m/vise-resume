package mocks

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/services"
)

func SetupTestServer() (*fiber.App, *database.DB, error) {

	err := os.Setenv("TOKEN_SECRET_KEY", "mockValue")
	if err != nil {
		return nil, nil, err
	}

	db, err := database.SetupMockDB()
	if err != nil {
		return nil, nil, err
	}

	server := services.NewServer(db)
	app, err := server.PrepareServer()
	if err != nil {
		return nil, nil, err
	}

	return app, db, nil
}
