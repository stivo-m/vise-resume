package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stivo-m/vise-resume/internal/core/mocks"
	"github.com/stivo-m/vise-resume/internal/core/test"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistrationWithoutBody(t *testing.T) {
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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
	app, _, err := mocks.SetupTestServer()
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

func TestLoginWithVerifiedEmail(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	user, _, _ := test.GetAuthenticatedTestUser(db)

	requestBody := fmt.Sprintf(`{"email":"%v", "password": "password"}`, user.Email)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"User was logged in successfully"`)
	assert.NotContains(t, string(body), `"email address is not verified"`)
}

func TestShowProfile(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, token, _ := test.GetAuthenticatedTestUser(db)

	req := httptest.NewRequest("GET", "/api/v1/auth/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"User profile"`)
	assert.NotContains(t, string(body), `"Unable to show profile"`)
}

func TestUpdateUserName(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, token, _ := test.GetAuthenticatedTestUser(db)

	requestBody := fmt.Sprintf(`{"full_name":"%v"}`, "Random name")
	req := httptest.NewRequest("PATCH", "/api/v1/auth/profile", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"User updated successfully"`)
	assert.NotContains(t, string(body), `"User update failed"`)
}
