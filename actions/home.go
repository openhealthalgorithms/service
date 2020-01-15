package actions

import (
    "net/http"

    "github.com/labstack/echo/v4"

    "github.com/openhealthalgorithms/service/pkg"
)

// HomeHandler func
func HomeHandler(c echo.Context) error {
    data := map[string]string{
        "message": "Welcome to Open Health Algorithm Service",
        "version": pkg.GetVersion(),
        "package": "ohas",
    }

    return c.JSON(http.StatusOK, data)
}

// VersionHandler func
func VersionHandler(c echo.Context) error {
    data := map[string]string{
        "version": pkg.GetVersion(),
    }
    
    return c.JSON(http.StatusOK, data)
}
