package user_service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/nishanthrk/aspire-lms/app/common/constants"
	"github.com/nishanthrk/aspire-lms/app/common/utility"
	cfg "github.com/nishanthrk/aspire-lms/app/configs"
	db "github.com/nishanthrk/aspire-lms/app/database"
	"github.com/nishanthrk/aspire-lms/app/dto"
	"github.com/nishanthrk/aspire-lms/app/models"
	"time"
)

// AccessClaims represents access token JWT claims
type AccessClaims struct {
	Username string `json:"username"`
	UserType string `json:"user_type"`
	UserId   string `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateCredentials validates the user's credentials for login
// Parameters:
// - request: dto.LoginRequest containing the login details (platform and identifier)
// Returns:
// - dto.LoginResponse with the authentication tokens
// - dto.HandleError with any error that occurred during the process
func (s *userService) ValidateCredentials(request dto.LoginRequest) (response dto.LoginResponse, handle dto.HandleError) {
	userModel := models.User{}
	var andCondition, orCondition []db.WhereCondition

	// Add platform validation condition
	andCondition = append(andCondition, db.WhereCondition{
		Key:       models.UsersColumns.UserType,
		Condition: "=",
		Value:     utility.ValidatePlatform(request.Platform),
	})

	// Add email identifier condition
	orCondition = append(orCondition, db.WhereCondition{
		Key:       models.UsersColumns.UserEmail,
		Condition: "=",
		Value:     request.Identifier,
	})

	// Add mobile number identifier condition
	orCondition = append(orCondition, db.WhereCondition{
		Key:       models.UsersColumns.MobileNumber,
		Condition: "=",
		Value:     request.Identifier,
	})

	// Find user by conditions
	userModel, _ = userModel.FindOneByCondition(&andCondition, &orCondition)
	if userModel.UserID == "" {
		handle.Status = -2
		handle.Errors = fmt.Errorf("user details not found")
		return
	}

	// Check if the password matches
	if userModel.UserPassword != utility.HashPassword(request.Password) {
		handle.Status = -3
		handle.Errors = fmt.Errorf("user details not found")
		return
	}

	// Generate authentication tokens
	response, err := s.GenerateAuth(userModel)
	if err != nil {
		handle.Status = -4
		handle.Errors = fmt.Errorf("user details not found")
		return
	}

	return
}

// GenerateAuth generates JWT access and refresh tokens for the authenticated user
// Parameters:
// - users: models.User representing the authenticated user
// Returns:
// - dto.LoginResponse with the authentication tokens and user details
// - error if any error occurred during token generation
func (s *userService) GenerateAuth(users models.User) (detail dto.LoginResponse, err error) {
	// Set the access token expiration time to 14 hours
	expireTime := time.Now().Add(time.Hour * 14)

	// Create access claims with user information and token metadata
	accessClaims := AccessClaims{
		users.UserName,
		users.UserType,
		utility.ToString(users.UserID),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    cfg.GetConfig().Tenant,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate the access token with the specified signing method and secret key
	accessClaimWithSecret := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessClaimWithSecret.SignedString([]byte(cfg.GetConfig().JWTAccessSecret))
	if err != nil {
		return
	}

	// Create refresh claims with user information and token metadata
	refreshClaims := AccessClaims{
		users.UserName,
		users.UserType,
		utility.ToString(users.UserID),
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    cfg.GetConfig().Tenant,
			NotBefore: jwt.NewNumericDate(expireTime),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}

	// Generate the refresh token with the specified signing method and secret key
	refreshClaimWithSecret := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshClaimWithSecret.SignedString([]byte(cfg.GetConfig().JWTRefreshSecret))
	if err != nil {
		return
	}

	// Return the generated tokens and user details in the response
	return dto.LoginResponse{
		Status: 1,
		Data: dto.TokenDetail{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			UserID:       utility.ToString(users.UserID),
			UserType:     users.UserType,
			TokenExpires: expireTime.Unix(),
		},
	}, nil
}

// GetUserObject retrieves the user object from the JWT token in the request context
// Parameters:
// - c: *fiber.Ctx representing the Fiber request context
// Returns:
// - models.User representing the user extracted from the JWT token
func (s *userService) GetUserObject(c *fiber.Ctx) (userModel models.User) {
	defer func() {
		if r := recover(); r != nil {
			// Handle any panic that might occur
		}
	}()

	// Initialize an empty AccessClaims struct to store token claims
	access := AccessClaims{}

	// Retrieve the JWT token from the request context
	user := c.Locals("user").(*jwt.Token)
	// Extract claims from the JWT token
	claims := user.Claims.(jwt.MapClaims)

	// Convert the claims map to JSON
	jsonData, err := utility.MapToJSON(claims)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return
	}

	// Convert the JSON data to AccessClaims struct
	err = utility.JSONToStruct(jsonData, &access)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		return
	}

	// Find the user by primary key (user ID) extracted from the claims
	userModel, _ = userModel.FindByPrimaryKey(access.UserId)
	return
}

// AllocateEmployeeForProcess assigns the least loaded employee to the loan application
// Parameters:
// - application: models.LoanApplication representing the loan application
// Returns:
// - models.LoanApplicationParticipant representing the assigned employee as a participant
// - error if an employee could not be found or assigned
func (s *userService) AllocateEmployeeForProcess(application models.LoanApplication) (participant models.LoanApplicationParticipant, err error) {
	user := models.User{}

	// Find the least loaded employee
	employeeId, _ := user.FindLeastLoadedEmployee()

	// Retrieve the employee details by primary key
	user, _ = user.FindByPrimaryKey(employeeId)
	if user.UserID == "" {
		err = fmt.Errorf("employee not configured")
		return
	}

	// Create a new loan application participant for the employee
	participant = models.LoanApplicationParticipant{
		ParticipantID:   uuid.New().String(),
		ApplicationID:   application.ApplicationID,
		ParticipantType: constants.UserTypeEmployee,
		UserID:          user.UserID,
	}

	return
}

// GetApplicationUser retrieves or creates a user based on the provided UserObject
// Parameters:
// - object: dto.UserObject containing user details
// Returns:
// - models.User representing the user found or created
func (s *userService) GetApplicationUser(object dto.UserObject) (userModel models.User) {
	var andCondition, orCondition []db.WhereCondition

	// Add condition to match the user type as 'Customer'
	andCondition = append(andCondition, db.WhereCondition{
		Key:       models.UsersColumns.UserType,
		Condition: "=",
		Value:     constants.UserTypeCustomer,
	})

	// Add conditions to match either user email or mobile number
	orCondition = append(orCondition, db.WhereCondition{
		Key:       models.UsersColumns.UserEmail,
		Condition: "=",
		Value:     object.UserEmail,
	})

	orCondition = append(orCondition, db.WhereCondition{
		Key:       models.UsersColumns.MobileNumber,
		Condition: "=",
		Value:     object.MobileNumber,
	})

	// Find the user by the specified conditions
	userModel, _ = userModel.FindOneByCondition(&andCondition, &orCondition)
	return
}
