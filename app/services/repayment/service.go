package repayment_service

import (
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
)

// RepaymentService defines the interface for repayment-related operations
type RepaymentService interface {
	// UpdateRepayment Processes a repayment request and updates the repayment schedule
	UpdateRepayment(request dto.RepaymentRequest, user models.User) (dto.RepaymentResponse, dto.HandleError)

	// CalculateRepaymentSchedule Calculates the repayment schedule for a given loan application
	CalculateRepaymentSchedule(application *models.LoanApplication) ([]models.Repayment, error)
}

// repaymentService is an implementation of RepaymentService
type repaymentService struct{}

// NewRepaymentService returns a new instance of RepaymentService
func NewRepaymentService() RepaymentService {
	return &repaymentService{}
}
