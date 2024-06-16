package loan_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nishanthrk/aspire-lms/app/common/validator"
	"github.com/nishanthrk/aspire-lms/app/dto"
	loanService "github.com/nishanthrk/aspire-lms/app/services/loan"
	repaymentService "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
	"net/http"
)

// CreateLoanApplication handles the creation of a loan application
// Parameters:
// - c: *fiber.Ctx representing the request context
// - loanService: loanService.LoanService for handling loan-related operations
// - userService: userService.UserService for handling user-related operations
// - repaymentService: repaymentService.RepaymentService for handling repayment-related operations
// Returns:
// - An error if there was an issue during the process; otherwise, it returns a JSON response with the created loan application details
func CreateLoanApplication(c *fiber.Ctx, loanService loanService.LoanService, userService userService.UserService, repaymentService repaymentService.RepaymentService) error {
	// Parse and validate the request body into params
	params := dto.ApplicationCreateRequest{}
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		// Return a 422 Unprocessable Entity status with the validation error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	// Call the loanService to create the loan application
	response, handle := loanService.CreateLoanApplication(params, userService, repaymentService)
	if handle.Status < 0 {
		// Return a 422 Unprocessable Entity status with the service error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": handle.Status,
			"error":  handle.Errors.Error(),
		})
	}

	// Return a 200 OK status with the created loan application details
	return c.Status(http.StatusOK).JSON(response)
}

// GetApplicationList retrieves the list of loan applications in which the authenticated user is a participant.
// Parameters:
// - c: *fiber.Ctx representing the Fiber context
// - loanService: loanService.LoanService for handling loan-related operations
// - userService: userService.UserService for handling user-related operations
// Returns:
// - An error if any occurs during the process, otherwise returns a JSON response with the application list.

func GetApplicationList(c *fiber.Ctx, loanService loanService.LoanService, userService userService.UserService) error {
	// Retrieve the authenticated user object from the context
	user := userService.GetUserObject(c)

	// Use the loan service to get the list of loan applications where the user is a participant
	response := loanService.GetParticipantApplications(user)

	// Return the response as JSON with an HTTP 200 OK status
	return c.Status(http.StatusOK).JSON(response)
}

// ApproveLoanApplication handles the approval of a loan application
// Parameters:
// - c: *fiber.Ctx representing the request context
// - loanService: loanService.LoanService for handling loan-related operations
// - userService: userService.UserService for handling user-related operations
// - repaymentService: repaymentService.RepaymentService for handling repayment-related operations
// Returns:
// - An error if there was an issue during the process; otherwise, it returns a JSON response with the approval details
func ApproveLoanApplication(c *fiber.Ctx, loanService loanService.LoanService, userService userService.UserService, repaymentService repaymentService.RepaymentService) error {
	// Initialize the params object and set the ApplicationID from the URL parameter
	params := dto.ApplicationApproveRequest{}
	params.ApplicationID = c.Params("applicationId")

	// Parse and validate the request body into params
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		// Return a 422 Unprocessable Entity status with the validation error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	// Retrieve the user object from the request context
	user := userService.GetUserObject(c)

	// Call the loanService to approve the loan application
	response, handle := loanService.ApproveLoanApplication(params, user, repaymentService)
	if handle.Status < 0 {
		// Return a 422 Unprocessable Entity status with the service error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": handle.Status,
			"error":  handle.Errors.Error(),
		})
	}

	// Return a 200 OK status with the approval details
	return c.Status(http.StatusOK).JSON(response)
}

// GetLoanApplication retrieves the details of a loan application
// Parameters:
// - c: *fiber.Ctx representing the request context
// - loanService: loanService.LoanService for handling loan-related operations
// - userService: userService.UserService for handling user-related operations
// Returns:
// - An error if there was an issue during the process; otherwise, it returns a JSON response with the loan application details
func GetLoanApplication(c *fiber.Ctx, loanService loanService.LoanService, userService userService.UserService) error {
	// Retrieve the user object from the request context
	user := userService.GetUserObject(c)

	// Call the loanService to get the loan application details
	response, handle := loanService.GetLoanApplication(c.Params("applicationId"), user)
	if handle.Status < 0 {
		// Return a 422 Unprocessable Entity status with the service error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": handle.Status,
			"error":  handle.Errors.Error(),
		})
	}

	// Return a 200 OK status with the loan application details
	return c.Status(http.StatusOK).JSON(response)
}
