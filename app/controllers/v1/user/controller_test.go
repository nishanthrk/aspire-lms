package user_controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/nishanthrk/aspire-lms/app/dto"
	userSvc "github.com/nishanthrk/aspire-lms/app/services/user"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := userSvc.NewMockUserService(ctrl)
	mockUserService.EXPECT().ValidateCredentials(gomock.Any()).Return(dto.LoginResponse{
		Status: 1,
		Data: dto.TokenDetail{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			UserID:       "user_id",
			UserType:     "EMPLOYEE",
			TokenExpires: 1718562847,
		},
	}, dto.HandleError{
		Status: 1,
	})

	app := fiber.New()
	app.Post("/v1/user/auth", func(c *fiber.Ctx) error {
		return Login(c, mockUserService)
	})

	request := dto.LoginRequest{
		Identifier: "9790970381",
		Password:   "12345678",
	}

	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/v1/user/auth", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Platform", "EMPLOYEE_API")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["status"].(float64))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.Equal(t, "user_id", data["user_id"])
	assert.Equal(t, "EMPLOYEE", data["user_type"])
	assert.Equal(t, float64(1718562847), data["token_expires"].(float64))
}

func TestLogin_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := userSvc.NewMockUserService(ctrl)

	app := fiber.New()
	app.Post("/v1/user/auth", func(c *fiber.Ctx) error {
		return Login(c, mockUserService)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/user/auth", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(-1), response["status"].(float64))
	assert.NotEmpty(t, response["error"])
}
