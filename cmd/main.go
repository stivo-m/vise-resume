package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/services"
)

func main() {
	db, err := database.SetupMockDB()
	if err != nil {
		log.Panicf("unable to connect to the database: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Panicf("unable to parse SERVER_PORT: %v", err)
	}

	server := services.NewServer(db)
	app, err := server.PrepareServer()

	if err != nil {
		log.Panicf("unable to prepare server: %v", err)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
