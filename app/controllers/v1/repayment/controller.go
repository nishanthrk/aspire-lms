package repayment_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nishanthrk/aspire-lms/app/common/validator"
	"github.com/nishanthrk/aspire-lms/app/dto"
	repaymentService "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
	"net/http"
)

// PayRepayment handles the repayment process for a loan application
// Parameters:
// - c: *fiber.Ctx representing the request context
// - repaymentService: repaymentService.RepaymentService for handling repayment-related operations
// - userService: userService.UserService for handling user-related operations
// Returns:
// - An error if there was an issue during the process; otherwise, it returns a JSON response with the repayment details
func PayRepayment(c *fiber.Ctx, repaymentService repaymentService.RepaymentService, userService userService.UserService) error {
	// Initialize a RepaymentRequest DTO to hold the request parameters
	params := dto.RepaymentRequest{}
	// Set the ApplicationID from the URL parameters
	params.ApplicationID = c.Params("applicationId")

	// Parse and validate the request body into the params object
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		// Return a 422 Unprocessable Entity status with the validation error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	// Retrieve the user object from the request context
	user := userService.GetUserObject(c)

	// Call the repaymentService to update the repayment details
	response, handle := repaymentService.UpdateRepayment(params, user)
	if handle.Status < 0 {
		// Return a 422 Unprocessable Entity status with the service error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": handle.Status,
			"error":  handle.Errors.Error(),
		})
	}

	// Return a 200 OK status with the repayment details
	return c.Status(http.StatusOK).JSON(response)
}
