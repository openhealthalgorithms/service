package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/openhealthalgorithms/service/pkg/logger"
	"github.com/openhealthalgorithms/service/pkg/service"
	"github.com/openhealthalgorithms/service/pkg/tools"
)

var (
	appName    = "ohas"
	appVersion = "v0.1"
	appCommit  = "0000000"

	serviceSrv *http.Server
	logEntry   *logrus.Entry
	pidFile    = filepath.Join(os.TempDir(), "ohad.pid")
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

	// Check if OS is supported.
	agentOS, err := tools.CheckOS()
	if err != nil {
		tools.FallbackLogger(err)
		return
	}

	// Check os Arch is supported.
	_, err = tools.CheckArch()
	if err != nil {
		tools.FallbackLogger(err)
		return
	}

	// Create a config for logger.
	logConfig := logger.Config{
		App:         appName,
		Version:     appVersion,
		Commit:      appCommit,
		Destination: "stdout",
		Level:       "info",
	}

	if agentOS == "windows" {
		logConfig.DisableColors = true
	}

	// Create a logger along with closer if its Destination is file.
	log, logCloser, err := logger.New(logConfig)
	if err != nil {
		tools.FallbackLogger(err)
		return
	}

	if logCloser != nil {
		defer logCloser()
	}

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
			Usage: "Debug mode makes output more verbose",
		},
	}

	logEntry = logger.Entry(logConfig, log)

	// The list of commands.
	app.Commands = []cli.Command{
		// Run command.
		cli.Command{
			Name:    "start",
			Aliases: []string{"start"},
			Usage:   "start the server",
			Action: func(c *cli.Context) error {
				if _, err := os.Stat(pidFile); err == nil {
					logEntry.Debugf("Already running or %s file exist.", pidFile)
					return nil
				}

				cmd := exec.Command(os.Args[0], "main")
				cmd.Start()
				logEntry.Debugf("Service process ID is : ", cmd.Process.Pid)
				savePID(cmd.Process.Pid)
				return nil
			},
		},
		cli.Command{
			Name:    "stop",
			Aliases: []string{"stop"},
			Usage:   "stop the server",
			Action: func(c *cli.Context) error {
				return stopServer(logEntry)
			},
		},
		cli.Command{
			Name:   "main",
			Action: func(c *cli.Context) error {
				return startServer(logEntry)
			},
		},
	}

	// Run the app
	// The error is handled by the deferred function
	err = app.Run(os.Args)
	if err != nil && err != context.Canceled {
		log.Error(err)
		return
	}
}

func startServer(logEntry *logrus.Entry) error {
	logEntry.Debug("server starting")

	// Make arrangement to remove PID file upon receiving the SIGTERM from kill command
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		signalType := <-ch
		signal.Stop(ch)
		logEntry.Debug("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		logEntry.Debugln("Received signal type : ", signalType)

		// remove PID file
		os.Remove(pidFile)

		os.Exit(0)
	}()

	srv := service.NewService()
	srv.StartHttpServer()

	logEntry.Debug("server started")

	return nil
}

func stopServer(logEntry *logrus.Entry) error {
	if _, err := os.Stat(pidFile); err == nil {
		data, err := ioutil.ReadFile(pidFile)
		if err != nil {
			logEntry.Debugln("Not running")
			return err
		}
		ProcessID, err := strconv.Atoi(string(data))

		if err != nil {
			logEntry.Debugln("Unable to read and parse process id found in ", pidFile)
			return err
		}

		process, err := os.FindProcess(ProcessID)

		if err != nil {
			logEntry.Debugf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
			return err
		}
		// remove PID file
		os.Remove(pidFile)

		logEntry.Debugf("Killing process ID [%v] now.\n", ProcessID)
		// kill process and exit immediately
		err = process.Kill()

		if err != nil {
			logEntry.Debugf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
			return err
		} else {
			logEntry.Debugf("Killed process ID [%v]\n", ProcessID)
			return err
		}
	} else {
		logEntry.Debugln("Not running.")
		return err
	}
}

func savePID(pid int) error {
	file, err := os.Create(pidFile)
	if err != nil {
		logEntry.Errorf("Unable to create pid file : %v\n", err)
		return err
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		logEntry.Errorf("Unable to create pid file : %v\n", err)
		return err
	}

	file.Sync() // flush to disk

	return nil
}
