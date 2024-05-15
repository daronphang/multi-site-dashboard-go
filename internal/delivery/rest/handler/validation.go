package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mcuadros/go-defaults"
)

// swagger:model HTTPValidationError
type HTTPValidationError struct {
	Message string `json:"message" default:"validation error"`
	Error string `json:"error"`
}

func newHTTPValidationError(c echo.Context, code int, err error) error {
	hve := &HTTPValidationError{Error: err.Error()}
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