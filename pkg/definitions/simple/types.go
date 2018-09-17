package simple

import (
	"github.com/go-playground/validator"
)

type FamilyDisease struct {
	Name     string `json:"name",structs:"name",validate:"required,alphaunicode"`
	Relative string `json:"relative",structs:"relative",validate:"required,alphaunicode"`
}

var validate *validator.Validate
