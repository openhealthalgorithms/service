package actions

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// Errs object
type Errs struct {
	Error    bool     `json:"error"`
	Messages []string `json:"messages"`
}

// ErrorResponse function
func ErrorResponse(c echo.Context, err error, code int) error {
	er := Errs{}
	er.Error = true
	er.Messages = strings.SplitN(err.Error(), "\n", -1)
	return c.JSON(code, er)
}
