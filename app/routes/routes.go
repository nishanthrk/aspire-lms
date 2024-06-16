package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes the application routes
// Parameters:
// - app: *fiber.App, the Fiber application instance
func SetupRoutes(app *fiber.App) {

	// Define the root route with a welcome message
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to aspire-lms",
		})
	})

	// Setup version 1 API routes
	SetupRoutesV1(app)
}
