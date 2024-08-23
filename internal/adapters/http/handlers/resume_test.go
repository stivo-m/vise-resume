package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stivo-m/vise-resume/internal/core/mocks"
	"github.com/stivo-m/vise-resume/internal/core/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateResumeSuccessful(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, token, err := test.GetAuthenticatedTestUser(db)
	assert.Nil(t, err)

	payload := `{"summary":"Test","skills":["Javascript","Go"],"experience":[{"company_name":"Test Company","role":"Software Engineer","start_date":"2019-01-02T15:04:05Z"}],"education":[{"school_name":"Test School","course":"Test Course","start_date":"2006-01-02T15:04:05Z"}]}`

	req := httptest.NewRequest("POST", "/api/v1/resume/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"Resume was created successfully"`)
	assert.NotContains(t, string(body), `"Resume creation failed"`)
}

func TestCreateResumeFailedWhenUnauthenticated(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, _, err = test.GetAuthenticatedTestUser(db)
	assert.Nil(t, err)

	payload := `{"summary":"Test","skills":["Javascript","Go"],"experience":[{"company_name":"Test Company","role":"Software Engineer","start_date":"2019-01-02T15:04:05Z"}],"education":[{"school_name":"Test School","course":"Test Course","start_date":"2006-01-02T15:04:05Z"}]}`

	req := httptest.NewRequest("POST", "/api/v1/resume/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"unauthorized"`)
}

func TestCreateResumeValidationErrors(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, token, err := test.GetAuthenticatedTestUser(db)
	assert.Nil(t, err)

	payload := `{"summary":"Test","skills":["Javascript","Go"],"experience":[{"company_name":"Test Company","start_date":"2019-01-02T15:04:05Z"}],"education":[{"school_name":"Test School","course":"Test Course"}]}`

	req := httptest.NewRequest("POST", "/api/v1/resume/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.NotContains(t, string(body), `"Resume was created successfully"`)
	assert.Contains(t, string(body), `"one or more of the required fields are invalid or missing"`)
}

func TestListResumeSuccess(t *testing.T) {
	app, db, err := mocks.SetupTestServer()
	assert.Nil(t, err)

	_, token, err := test.GetAuthenticatedTestUser(db)
	assert.Nil(t, err)
	payload := `{"summary":"Test","skills":["Javascript","Go"],"experience":[{"company_name":"Test Company","role":"Software Engineer","start_date":"2019-01-02T15:04:05Z"}],"education":[{"school_name":"Test School","course":"Test Course","start_date":"2006-01-02T15:04:05Z"}]}`

	req := httptest.NewRequest("POST", "/api/v1/resume/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"Resume was created successfully"`)
	assert.NotContains(t, string(body), `"Resume creation failed"`)

	req = httptest.NewRequest("GET", "/api/v1/resume/list", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err = app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	body = make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Contains(t, string(body), `"Resumes obtained successfully"`)
	assert.NotContains(t, string(body), `"Unable to find resumes"`)
}
