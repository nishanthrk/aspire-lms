package routes

import (
	"github.com/gofiber/fiber/v2"
	loanController "github.com/nishanthrk/aspire-lms/app/controllers/v1/loan"
	repaymentController "github.com/nishanthrk/aspire-lms/app/controllers/v1/repayment"
	userController "github.com/nishanthrk/aspire-lms/app/controllers/v1/user"
	"github.com/nishanthrk/aspire-lms/app/logger"
	"github.com/nishanthrk/aspire-lms/app/middlewares"
	loanService "github.com/nishanthrk/aspire-lms/app/services/loan"
	repaymentService "github.com/nishanthrk/aspire-lms/app/services/repayment"
	userService "github.com/nishanthrk/aspire-lms/app/services/user"
)

// SetupRoutesV1 sets up the version 1 routes for the aspire-lms API
// Parameters:
// - app: *fiber.App representing the Fiber application instance
func SetupRoutesV1(app *fiber.App) {
	// Create a new route group for version 1 endpoints with optional authentication middleware
	v1 := app.Group("/v1", middlewares.OptionalAuth())

	// Define a health check endpoint
	v1.Get("/", func(c *fiber.Ctx) error {
		message := "Welcome to go-fiber aspire-lms with 12 factor"
		logger.Sugar.Info(message)
		return c.JSON(fiber.Map{
			"status":  1,
			"message": message,
		})
	})

	// Initialize the service instances
	loanSvc := loanService.NewLoanService()
	repaymentSvc := repaymentService.NewRepaymentService()
	userSvc := userService.NewUserService()

	// Define the user-related routes
	userRoute := v1.Group("user")
	userRoute.Post("/auth", func(c *fiber.Ctx) error {
		return userController.Login(c, userSvc)
	})

	// Define the application-related routes
	applicationRoute := v1.Group("/application")

	// Route for creating a new loan application
	applicationRoute.Post("/", func(c *fiber.Ctx) error {
		return loanController.CreateLoanApplication(c, loanSvc, userSvc, repaymentSvc)
	})

	// Define a restricted route group for application-related operations that require authentication
	restrictedApplicationRoute := applicationRoute.Group("/", middlewares.RequireLoggedIn())

	restrictedApplicationRoute.Get("/", func(c *fiber.Ctx) error {
		return loanController.GetApplicationList(c, loanSvc, userSvc)
	})

	// Route for getting loan application details
	restrictedApplicationRoute.Get("/:applicationId", func(c *fiber.Ctx) error {
		return loanController.GetLoanApplication(c, loanSvc, userSvc)
	})

	// Route for making a repayment
	restrictedApplicationRoute.Post("/:applicationId/repayment", func(c *fiber.Ctx) error {
		return repaymentController.PayRepayment(c, repaymentSvc, userSvc)
	})

	adminApplicationRoute := restrictedApplicationRoute.Group("/:applicationId/approve",
		middlewares.RequireAdmin)

	// Route for approving a loan application
	adminApplicationRoute.Post("/", func(c *fiber.Ctx) error {
		return loanController.ApproveLoanApplication(c, loanSvc, userSvc, repaymentSvc)
	})
}
