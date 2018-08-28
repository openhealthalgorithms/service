// Package config provides configuration constants and variables for the agent.
//
// WARNING: This package must not import any other agent's package.
// Otherwide the circle dependency might occur.
package config

import (
	"time"
)

// Default contains default configuration settings.
var Default Settings = NewSettings()

// NewSettings returns default configuration settings.
//
// WARNING! All the configuration MUST be done by setting appropriate values here!
func NewSettings() Settings {
	agent := &Agent{
		// The number of tasks being executing simultaneously.
		maxTaskExecutors: 4,
		// An URL which is used when --local flag is passed.
		// Must not be changed.
		localServerURL: "http://127.0.0.1:8080",

		// Interval between iterations.
		loopInterval:       60 * time.Second,
		startDelayDuration: 2 * time.Minute,

		// Default Timeout for shell commands.
		// Must not be changed.
		cmdexecTimeout: 10 * time.Second,
	}

	api := &API{
		Methods: map[string]string{
			EPCompanyConfiguration: "/api/v2.1/company-configuration",
			EPConnections:          "/api/v2.1/connections",
			EPFile:                 "/api/v2/file",
			EPFileSearch:           "/api/v2/file-search",
			EPFileSystem:           "/api/v2/filesystem-snapshot",
			EPLog:                  "/api/v2/log",
			EPRegister:             "/api/v2/agent/register",
			EPReport:               "/api/v2/report",
		},
	}

	common := &Common{
		appAuthorName:  "CF Automation",
		appAuthorEmail: "ms.cf.automation@aurea.com",
		licenseKey:     "unknown",
		serverURL:      "https://gravity.devfactory.com",
		logDestination: "combined",
		logLevel:       "info",
		logMaxFileSize: 10 * 1 << 20,

		// Daemon graceful timeout.
		gracefulTimeout: 15 * time.Second,
		connMapDuration: 10 * 60 * 60,
	}

	http := &HTTP{
		// This value sets actual timeout for http client.
		clientTimeout: 180 * time.Second,

		// Settings for retry in httplib.Service.
		backoffMinTimeout: 100 * time.Millisecond,
		backoffMaxTimeout: 5 * time.Minute,
		backoffFactor:     2,

		// Low level http transport settings.
		// They must not be changed at all.
		dialTimeout:           30 * time.Second,
		keepAlive:             30 * time.Second,
		tLSHandshakeTimeout:   10 * time.Second,
		idleConnTimeout:       90 * time.Second,
		expectContinueTimeout: 1 * time.Second,
	}

	settings := Settings{
		Agent:  agent,
		API:    api,
		Common: common,
		HTTP:   http,
		Operations: operationsList,
	}

	return settings
}
