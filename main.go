package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	a "github.com/openhealthalgorithms/service/actions"
	c "github.com/openhealthalgorithms/service/config"
	d "github.com/openhealthalgorithms/service/database"
	m "github.com/openhealthalgorithms/service/models"
	t "github.com/openhealthalgorithms/service/tools"
)

var (
	dbFile = filepath.Join(t.GetCurrentDirectory(), "logs.db")

	currentSettings c.Settings

	sqlite *d.SqliteDb
)

func init() {
	currentSettings = c.CurrentSettings()

	dbFile = currentSettings.LogFile
}

func main() {
	var err error

	// Check if the DB file exists
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		f, err := os.Create(dbFile)
		if err != nil {
			fmt.Println("Error:", err)
		}
		f.Close()
	}

	sqlite, err = d.InitDb(dbFile)
	if err != nil {
		fmt.Printf("Error in DB: %v\n", err)
		os.Exit(1)
	}

	err = sqlite.Migrate()
	if err != nil {
		fmt.Printf("Error in DB: %v\n", err)
		os.Exit(1)
	}
	defer sqlite.Closer()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(a.ServerHeader())
	e.Use(a.CurrentConfig())

	e.Validator = &m.CustomValidator{Validator: validator.New()}

	// Routes
	e.GET("/", a.HomeHandler)
	e.GET("/api", a.HomeHandler)
	e.GET("/api/version", a.VersionHandler)
	e.POST("/api/algorithm", a.AlgorithmHandler)

	// Start server
	e.Logger.Fatal(e.Start(":" + currentSettings.Port))
}
