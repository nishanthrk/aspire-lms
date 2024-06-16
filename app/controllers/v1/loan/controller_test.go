package loan_controller

import (
	"bytes"
	"encoding/json"
	"github.com/nishanthrk/aspire-lms/app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/nishanthrk/aspire-lms/app/dto"
	loanSvc "github.com/nishanthrk/aspire-lms/app/services/loan"
	repaymentSvc "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userSvc "github.com/nishanthrk/aspire-lms/app/services/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoanApplication_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoanService := loanSvc.NewMockLoanService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)
	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)

	mockLoanService.EXPECT().CreateLoanApplication(gomock.Any(), mockUserService, mockRepaymentService).Return(dto.ApplicationCreateResponse{
		Status: 1,
		Data: struct {
			ApplicationID string `json:"application_id"`
			UserID        string `json:"user_id"`
		}{
			ApplicationID: "application_id",
			UserID:        "user_id",
		},
		Message: "Loan application created successfully",
	}, dto.HandleError{
		Status: 1,
	})

	app := fiber.New()
	app.Post("/v1/application", func(c *fiber.Ctx) error {
		return CreateLoanApplication(c, mockLoanService, mockUserService, mockRepaymentService)
	})

	request := dto.ApplicationCreateRequest{
		User: dto.UserObject{
			UserName:     "John Doe",
			UserEmail:    "john.doe@example.com",
			MobileNumber: "1234567890",
			CountryCode:  "IND",
			Kyc: dto.KycObject{
				KycType:   "PAN",
				KycNumber: "ABCDE1234F",
			},
		},
		LoanApplication: dto.LoanApplicationObject{
			LoanAmount:    100000.00,
			CurrencyCode:  "INR",
			InterestRate:  7.5,
			LoanTerm:      12,
			LoanTermUnit:  "WEEKLY",
			Income:        600000.00,
			CreditScore:   750,
			ExistingDebts: 100000.00,
			CountryCode:   "IND",
		},
	}

	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/v1/application", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["status"].(float64))
	assert.Equal(t, "Loan application created successfully", response["message"].(string))
	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["application_id"])
	assert.NotEmpty(t, data["user_id"])
}

func TestCreateLoanApplication_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoanService := loanSvc.NewMockLoanService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)
	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)

	app := fiber.New()
	app.Post("/v1/application", func(c *fiber.Ctx) error {
		return CreateLoanApplication(c, mockLoanService, mockUserService, mockRepaymentService)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/application", bytes.NewReader([]byte(`invalid json`)))
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

func TestApproveLoanApplication_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoanService := loanSvc.NewMockLoanService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)
	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)

	mockUserService.EXPECT().GetUserObject(gomock.Any()).Return(models.User{
		UserID:   "53297921-01d9-4311-94f3-54cbb971c5a0",
		UserName: "John Doe",
		UserType: "EMPLOYEE",
	})
	mockLoanService.EXPECT().ApproveLoanApplication(gomock.Any(), gomock.Any(), mockRepaymentService).Return(dto.ApplicationApproveResponse{
		Status: 1,
		Data: dto.ApproveObject{
			ApplicationID: "2e560272-45eb-442e-9a67-a2f8c423e063",
		},
		Message: "Loan application approved successfully",
	}, dto.HandleError{
		Status: 1,
	})

	app := fiber.New()
	app.Post("/v1/application/:applicationId/approve", func(c *fiber.Ctx) error {
		return ApproveLoanApplication(c, mockLoanService, mockUserService, mockRepaymentService)
	})

	request := dto.ApplicationApproveRequest{
		ApprovedAmount: 100000,
		Override:       true,
	}

	requestBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/v1/application/2e560272-45eb-442e-9a67-a2f8c423e063/approve", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkFwcHJvdmVyIE9uZSIsInVzZXJfdHlwZSI6IkVNUExPWUVFIiwidXNlcl9pZCI6IjUzMjk3OTIxLTAxZDktNDMxMS05NGYzLTU0Y2JiOTcxYzVhMCIsImlzcyI6IkFTUElSRSIsImV4cCI6MTcxODUxNDA2OX0.QuKDYgSExuTEFatB149t-Upf5pFzbDzAVR3Ar7YDRyY")
	req.Header.Set("X-Platform", "CUSTOMER_API")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["status"].(float64))
	assert.Equal(t, "Loan application approved successfully", response["message"].(string))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "2e560272-45eb-442e-9a67-a2f8c423e063", data["application_id"])
}

func TestApproveLoanApplication_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLoanService := loanSvc.NewMockLoanService(ctrl)
	mockUserService := userSvc.NewMockUserService(ctrl)
	mockRepaymentService := repaymentSvc.NewMockRepaymentService(ctrl)

	app := fiber.New()
	app.Post("/v1/application/:applicationId/approve", func(c *fiber.Ctx) error {
		return ApproveLoanApplication(c, mockLoanService, mockUserService, mockRepaymentService)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/application/2e560272-45eb-442e-9a67-a2f8c423e063/approve", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkFwcHJvdmVyIE9uZSIsInVzZXJfdHlwZSI6IkVNUExPWUVFIiwidXNlcl9pZCI6IjUzMjk3OTIxLTAxZDktNDMxMS05NGYzLTU0Y2JiOTcxYzVhMCIsImlzcyI6IkFTUElSRSIsImV4cCI6MTcxODUxNDA2OX0.QuKDYgSExuTEFatB149t-Upf5pFzbDzAVR3Ar7YDRyY")
	req.Header.Set("X-Platform", "CUSTOMER_API")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, float64(-1), response["status"].(float64))
	assert.NotEmpty(t, response["error"])
}
