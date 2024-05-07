package handler

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

// swagger:model HTTPValidationError
type HTTPValidationError struct {
	Message string `json:"message" default:"Invalid payload"`
	Errors []string `json:"errors"`
}

func newHTTPValidationError(c echo.Context, code int, err error) error {
	errorMsgs := strings.Split(err.Error(), ";")
	hve := &HTTPValidationError{Errors: errorMsgs}
	defaults.SetDefaults(hve)
	return c.JSON(code, hve)
}

// v is an argument that is a pointer to a value of the type that implements 
// the interface you want to validate with.
// Errors are handled in the routes.
func bindAndValidateRequestBody(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		return err
	}
	if err := c.Validate(v); err != nil {
		return err
	}
	return nil
}