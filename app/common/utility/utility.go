package utility

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nishanthrk/aspire-lms/app/common/constants"
)

func MapToJSON(inputMap map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(inputMap)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JSONToStruct converts JSON string to a struct
func JSONToStruct(jsonData string, resultStruct interface{}) error {
	err := json.Unmarshal([]byte(jsonData), resultStruct)
	return err
}

func ValidatePlatform(platform string) string {
	switch platform {
	case constants.SystemEmployeeAPI:
		return constants.UserTypeEmployee
	case constants.SystemCustomerAPI:
		return constants.UserTypeCustomer
	default:
		return ""
	}
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// HashPassword returns a hashed password
func HashPassword(password string) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}
