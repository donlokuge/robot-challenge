package robot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	repo := NewMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	app.Post("/commands", handler.SubmitCommands)
	app.Get("/status/:taskID", handler.GetStatus)
	app.Delete("/commands/:taskID", handler.CancelTask)
	app.Get("/robots", handler.GetRobots)

	return app
}

func TestSubmitCommands(t *testing.T) {
	app := setupTestApp()

	payload := SubmitCommandsRequest{Commands: "N E S W"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/commands", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)

	var response SubmitCommandsResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response.TaskID)
}

func TestGetRobots(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/robots", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSubmitCommands_Invalid(t *testing.T) {
	app := setupTestApp()

	payload := SubmitCommandsRequest{Commands: "N X"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/commands", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
