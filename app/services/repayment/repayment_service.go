package repayment_service

import (
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/guregu/null"
	db "github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
	"math"
	"strings"
	"time"
)

// CalculateRepaymentSchedule generates a repayment schedule for a given loan application.
// Parameters:
// - application: pointer to the loan application model containing loan details
// Returns:
// - repayments: a slice of Repayment models representing the repayment schedule
// - err: any error that occurred during the calculation
func (s *repaymentService) CalculateRepaymentSchedule(application *models.LoanApplication) (repayments []models.Repayment, err error) {
	var paymentInterval time.Duration

	// Determine the payment interval based on the loan term unit
	switch strings.ToUpper(application.LoanTermUnit) {
	case "WEEKLY":
		paymentInterval = 7 * 24 * time.Hour
	case "MONTHLY":
		paymentInterval = time.Hour * 24 * 30
	default:
		err = fmt.Errorf("invalid frequency: %s", application.LoanTermUnit)
		return
	}

	var probableAmount float64
	var repaymentStartDate time.Time

	// Set probable amount and repayment start date based on the application status
	if application.Status == models.LoanApplicationStatusApproved {
		probableAmount = application.ApprovedAmount.Float64
		repaymentStartDate = application.ApprovedDate.Time.Add(paymentInterval)
	} else {
		probableAmount = application.LoanAmount
		repaymentStartDate = time.Now().Add(paymentInterval)
	}

	application.RepaymentStartDate = null.TimeFrom(repaymentStartDate)

	// Calculate principal amount and total repayment amount
	principalAmount := money.NewFromFloat(probableAmount, application.CurrencyCode)
	interestAmount, _ := calculateInterest(*application)
	totalRepaymentAmount := money.NewFromFloat(probableAmount+interestAmount, application.CurrencyCode)

	// Split the principal amount and total repayment amount into equal installments
	principles, err := principalAmount.Split(application.LoanTerm)
	if err != nil {
		return
	}
	emis, err := totalRepaymentAmount.Split(application.LoanTerm)
	if err != nil {
		return
	}

	// Generate repayment schedule
	for i := 0; i < application.LoanTerm; i++ {
		paymentDate := addPaymentInterval(repaymentStartDate, i, application.LoanTermUnit)
		principle := principles[i]
		emi := emis[i]
		interest, _ := emi.Subtract(principle)

		repayments = append(repayments, models.Repayment{
			RepaymentID:       uuid.New().String(),
			InstallmentNumber: i + 1,
			ApplicationID:     application.ApplicationID,
			InstallmentDate:   paymentDate,
			AmountDue:         emi.AsMajorUnits(),
			PrincipleAmount:   principle.AsMajorUnits(),
			InterestAmount:    interest.AsMajorUnits(),
			Status:            models.LoanApplicationStatusPending,
		})
	}

	return
}

// addPaymentInterval calculates the next payment date based on the term unit and iteration.
// Parameters:
// - startDate: the initial payment start date
// - iteration: the current installment number
// - termUnit: the unit of the loan term, either "MONTHLY" or "WEEKLY"
// Returns:
// - time.Time: the calculated next payment date
func addPaymentInterval(startDate time.Time, iteration int, termUnit string) time.Time {
	if termUnit == "MONTHLY" {
		return startDate.AddDate(0, iteration, 0)
	} else if termUnit == "WEEKLY" {
		return startDate.AddDate(0, 0, iteration*7)
	}
	return startDate
}

// calculateEMI calculates the Equated Monthly Installment (EMI) for a given loan.
// Parameters:
// - principal: the principal loan amount
// - annualRate: the annual interest rate in percentage
// - loanTerm: the loan term in months or weeks
// - loanTermUnit: the unit of the loan term, either "MONTHLY" or "WEEKLY"
// Returns:
// - float64: the calculated EMI amount
// - error: any error that occurred during the calculation
func calculateEMI(principal float64, annualRate float64, loanTerm int, loanTermUnit string) (float64, error) {
	var n int     // Number of installments
	var r float64 // Monthly or weekly interest rate

	// Determine the number of installments and interest rate based on the term unit
	switch strings.ToUpper(loanTermUnit) {
	case "MONTHLY":
		n = loanTerm
		r = (annualRate / 12) / 100 // Monthly interest rate
	case "WEEKLY":
		n = loanTerm
		r = (annualRate / 52) / 100 // Weekly interest rate
	default:
		return 0, fmt.Errorf("invalid loan term unit: %s", loanTermUnit)
	}

	// EMI calculation using the formula: EMI = [P * r * (1 + r)^n] / [(1 + r)^n – 1]
	emi := (principal * r * math.Pow(1+r, float64(n))) / (math.Pow(1+r, float64(n)) - 1)
	return emi, nil
}

// calculateInterest calculates the total interest for a given loan application.
// Parameters:
// - application: the loan application details
// Returns:
// - float64: the calculated total interest amount
// - error: any error that occurred during the calculation
func calculateInterest(application models.LoanApplication) (float64, error) {
	var principalAmount float64

	// Determine the principal amount based on the application status
	if application.Status == models.LoanApplicationStatusApproved {
		principalAmount = application.ApprovedAmount.Float64
	} else {
		principalAmount = application.LoanAmount
	}

	// Calculate the total repayment amount using the EMI formula
	totalRepaymentAmount, err := calculateEMI(principalAmount, application.InterestRate, application.LoanTerm, application.LoanTermUnit)
	if err != nil {
		return 0, err
	}

	// Calculate the total repayment and interest
	totalRepayment := totalRepaymentAmount * float64(application.LoanTerm)
	interest := totalRepayment - principalAmount
	return interest, nil
}

// UpdateRepayment updates the repayment details for a given request and user
// Parameters:
// - request: dto.RepaymentRequest containing the application ID and payment amount
// - user: models.User representing the user making the repayment
// Returns:
// - dto.RepaymentResponse with the repayment result
// - dto.HandleError with any error that occurred during the process
func (s *repaymentService) UpdateRepayment(request dto.RepaymentRequest, user models.User) (
	response dto.RepaymentResponse, handle dto.HandleError) {
	// Find the loan application by its primary key
	application := models.LoanApplication{}
	application, _ = application.FindByPrimaryKey(request.ApplicationID)

	// Check if the application is still pending
	if application.Status == models.LoanApplicationStatusPending {
		handle.Status = -1
		handle.Errors = fmt.Errorf("application still in: %v", models.LoanApplicationStatusPending)
		return
	}

	// Check if the application has already been paid
	if application.Status == models.LoanApplicationStatusPaid {
		handle.Status = -2
		handle.Errors = fmt.Errorf("all the installments are cleared")
		return
	}

	// Create a new payment entry
	payment := models.Payment{
		PaymentID:     uuid.New().String(),
		Amount:        request.PaymentAmount,
		ApplicationID: application.ApplicationID,
		CurrencyCode:  application.CurrencyCode,
		Status:        "RECEIVED",
	}

	var participantCondition, repaymentCondition []db.WhereCondition

	// Set conditions to find the participant
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

	// Find the participant
	participant := models.LoanApplicationParticipant{}
	participant, _ = participant.FindOneByCondition(participantCondition)

	// Check if the user has permission to update this application
	if participant.ParticipantID == "" {
		handle.Status = -1
		handle.Errors = fmt.Errorf("don't have permission to this application: %v", application.ApplicationID)
		return
	}

	// Set conditions to find pending repayments
	repaymentCondition = append(repaymentCondition, db.WhereCondition{
		Key:       models.RepaymentColumns.ApplicationID,
		Condition: "=",
		Value:     application.ApplicationID,
	})

	repaymentCondition = append(repaymentCondition, db.WhereCondition{
		Key:       models.RepaymentColumns.Status,
		Condition: "=",
		Value:     models.LoanApplicationStatusPending,
	})

	// Find all pending repayments
	repayment := models.Repayment{}
	repayments, _ := repayment.FindAllByCondition(repaymentCondition)

	// Check if any repayments are found
	if len(repayments) == 0 {
		handle.Status = -3
		handle.Errors = fmt.Errorf("failed to find repayments")
		return
	}

	// Begin a new transaction
	tx := db.MysqlDB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			handle.Status = -4
			handle.Errors = fmt.Errorf("%v", r)
			return
		}
	}()

	// Save the payment
	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		handle.Status = -5
		handle.Errors = err
		return
	}

	var repaymentPaymentLogs []models.RepaymentPaymentLog
	var updateRepayments []models.Repayment
	repaymentAmount := request.PaymentAmount

	// Iterate over each repayment
	for i := range repayments {
		repayment = repayments[i]

		// Process only pending repayments
		if repayment.Status != models.LoanApplicationStatusPaid {
			amountDueRemaining := repayment.AmountDue - repayment.AmountPaid

			// Create a new repayment payment log entry
			repaymentPaymentLog := models.RepaymentPaymentLog{
				LogID:       uuid.New().String(),
				RepaymentID: repayment.RepaymentID,
				PaymentID:   payment.PaymentID,
				Amount:      0,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			// If repayment amount is greater than or equal to the remaining amount due
			if repaymentAmount >= amountDueRemaining {
				repaymentPaymentLog.Amount = amountDueRemaining
				repayment.AmountPaid += amountDueRemaining
				repayment.Status = models.LoanApplicationStatusPaid
				repayment.PaymentDate = null.TimeFrom(time.Now())
				repayment.OutstandingBalance = null.FloatFrom(0)
				repaymentAmount -= amountDueRemaining
			} else {
				// If repayment amount is less than the remaining amount due
				repaymentPaymentLog.Amount = repaymentAmount
				repayment.AmountPaid += repaymentAmount
				repayment.OutstandingBalance = null.FloatFrom(repayment.AmountDue - repayment.AmountPaid)
				repaymentAmount = 0
			}

			// Add the updated repayment and log to the respective slices
			updateRepayments = append(updateRepayments, repayment)
			repaymentPaymentLogs = append(repaymentPaymentLogs, repaymentPaymentLog)

			// Break the loop if the repayment amount is exhausted
			if repaymentAmount == 0 {
				break
			}
		}
	}

	// Save the repayment payment logs
	if err := tx.Save(&repaymentPaymentLogs).Error; err != nil {
		tx.Rollback()
		handle.Status = -6
		handle.Errors = err
		return
	}

	// Save the updated repayments
	if err := tx.Save(&updateRepayments).Error; err != nil {
		tx.Rollback()
		handle.Status = -7
		handle.Errors = err
		return
	}

	// Check if all repayments are paid
	allPaid := true
	for _, _repayment := range updateRepayments {
		if _repayment.Status != models.LoanApplicationStatusPaid {
			allPaid = false
			break
		}
	}

	// Update the application status if all repayments are paid
	if allPaid {
		application.Status = models.LoanApplicationStatusPaid
		if err := tx.Save(&application).Error; err != nil {
			tx.Rollback()
			handle.Status = -8
			handle.Errors = err
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		handle.Status = -9
		handle.Errors = err
		return
	}

	response.Status = 1
	response.Message = "Payment processed successfully"
	response.Data.PaymentId = payment.PaymentID
	return
}
