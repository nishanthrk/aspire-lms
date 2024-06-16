package loan_service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nishanthrk/aspire-lms/app/common/constants"
	"github.com/nishanthrk/aspire-lms/app/common/utility"
	db "github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
	repaymentService "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
	"time"
)

// ApproveLoanApplication approves a loan application and generates the repayment schedule if needed
// Parameters:
// - request: dto.ApplicationApproveRequest containing the application ID, approved amount, and override flag
// - user: models.User representing the user performing the approval
// - repaymentSvc: repaymentService.RepaymentService for generating the repayment schedule
// Returns:
// - dto.ApplicationApproveResponse with the approval result
// - dto.HandleError with any error that occurred during the process
func (s *loanService) ApproveLoanApplication(request dto.ApplicationApproveRequest, user models.User, repaymentSvc repaymentService.RepaymentService) (
	response dto.ApplicationApproveResponse, handle dto.HandleError) {

	// Find the application by its primary key
	application := models.LoanApplication{}
	application, _ = application.FindByPrimaryKey(request.ApplicationID)
	if application.ApplicationID == "" {
		handle.Status = -1
		handle.Errors = fmt.Errorf("application details not found")
		return
	}

	// Check if the application is already processed
	if application.Status != models.LoanApplicationStatusPending {
		response = dto.ApplicationApproveResponse{
			Status:  1,
			Message: "Loan application already processed",
			Data: dto.ApproveObject{
				ApplicationID: application.ApplicationID,
			},
		}
		return
	}

	// Validate the approved amount
	if request.ApprovedAmount > application.EligibleLoanAmount && !request.Override {
		handle.Status = -3
		handle.Errors = fmt.Errorf("application already processed")
		return
	}

	// Check if the user has permission to approve the application
	var participantCondition []db.WhereCondition

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.UserID,
		Condition: "=",
		Value:     user.UserID,
	})

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.ApplicationID,
		Condition: "=",
		Value:     application.ApplicationID,
	})

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.ParticipantType,
		Condition: "=",
		Value:     constants.UserTypeEmployee,
	})

	participant := models.LoanApplicationParticipant{}
	participant, _ = participant.FindOneByCondition(participantCondition)

	if participant.ParticipantID == "" {
		handle.Status = -4
		handle.Errors = fmt.Errorf("dont have permission to approve the application")
		return
	}

	tx := db.MysqlDB.Begin()

	// Update the application status and approval details
	application.Status = models.LoanApplicationStatusApproved
	application.ApprovedAmount = null.FloatFrom(request.ApprovedAmount)
	application.ApprovedDate = null.TimeFrom(time.Now())

	// Check if the application needs a new repayment schedule
	difference := application.ApprovedDate.Time.Sub(application.ApplicationDate)

	// If the approved amount and loan amount not same or loan is approved beyond application date regenerate repayment
	if ((difference.Hours() / 24) > 1) || request.ApprovedAmount != application.LoanAmount {

		var repaymentCondition []db.WhereCondition
		repaymentCondition = append(repaymentCondition, db.WhereCondition{
			Key:       models.RepaymentColumns.ApplicationID,
			Condition: "=",
			Value:     application.ApplicationID,
		})

		repayment := models.Repayment{}
		oldRepayment, _ := repayment.FindAllByCondition(repaymentCondition)

		// Delete all the old repayment of the application
		if len(oldRepayment) > 0 {
			err := tx.Delete(&oldRepayment).Error
			if err != nil {
				tx.Rollback()
				handle.Status = -5
				handle.Errors = err
				return
			}
		}

		// Generate new repayment for application for approved date
		repayments, err := repaymentSvc.CalculateRepaymentSchedule(&application)
		if err != nil {
			tx.Rollback()
			handle.Status = -5
			handle.Errors = err
			return
		}

		err = tx.Save(&repayments).Error
		if err != nil {
			tx.Rollback()
			handle.Status = -5
			handle.Errors = err
			return
		}
	}

	err := tx.Save(&application).Error
	if err != nil {
		tx.Rollback()
		handle.Status = -5
		handle.Errors = err
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		handle.Status = -6
		handle.Errors = err
		return
	}

	response = dto.ApplicationApproveResponse{
		Status:  1,
		Message: "Loan application approved successfully",
		Data: dto.ApproveObject{
			ApplicationID: application.ApplicationID,
		},
	}

	return
}

// GetParticipantApplications retrieves the list of loan applications in which the user is a participant.
// Parameters:
// - user: models.User representing the user for whom to retrieve the applications.
// Returns:
// - dto.ApplicationListResponse with the list of loan applications and a status message.
func (s *loanService) GetParticipantApplications(user models.User) (response dto.ApplicationListResponse) {
	// Define conditions to check if the user has participated in any loan applications
	var participantCondition []db.WhereCondition

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.UserID,
		Condition: "=",
		Value:     user.UserID,
	})

	// Find all loan applications where the user is a participant
	participant := models.LoanApplicationParticipant{}
	participants, _ := participant.FindAllByCondition(participantCondition)

	// Prepare the response data
	var data []dto.LoanApplicationObject
	for _, p := range participants {
		data = append(data, p.LoanApplication.GetLoanApplicationDTO())
	}

	// Set response status and message based on whether the user has participated in any applications
	if len(data) > 0 {
		response.Status = 1
		response.Message = fmt.Sprintf("Total %v application(s) found", len(data))
		response.Data = &data
	} else {
		response.Status = 1
		response.Message = "You have not participated in any application"
	}
	return
}

// GetLoanApplication retrieves the loan application details and its repayment schedule
// Parameters:
// - applicationId: string containing the application ID
// - user: models.User representing the user requesting the details
// Returns:
// - dto.ApplicationDetailsResponse with the application details and repayment schedule
// - dto.HandleError with any error that occurred during the process
func (s *loanService) GetLoanApplication(applicationId string, user models.User) (
	response dto.ApplicationDetailsResponse, handle dto.HandleError) {

	application := models.LoanApplication{}
	application, _ = application.FindByPrimaryKey(applicationId)
	if application.ApplicationID == "" {
		handle.Status = -1
		handle.Errors = fmt.Errorf("application %v not found", application.ApplicationID)
		return
	}

	var participantCondition, repaymentCondition []db.WhereCondition

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.UserID,
		Condition: "=",
		Value:     user.UserID,
	})

	participantCondition = append(participantCondition, db.WhereCondition{
		Key:       models.LoanApplicationParticipantColumns.ApplicationID,
		Condition: "=",
		Value:     application.ApplicationID,
	})

	participant := models.LoanApplicationParticipant{}
	participant, _ = participant.FindOneByCondition(participantCondition)

	if participant.ParticipantID == "" {
		handle.Status = -3
		handle.Errors = fmt.Errorf("dont have permission to this application: %v", application.ApplicationID)
		return
	}

	repaymentCondition = append(repaymentCondition, db.WhereCondition{
		Key:       models.RepaymentColumns.ApplicationID,
		Condition: "=",
		Value:     application.ApplicationID,
	})

	repayment := models.Repayment{}
	repayments, _ := repayment.FindAllByCondition(repaymentCondition)

	var repaymentDTOs []dto.Repayment
	for _, r := range repayments {
		repaymentDTOs = append(repaymentDTOs, r.GetRepaymentDTO())
	}

	response = dto.ApplicationDetailsResponse{
		Status: 1,
	}
	response.Data.LoanApplicationObject = application.GetLoanApplicationDTO()
	response.Data.Repayment = repaymentDTOs

	return
}

// CreateLoanApplication creates a new loan application and its initial repayment schedule
// Parameters:
// - request: dto.ApplicationCreateRequest containing the user and loan application details
// - userSvc: userService.UserService for handling user-related operations
// - repaymentSvc: repaymentService.RepaymentService for generating the repayment schedule
// Returns:
// - dto.ApplicationCreateResponse with the creation result
// - dto.HandleError with any error that occurred during the process
func (s *loanService) CreateLoanApplication(request dto.ApplicationCreateRequest,
	userSvc userService.UserService, repaymentSvc repaymentService.RepaymentService) (
	response dto.ApplicationCreateResponse, handle dto.HandleError) {
	tx := db.MysqlDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			handle.Status = -1
			handle.Errors = fmt.Errorf("%v", r)
		}
	}()

	// Check if the user exist
	user := userSvc.GetApplicationUser(request.User)
	if user.UserID == "" {
		user = models.User{
			UserID:       uuid.New().String(),
			UserType:     constants.UserTypeCustomer,
			UserName:     request.User.UserName,
			UserEmail:    request.User.UserEmail,
			MobileNumber: request.User.MobileNumber,
		}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			handle.Status = -2
			handle.Errors = err
			return
		}
	}

	// Insert user identification
	identification := models.UserKyc{
		KycID:       uuid.New().String(),
		UserID:      user.UserID,
		KycType:     request.User.Kyc.KycType,
		KycNumber:   request.User.Kyc.KycNumber,
		CountryCode: request.User.CountryCode,
	}
	if err := tx.Create(&identification).Error; err != nil {
		tx.Rollback()
		handle.Status = -3
		handle.Errors = err
		return
	}

	// Insert loan application
	loanApplication := models.LoanApplication{
		ApplicationID:      uuid.New().String(),
		LoanAmount:         request.LoanApplication.LoanAmount,
		CurrencyCode:       request.LoanApplication.CurrencyCode,
		InterestRate:       request.LoanApplication.InterestRate,
		LoanTerm:           request.LoanApplication.LoanTerm,
		LoanTermUnit:       request.LoanApplication.LoanTermUnit,
		Income:             request.LoanApplication.Income,
		CreditScore:        request.LoanApplication.CreditScore,
		ExistingDebts:      request.LoanApplication.ExistingDebts,
		CountryCode:        request.LoanApplication.CountryCode,
		ApplicationDate:    time.Now(),
		Status:             models.LoanApplicationStatusPending,
		EligibleLoanAmount: calculateEligibleLoanAmount(request.LoanApplication),
	}

	// Generate repayment for application
	repayments, err := repaymentSvc.CalculateRepaymentSchedule(&loanApplication)
	if err != nil {
		tx.Rollback()
		handle.Status = -5
		handle.Errors = err
		return
	}

	// Repayment date might be updated so saving it again
	if err := tx.Create(&loanApplication).Error; err != nil {
		tx.Rollback()
		handle.Status = -4
		handle.Errors = err
		return
	}

	if err := tx.Create(&repayments).Error; err != nil {
		tx.Rollback()
		handle.Status = -6
		handle.Errors = err
		return
	}

	var participants []models.LoanApplicationParticipant

	// Insert customer to application participant
	participants = append(participants, models.LoanApplicationParticipant{
		ParticipantID:   uuid.New().String(),
		ApplicationID:   loanApplication.ApplicationID,
		ParticipantType: constants.UserTypeCustomer,
		UserID:          user.UserID,
	})

	// Fetch employee for processing the application
	approveParticipant, err := userSvc.AllocateEmployeeForProcess(loanApplication)
	if err != nil {
		tx.Rollback()
		handle.Status = -7
		handle.Errors = err
		return
	}

	// Insert employee for approving to application participant
	participants = append(participants, approveParticipant)

	if err := tx.Create(&participants).Error; err != nil {
		tx.Rollback()
		handle.Status = -8
		handle.Errors = err
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		handle.Status = -9
		handle.Errors = err
		return
	}

	response.Status = 1
	response.Message = "Loan application created successfully"
	response.Data.ApplicationID = utility.ToString(loanApplication.ApplicationID)
	response.Data.UserID = utility.ToString(user.UserID)

	return
}

// calculateEligibleLoanAmount calculates the eligible loan amount based on the applicant's income, credit score, and existing debts
// Parameters:
// - application: dto.LoanApplicationObject containing the loan application details
// Returns:
// - float64 representing the eligible loan amount
func calculateEligibleLoanAmount(application dto.LoanApplicationObject) (eligibleLoanAmount float64) {
	income := application.Income
	creditScore := application.CreditScore
	existingDebts := application.ExistingDebts
	foir := existingDebts / income

	config := models.LoanEligibilityConfig{}

	var condition []db.WhereCondition

	condition = append(condition, db.WhereCondition{
		Key:       models.LoanEligibilityConfigColumns.CountryCode,
		Condition: "=",
		Value:     application.CountryCode,
	})

	condition = append(condition, db.WhereCondition{
		Key:       models.LoanEligibilityConfigColumns.MaxFoir,
		Condition: ">",
		Value:     foir,
	})

	condition = append(condition, db.WhereCondition{
		Key:       models.LoanEligibilityConfigColumns.MinCreditScore,
		Condition: "<=",
		Value:     creditScore,
	})

	condition = append(condition, db.WhereCondition{
		Key:       models.LoanEligibilityConfigColumns.MaxCreditScore,
		Condition: ">=",
		Value:     creditScore,
	})

	config, _ = config.FindOneByCondition(condition)

	if config.ID != 0 {
		maxLoanAmount := (income * config.MaxFoir) - existingDebts

		var creditScoreFactor float64
		switch {
		case creditScore >= ((config.MaxCreditScore + config.MinCreditScore) / 2):
			creditScoreFactor = config.CreditScoreFactorHigh
		case creditScore >= (((config.MaxCreditScore+config.MinCreditScore)/2 + config.MinCreditScore) / 2):
			creditScoreFactor = config.CreditScoreFactorMedium
		case creditScore >= (((config.MaxCreditScore+config.MinCreditScore)/2 + config.MaxCreditScore) / 2):
			creditScoreFactor = config.CreditScoreFactorLow
		default:
			creditScoreFactor = 0
		}

		eligibleLoanAmount = minEligible(config.BaseLoanAmount, maxLoanAmount) * creditScoreFactor
	}

	return
}

// minEligible returns the minimum of two float64 values
// Parameters:
// - a: float64
// - b: float64
// Returns:
// - float64 representing the minimum value
func minEligible(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
