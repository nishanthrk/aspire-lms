package user_controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nishanthrk/aspire-lms/app/common/validator"
	"github.com/nishanthrk/aspire-lms/app/dto"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
	"net/http"
)

// Login handles user authentication by validating the provided credentials
// Parameters:
// - c: *fiber.Ctx representing the request context
// - userService: userService.UserService for handling user-related operations
// Returns:
// - An error if there was an issue during the authentication process; otherwise, it returns a JSON response with the authentication details
func Login(c *fiber.Ctx, userService userService.UserService) error {
	// Initialize a LoginRequest DTO to hold the request parameters
	var params dto.LoginRequest
	// Set the Platform from the request headers
	params.Platform = c.Get("X-Platform")

	// Parse and validate the request body into the params object
	if err := validator.ParseBodyAndValidate(c, &params); err != nil {
		// Return a 422 Unprocessable Entity status with the validation error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": -1,
			"error":  err,
		})
	}

	// Call the userService to validate the provided credentials
	response, handle := userService.ValidateCredentials(params)
	if handle.Status < 0 {
		// Return a 422 Unprocessable Entity status with the service error
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"status": handle.Status,
			"error":  handle.Errors.Error(),
		})
	}

	// Return a JSON response with the authentication details
	return c.JSON(response)
}
