package models

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator type
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
