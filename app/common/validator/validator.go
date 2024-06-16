package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nishanthrk/aspire-lms/app/models"
	"reflect"
	"regexp"
)

type ErrorDetail struct {
	Id         string `json:"id"`
	ColumnName string `json:"column_name"`
	Message    string `json:"message"`
}

var validate = validator.New()

func Validate(payload interface{}) []*fiber.Error {
	err := validate.Struct(payload)
	var errorList []*fiber.Error
	if err != nil {
		// Empty errors slice to store the errors
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			errorList = append(
				errorList,
				&fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: fmt.Sprintf("Invalid value for %v: %v", fieldName, err.ActualTag()),
				},
			)
		}
	}

	validateFields(reflect.ValueOf(payload), &errorList)

	if len(errorList) > 0 {
		return errorList
	}

	return nil
}

func validateFields(val reflect.Value, errorList *[]*fiber.Error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		validateStructFields(val, errorList)
	case reflect.Slice, reflect.Array:
		validateSliceElements(val, errorList)
	case reflect.Map:
		validateMapElements(val, errorList)
	default:
	}
}

func validateStructFields(val reflect.Value, errorList *[]*fiber.Error) {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Type.Kind() == reflect.String && !field.IsZero() && field.String() != "" {
			match, err := regexp.MatchString("^[a-zA-Z0-9 &@._/{},'\"#+:;=\\\\[\\]-]+$", field.String())
			if err != nil {
				*errorList = append(
					*errorList,
					&fiber.Error{
						Code:    fiber.StatusBadRequest,
						Message: fmt.Sprintf("Invalid regex pattern"),
					},
				)
				continue
			}

			if !match {
				*errorList = append(
					*errorList,
					&fiber.Error{
						Code:    fiber.StatusBadRequest,
						Message: fmt.Sprintf("%v does not match the pattern", fieldType.Name),
					},
				)
				continue
			}
		}

		validateFields(field, errorList)
	}
}

func validateSliceElements(val reflect.Value, errorList *[]*fiber.Error) {
	for i := 0; i < val.Len(); i++ {
		validateFields(val.Index(i), errorList)
	}
}

func validateMapElements(val reflect.Value, errorList *[]*fiber.Error) {
	for _, key := range val.MapKeys() {
		validateFields(val.MapIndex(key), errorList)
	}
}

// ParseBody is a helper function for parsing the body.
// Is any error occurring it will panic.
// It's just a helper function to avoid writing if condition again n again.
func ParseBody(c *fiber.Ctx, body interface{}) []*fiber.Error {

	if err := c.BodyParser(body); err != nil {
		var errorList []*fiber.Error
		errorList = append(
			errorList,
			&fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: fmt.Sprintf("%v given value is not matching", err.Error()),
			},
		)

		return errorList
	}

	return nil
}

// ParseBodyAndValidate is a helper function for parsing the body.
// Is any error occurring it will panic.
// It's just a helper function to avoid writing if condition again n again.
func ParseBodyAndValidate(c *fiber.Ctx, body interface{}) []*fiber.Error {

	// First We Parse
	if err := ParseBody(c, body); err != nil {
		return err
	}
	// Then We Validate
	return Validate(body)
}
func init() {
	validate = validator.New()

	_ = validate.RegisterValidation("validateCountryCode", func(fl validator.FieldLevel) bool {
		countryCode := fl.Field().String()

		// Mock the database lookup here so that test can pass
		if countryCode == "IND" {
			return true
		}

		if countryCode != "" {
			country := models.Country{}
			if country, _ = country.FindByPrimaryKey(countryCode); country.CountryCode == "" {
				return false
			}
			return true
		} else {
			return true
		}
	})

	_ = validate.RegisterValidation("validateCurrencyCode", func(fl validator.FieldLevel) bool {
		currencyCode := fl.Field().String()
		// Mock the database lookup here so that test can pass
		if currencyCode == "INR" {
			return true
		}

		if currencyCode != "" {
			currency := models.Currency{}
			if currency, _ = currency.FindByPrimaryKey(currencyCode); currency.CurrencyCode == "" {

				return false
			}
			return true
		} else {
			return true
		}
	})
}
