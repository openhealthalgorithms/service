package main

import (
	"context"
	"os"

	"github.com/openhealthalgorithms/service/pkg/tools"
	"github.com/urfave/cli"
)

var (
	appName    = "oha-service"
	appVersion = "v0.1"
	appCommit  = "0000000"
)

func main() {
	// To be able to exit properly in case of error.
	var err error
	// Deferred calls are run in LIFO order.
	// os.Exit does not run any other deferred calls.
	// This way allows us to exit the app releasing resources properly.
	defer func() {
		if err != nil && err != context.Canceled {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// Create an instance of the cli.App to process commands and options (flags).
	// For each mode we have a command.
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "Provides an API service for Open Health Algorithms"
	app.Version = appVersion
	// Create Metadata to carry commit version
	app.Metadata = make(map[string]interface{})
	app.Metadata["commit"] = appCommit
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Open Health Algorithms",
			Email: "contact@openhealthalgorithms.org",
		},
	}

	// The list of global flags.
	app.Flags = []cli.Flag{
		// Debug mode sets level for logger to debug.
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug mode makes output more verbose. Default - off",
		},
		// Local mode makes agent send data to localhost.
		cli.BoolFlag{
			Name:  "grace",
			Usage: "Grace mode enables graceful shutdown. Default - off",
		},
	}

	// The list of commands.
	app.Commands = []cli.Command{
		// Run command.
		cli.Command{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run discovery in cli mode",
			Action: func(c *cli.Context) error {
				// Check license-key.
				if err := checkLicenseKey(c); err != nil {
					cli.ShowAppHelp(c)
					tools.FallbackLogger(err)
					return err
				}
				// Validate server-url.
				if err := validateServerURL(c); err != nil {
					cli.ShowAppHelp(c)
					tools.FallbackLogger(err)
					return err
				}

				return setupAndRun(c)
			},
		},
		// Daemon command.
		cli.Command{
			Name:    "daemon",
			Aliases: []string{"d"},
			Usage:   "run discovery in daemon mode",
			Action: func(c *cli.Context) error {
				// Check license-key.
				if err := checkLicenseKey(c); err != nil {
					cli.ShowAppHelp(c)
					tools.FallbackLogger(err)
					return err
				}
				// Validate server-url.
				if err := validateServerURL(c); err != nil {
					cli.ShowAppHelp(c)
					tools.FallbackLogger(err)
					return err
				}

				return setupAndRun(c)
			},
		},
	}

	// Run the app
	// The error is handled by the deferred function
	err = app.Run(os.Args)
}
