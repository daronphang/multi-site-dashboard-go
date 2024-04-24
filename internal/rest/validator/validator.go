package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

type HTTPValidationError struct {
	Message string `json:"message" default:"Invalid payload"`
	Errors []string `json:"errors"`
}

type CustomValidator struct {
    Validator *validator.Validate
}

func NewHTTPValidationError(c echo.Context, code int, err error) error {
	errorMsgs := strings.Split(err.Error(), ";")
	hve := &HTTPValidationError{Errors: errorMsgs}
	defaults.SetDefaults(hve)
	return c.JSON(code, hve)
}

func ValidatePayload(c echo.Context, v interface{}) error {
	// v is an argument that is a pointer to a value of the type that implements 
	// the interface you want to validate with.
	// Errors are handled in the routes.
	if err := c.Bind(v); err != nil {
		return err
	}
	if err := c.Validate(v); err != nil {
		return err
	}
	return nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Returns an array of type FieldError
		var ve validator.ValidationErrors
		var errorMsg string
		if !errors.As(err, &ve) {
			errorMsg = err.Error()
		} else {
			errorMsg = flattenAndTranslateErrors(&ve)
		}
	  return errors.New(errorMsg)
	}
	return nil
  }

// Substitutes tag specified in struct for error message if thrown
func msgForTag(fe validator.FieldError) string {
	var msg string
	switch fe.Tag() {
	case "required":
		msg = "field is required"
	case "email":
		msg = "invalid email"
	case "boolean":
		msg = "field must be boolean"
	case "number":
		msg = "field must be number"
	case "numeric":
		msg = "field must be numeric"
	case "url":
		msg = "field must be url"
	case "len":
		msg = fmt.Sprintf("field must have length of %s", fe.Param())
	case "oneof":
		msg = fmt.Sprintf("field must be one of %s", fe.Param())
	case "gt":
		msg = fmt.Sprintf("field must be greater than %s", fe.Param())
	case "lt":
		msg = fmt.Sprintf("field must be less than %s", fe.Param())
	case "max":
		msg = fmt.Sprintf("maximum allowed is %s", fe.Param())
	case "min":
		msg = fmt.Sprintf("minimum allowed is %s", fe.Param())
	case "eq":
		msg = fmt.Sprintf("field must be equal to %s", fe.Param())
	case "ne":
		msg = fmt.Sprintf("field must not be equal to %s", fe.Param())
	default:
		msg = fe.Error()
	}
	return msg
}

// func translateErrors(ve *validator.ValidationErrors) []string {
// 	output := make([]string, len(*ve))
// 	for _, fe := range *ve {
// 		output = append(output, fmt.Sprintf("%s: %s", fe.Field(), msgForTag(fe)))
// 	}
// 	return output
// }

func flattenAndTranslateErrors(ve *validator.ValidationErrors) string {
	b := new(bytes.Buffer)
	for _, fe := range *ve {
		fmt.Fprintf(b, "%s: %s;", fe.Field(), msgForTag(fe))
	}
	s := b.String()
	return s[:len(s) - 1]
}
