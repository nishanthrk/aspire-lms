package repayment_controller

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
	repaymentSvc "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userSvc "github.com/nishanthrk/aspire-lms/app/services/user"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPayRepayment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)

	mockRepaymentService.EXPECT().UpdateRepayment(gomock.Any(), gomock.Any()).Return(dto.RepaymentResponse{
		Status:  1,
		Message: "Payment process successfully",
		Data: struct {
			PaymentId string `json:"payment_id"`
		}{
			PaymentId: "payment_id",
		},
	}, dto.HandleError{
		Status: 1,
	})

	mockUserService.EXPECT().GetUserObject(gomock.Any()).Return(models.User{
		UserID: "user_id",
	})

	app := fiber.New()
	app.Post("/application/:applicationId/repayment", func(c *fiber.Ctx) error {
		return PayRepayment(c, mockRepaymentService, mockUserService)
	})

	request := dto.RepaymentRequest{
		ApplicationID: "application_id",
		PaymentAmount: 1000.0,
	}

	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/application/application_id/repayment", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["status"].(float64))
	assert.Equal(t, "Payment process successfully", response["message"].(string))
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["payment_id"])
}

func TestPayRepayment_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)

	app := fiber.New()
	app.Post("/application/:applicationId/repayment", func(c *fiber.Ctx) error {
		return PayRepayment(c, mockRepaymentService, mockUserService)
	})

	req := httptest.NewRequest(http.MethodPost, "/application/application_id/repayment", bytes.NewReader([]byte(`invalid json`)))
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
