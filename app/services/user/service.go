package user_service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
)

// UserService defines the interface for user-related operations
type UserService interface {
	// ValidateCredentials validates the user's login credentials
	ValidateCredentials(request dto.LoginRequest) (response dto.LoginResponse, handle dto.HandleError)

	// GetUserObject retrieves the user object from the JWT token present in the request context
	GetUserObject(c *fiber.Ctx) models.User

	// GenerateAuth generates authentication tokens for the user
	GenerateAuth(users models.User) (dto.LoginResponse, error)

	// AllocateEmployeeForProcess allocates an employee for processing the loan application
	AllocateEmployeeForProcess(application models.LoanApplication) (models.LoanApplicationParticipant, error)

	// GetApplicationUser retrieves the user object based on the provided details
	GetApplicationUser(object dto.UserObject) models.User
}

// userService is an implementation of UserService
type userService struct{}

// NewUserService returns a new instance of UserService
func NewUserService() UserService {
	return &userService{}
}
