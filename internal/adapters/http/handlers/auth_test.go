package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stivo-m/vise-resume/internal/core/domain"
	"github.com/stivo-m/vise-resume/internal/core/mocks"
	"github.com/stretchr/testify/assert"
)

func GenerateFakeUser() domain.User {
	return domain.User{
		Base: domain.Base{
			ID:        gofakeit.UUID(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		},
		FullName: gofakeit.Name(),
		Email:    gofakeit.Email(),
	}
}

func TestUserRegistrationWithoutBody(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", nil)
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"The request body is invalid"`)
}

func TestUserRegistrationInvalidBody(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	payload := `{?}`

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"The request body is invalid"`)
}

func TestUserRegistrationWithValidationErrors(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	payload := `{}`

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"one or more of the required fields are invalid or missing"`)
}

func TestUserRegistrationSuccessful(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	email := gofakeit.Email()
	payload := fmt.Sprintf(`{"full_name":"John Doe","email":"%v", "password": "password"}`, email)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"User was registered successfully"`)
	assert.Contains(t, string(body), `"id"`)
}

func TestUserRegistrationFailsForSameEmail(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	email := gofakeit.Email()
	payload := fmt.Sprintf(`{"full_name":"John Doe","email":"%v", "password": "password"}`, email)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	req = httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"Registration failed"`)
	assert.NotContains(t, string(body), `"id"`)
}

// Login Tests

func TestLoginWithNonExistingUser(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	email := gofakeit.Email()
	payload := fmt.Sprintf(`{"email":"%v", "password": "password"}`, email)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"message":"Login failed"`)
}

func TestLoginWithUnverifiedEmail(t *testing.T) {
	app, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	email := gofakeit.Email()
	payload := fmt.Sprintf(`{"full_name":"John Doe","email":"%v", "password": "password"}`, email)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	payload = fmt.Sprintf(`{"email":"%v", "password": "password"}`, email)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"email address is not verified"`)
}
