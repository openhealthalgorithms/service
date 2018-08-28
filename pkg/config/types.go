package config

import (
	"time"
)

// Settings represents settings for the application.
type Settings struct {
	Agent  *Agent
	API    *API
	Common *Common
	HTTP   *HTTP
	Operations []Operation
}

// Agent holds settings for agent.
type Agent struct {
	localServerURL          string
	maxTaskExecutors        int
	loopInterval            time.Duration
	startDelayDuration      time.Duration
	cmdexecTimeout          time.Duration
}

// MaxTaskExecutors returns the corresponding private value.
func (c *Agent) MaxTaskExecutors() int {
	return c.maxTaskExecutors
}

// LocalServerURL returns the corresponding private value.
func (c *Agent) LocalServerURL() string {
	return c.localServerURL
}

// LoopInterval returns the corresponding private value.
func (c *Agent) LoopInterval() time.Duration {
	return c.loopInterval
}

// StartDelayDuration returns the corresponding private value.
func (c *Agent) StartDelayDuration() time.Duration {
	return c.startDelayDuration
}

// CmdexecTimeout returns the corresponding private value.
func (c *Agent) CmdexecTimeout() time.Duration {
	return c.cmdexecTimeout
}

// API holds settings related to Gravity Backend.
type API struct {
	Methods map[string]string
}

// Common holds settings for various things.
type Common struct {
	appAuthorName   string
	appAuthorEmail  string
	serverURL       string
	licenseKey      string
	logLevel        string
	logDestination  string
	gracefulTimeout time.Duration
	logMaxFileSize  int64
	connMapDuration int
}

// AppAuthorName returns the corresponding private value.
func (c *Common) AppAuthorName() string {
	return c.appAuthorName
}

// AppAuthorEmail returns the corresponding private value.
func (c *Common) AppAuthorEmail() string {
	return c.appAuthorEmail
}

// ServerURL returns the corresponding private value.
func (c *Common) ServerURL() string {
	return c.serverURL
}

// LicenseKey returns the corresponding private value.
func (c *Common) LicenseKey() string {
	return c.licenseKey
}

// LogLevel returns the corresponding private value.
func (c *Common) LogLevel() string {
	return c.logLevel
}

// LogDestination returns the corresponding private value.
func (c *Common) LogDestination() string {
	return c.logDestination
}

// GracefulTimeout returns the corresponding private value.
func (c *Common) GracefulTimeout() time.Duration {
	return c.gracefulTimeout
}

// LogMaxFileSize returns the corresponding private value.
func (c *Common) LogMaxFileSize() int64 {
	return c.logMaxFileSize
}

// ConnMapDuration returns the corresponding private value.
func (c *Common) ConnMapDuration() int {
	return c.connMapDuration
}

// HTTP holds settings for http layer.
type HTTP struct {
	clientTimeout         time.Duration
	dialTimeout           time.Duration
	keepAlive             time.Duration
	tLSHandshakeTimeout   time.Duration
	idleConnTimeout       time.Duration
	expectContinueTimeout time.Duration
	backoffMinTimeout     time.Duration
	backoffMaxTimeout     time.Duration
	backoffFactor         int
}

// ClientTimeout returns the corresponding private value.
func (c *HTTP) ClientTimeout() time.Duration {
	return c.clientTimeout
}

// DialTimeout returns the corresponding private value.
func (c *HTTP) DialTimeout() time.Duration {
	return c.dialTimeout
}

// KeepAlive returns the corresponding private value.
func (c *HTTP) KeepAlive() time.Duration {
	return c.keepAlive
}

// TLSHandshakeTimeout returns the corresponding private value.
func (c *HTTP) TLSHandshakeTimeout() time.Duration {
	return c.tLSHandshakeTimeout
}

// IdleConnTimeout returns the corresponding private value.
func (c *HTTP) IdleConnTimeout() time.Duration {
	return c.idleConnTimeout
}

// ExpectContinueTimeout returns the corresponding private value.
func (c *HTTP) ExpectContinueTimeout() time.Duration {
	return c.expectContinueTimeout
}

// BackoffMinTimeout returns the corresponding private value.
func (c *HTTP) BackoffMinTimeout() time.Duration {
	return c.backoffMinTimeout
}

// BackoffMaxTimeout returns the corresponding private value.
func (c *HTTP) BackoffMaxTimeout() time.Duration {
	return c.backoffMaxTimeout
}

// BackoffFactor returns the corresponding private value.
func (c *HTTP) BackoffFactor() int {
	return c.backoffFactor
}

// OpOutType represents type of output of Operation.
type OpOutType int

const (
	OpJSON OpOutType = iota
	OpRawBytes
	OpFileList
)

// OpParamType represents type of OpParams of Operation.
type OpParamsType int

const (
	OpParamsEmpty OpParamsType = iota
	OpParamsConfig
	OpParamsAPI
)

// Operation represents a single operation.
//
// An Operation represents a high level structure.
// Operation -> Task -> Task Executor -> Plugin -> Output.
type Operation struct {
	Name string
	Endpoint string
	SendMethod string
	Params OpParams
	OutType OpOutType
	Timeout time.Duration
	SaveToFile bool
	TaskMap map[string]string
}

// OpParams represents params for an Operation.
type OpParams struct {
	Key string
	Template string
	Type OpParamsType
	Required bool
}
