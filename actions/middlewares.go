package actions

import (
    "strings"

    "github.com/labstack/echo/v4"
    "github.com/pkg/errors"

    "github.com/openhealthalgorithms/service/config"
    "github.com/openhealthalgorithms/service/pkg"
)

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Response().Header().Set(echo.HeaderServer, "ohas/"+strings.TrimPrefix(pkg.GetVersion(), "v"))
            return next(c)
        }
    }
}

// CurrentConfig middleware adds the list of current configuration to context
func CurrentConfig() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            currentSettings := config.CurrentSettings()
            if currentSettings.CloudEnable {
                projectName := strings.Replace(c.Request().URL.Path, "/api/algorithm", "", 1)
                if len(projectName) > 1 {
                    projectName = projectName[1:]

                    // check for authorization
                    authorizationToken := c.Request().Header.Get("Authorization")
                    if len(authorizationToken) == 0 {
                        return errors.New("authorization token missing")
                    }
                    if !strings.HasPrefix(authorizationToken, "Bearer ") || len(authorizationToken) != 71 {
                        return errors.New("invalid token format. should be in the format of 'Bearer YOUR_TOKEN'")
                    }

                    bearerToken := strings.TrimPrefix(authorizationToken, "Bearer ")
                    if len(bearerToken) != 64 {
                        return errors.New("invalid token for the api")
                    }

                    // check for api token in the database and get the project name
                    projectForToken, err := checkAPIToken(bearerToken,
                        currentSettings.CloudDBHost,
                        currentSettings.CloudDBName,
                        currentSettings.CloudDBUser,
                        currentSettings.CloudDBPassword,
                    )
                    if err != nil {
                        return err
                    }
                    if projectForToken != projectName {
                        return errors.New("invalid token for the project")
                    }
                }
            } else {
                currentSettings.CloudEnable = false
            }
            c.Set("current_config", currentSettings)
            return next(c)
        }
    }
}
