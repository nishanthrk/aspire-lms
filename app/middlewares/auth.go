package middlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nishanthrk/aspire-lms/app/common/constants"
	cfg "github.com/nishanthrk/aspire-lms/app/configs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// RequireLoggedIn ensures access only to login users by checking for token presence and validity
func RequireLoggedIn() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(cfg.GetConfig().JWTAccessSecret),
		ErrorHandler: jwtError,
	})
}

func OptionalAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.GetConfig().JWTAccessSecret), // Replace with your actual secret key
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if c.Get("X-Platform") == "" {
				var errorList []*fiber.Error
				errorList = append(
					errorList,
					&fiber.Error{
						Code:    fiber.StatusBadRequest,
						Message: "Missing or X-Platform",
					},
				)
				return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errorList})
			}
			// Handle authentication errors (optional)
			fmt.Println("Authentication error:", err)
			return c.Next()
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if c.Get("X-Platform") == "" {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Missing or X-Platform",
			},
		)
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errorList})
	}

	if err.Error() == "Missing or malformed JWT" {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing or Malformed Authentication Token",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
	}

	var errorList []*fiber.Error
	errorList = append(
		errorList,
		&fiber.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid or Expired Authentication Token",
		},
	)
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
}

// RequireAdmin Ensures A route Can Only Be Accessed by an Admin user
// This function can be extended to handle different roles
func RequireAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := claims["user_type"].(string)

	var errorList []*fiber.Error

	if role != constants.UserTypeEmployee {
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "You're Not Authorized",
			},
		)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"errors": errorList})
	}
	return c.Next()
}
