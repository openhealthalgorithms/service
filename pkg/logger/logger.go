// Package logger provides configured logger.
package logger

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openhealthalgorithms/service/pkg/tools"
)

var (
	// ErrPathEmpty is returned when Path is empty for Destination file.
	ErrPathEmpty = errors.New("path is empty")
)

// Config holds options to create a logger.
type Config struct {
	App           string
	Version       string
	Commit        string
	Destination   string
	Level         string
	Path          string
	Platform      string
	MaxFileSize   int64
	DisableColors bool
	Extra         map[string]string
}

// New returns a new logger based on config.
func New(c Config) (*logrus.Logger, func(), error) {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: c.DisableColors,
	}

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     setLevel(c.Level),
	}

	var logFile *os.File
	var logCloser func()
	var err error

	if c.Destination == "file" || c.Destination == "combined" {
		if c.Path == "" {
			return nil, nil, errors.Wrapf(ErrPathEmpty, "create logger")
		}

		fInfo, err1 := os.Stat(c.Path)
		if err1 != nil {
			if !os.IsNotExist(err1) {
				return nil, nil, errors.Wrapf(err1, "get logfile info")
			}
		}

		if fInfo == nil {
			logFile, logCloser, err = tools.CreateOrWriteFile(c.Path)
		} else {
			if fInfo.Size() >= c.MaxFileSize {
				logFile, logCloser, err = tools.CreateOrWriteFile(c.Path)
			} else {
				logFile, logCloser, err = tools.CreateOrAppendFile(c.Path)
			}
		}

		if err != nil {
			return nil, nil, err
		}

		switch c.Destination {
		case "file":
			logger.Out = logFile
		case "combined":
			multiOut := io.MultiWriter(os.Stderr, logFile)
			logger.Out = multiOut
		}
	}

	return logger, logCloser, nil
}

// Entry sets additional fields for logger and returns an Entry.
//
// After getting an Entry seetings of the logger should not be changed.
func Entry(c Config, log *logrus.Logger) *logrus.Entry {
	entry := log.WithFields(logrus.Fields{
		"app": c.App,
		"ver": c.Version,
	})

	if c.Commit != "" {
		entry = entry.WithField("commit", c.Commit)
	}

	if c.Platform != "" {
		entry = entry.WithField("platform", c.Platform)
	}

	if len(c.Extra) > 0 {
		for k, v := range c.Extra {
			entry = entry.WithField(k, v)
		}
	}

	return entry
}

// setLevel sets logrus.Level based on given string value.
func setLevel(v string) logrus.Level {
	var l logrus.Level

	switch v {
	case "debug":
		l = logrus.DebugLevel
	case "info":
		l = logrus.InfoLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	case "fatal":
		l = logrus.FatalLevel
	case "panic":
		l = logrus.PanicLevel
	default:
		l = logrus.DebugLevel
	}

	return l
}
