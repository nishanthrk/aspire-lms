package loan_service

import (
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
	repaymentService "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
)

// LoanService defines the interface for loan-related operations
type LoanService interface {
	// CreateLoanApplication creates a new loan application
	CreateLoanApplication(params dto.ApplicationCreateRequest, userSvc userService.UserService,
		repaymentSvc repaymentService.RepaymentService) (dto.ApplicationCreateResponse, dto.HandleError)
	// ApproveLoanApplication approves an existing loan application
	ApproveLoanApplication(params dto.ApplicationApproveRequest, user models.User,
		repaymentSvc repaymentService.RepaymentService) (dto.ApplicationApproveResponse, dto.HandleError)
	// GetLoanApplication retrieves the details of an existing loan application
	GetLoanApplication(applicationId string, user models.User) (dto.ApplicationDetailsResponse, dto.HandleError)

	// GetParticipantApplications retrieves the applications of participant
	GetParticipantApplications(user models.User) dto.ApplicationListResponse
}

// loanService is the implementation of the LoanService interface
type loanService struct{}

// NewLoanService returns a new instance of LoanService
func NewLoanService() LoanService {
	return &loanService{}
}
