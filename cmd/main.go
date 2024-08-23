package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/adapters/database"
	"github.com/stivo-m/vise-resume/internal/core/services"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

func setupPostmanCollections(app *fiber.App, port int) {
	// Generate the Postman collection
	collection := utils.GeneratePostmanCollection(app, port)

	// Write the collection to a file
	file, err := os.Create("postman_collection.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print the JSON
	if err := encoder.Encode(collection); err != nil {
		fmt.Println("Error writing JSON to file:", err)
	}

	fmt.Println("Postman collection has been generated and saved to postman_collection.json")
}

func main() {
	db, err := database.NewDatabase()
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

	// List available routes
	if len(os.Args) > 1 && os.Args[1] == "list:routes" {
		utils.ListRoutes(app)
		return
	}

	// Generate postman collections
	if len(os.Args) > 1 && os.Args[1] == "generate:postman" {
		setupPostmanCollections(app, port)
		return
	}

	// Run migrations
	if len(os.Args) > 1 && os.Args[1] == "migrations:run" {
		db.AutoMigrate()
		return
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
