package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func FormatApiResponse(message string, data interface{}) dto.ApiResponse[any] {
	var errorMap dto.ApiResponse[any]
	errorMap.Message = message
	errorMap.Data = data

	return errorMap
}

func GetValidationMessage(err validator.FieldError, field string) string {
	switch err.Tag() {
	case "required":
		return "The '" + field + "' field is required"
	case "min":
		return "The '" + field + "' field must be at least " + err.Param() + " characters long"
	case "oneof":
		return "The '" + field + "' field should be one of " + err.Param()
	default:
		return "The '" + field + "' field is invalid"
	}
}

func GetJSONFieldName(dto interface{}, structField string) string {
	r := reflect.TypeOf(dto).Elem()
	field, _ := r.FieldByName(structField)
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return strings.ToLower(structField)
	}
	return strings.Split(jsonTag, ",")[0]
}

func ExtractUuidFromContext(ctx context.Context, key string) (*uuid.UUID, error) {
	id, ok := ctx.Value(key).(string)

	if !ok {
		return nil, errors.New("given key is not set")
	}

	userId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	return &userId, nil
}

// Function to list routes
func ListRoutes(app *fiber.App) {
	routes := app.GetRoutes()
	fmt.Printf("%-10s %-50s\n", "METHOD", "PATH")
	fmt.Println("---------------------------------------------------------------")

	for _, route := range routes {
		fmt.Printf("%-10s %-50s\n", route.Method, route.Path)
	}
}

// Helper function to find or create a folder by name
func findOrCreateFolder(folders *[]dto.PostmanItem, name string) *dto.PostmanItem {
	for i := range *folders {
		if (*folders)[i].Name == name {
			return &(*folders)[i]
		}
	}
	// If folder doesn't exist, create it
	folder := dto.PostmanItem{Name: name}
	*folders = append(*folders, folder)
	return &(*folders)[len(*folders)-1]
}

// Helper function to generate a more readable name from a route path
func generateReadableName(routePath string) string {
	// Trim the /api/v1 prefix
	trimmedPath := strings.TrimPrefix(routePath, "/api/v1")
	parts := strings.Split(trimmedPath, "/")

	// Capitalize each part and join them with spaces
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	return strings.Join(parts, " - ") // Join with a dash for readability
}

func GeneratePostmanCollection(app *fiber.App, port int) dto.PostmanCollection {
	routes := app.GetRoutes()
	url := os.Getenv("SERVER_URL")

	// Create base info for Postman Collection
	collection := dto.PostmanCollection{
		Info: dto.PostmanInfo{
			Name:   "Vise Resume API Collection",
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []dto.PostmanVariable{
			{
				Key:   "base_url",
				Value: fmt.Sprintf("%s:%d", url, port) + "/api/v1", // Set this to your desired base URL
			},
			{
				Key:   "access_token",
				Value: "", // Access token can be set dynamically via environment variables
			},
		},
	}

	// Group items under `/api/v1`
	apiGroup := dto.PostmanItem{
		Name: "API v1",
	}

	// Iterate over Fiber routes and group them based on path
	for _, route := range routes {
		// Skip routes that don't start with "/api/v1"
		if !strings.HasPrefix(route.Path, "/api/v1") {
			continue
		}

		// Determine the prefix group for the folder (e.g., /api/v1/auth or /api/v1/resume)
		// Only append the path after /api/v1
		trimmedPath := strings.TrimPrefix(route.Path, "/api/v1")
		parts := strings.Split(trimmedPath, "/")
		if len(parts) < 2 {
			continue
		}
		groupName := parts[1]

		// Find or create the group folder
		folder := findOrCreateFolder(&collection.Item, groupName)

		readableName := generateReadableName(route.Path)

		// Create a Postman item for each route
		item := dto.PostmanItem{
			Name: readableName,
			Request: &dto.PostmanRequest{
				Method: route.Method,
				Url: dto.PostmanUrl{
					Raw:  fmt.Sprintf("{{base_url}}%s", trimmedPath),
					Host: []string{"{{base_url}}"},
					Path: parts,
				},
				Header: []dto.PostmanHeader{
					{
						Key:   "Authorization",
						Value: "Bearer {{access_token}}",
						Type:  "text",
					},
				},
			},
		}

		// Add the route to the appropriate folder
		folder.Item = append(folder.Item, item)
	}

	// Add the API group to the collection if it has items
	if len(apiGroup.Item) > 0 {
		collection.Item = append(collection.Item, apiGroup)
	}

	return collection
}
