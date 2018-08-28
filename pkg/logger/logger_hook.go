package logger

import (
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openhealthalgorithms/service/pkg/httplib"
)

var (
	// ErrAuthInvalid returned when something went wrong with auth from a backend side.
	// Either we got an expireation for token, an error has occured during the auth.
	ErrAuthInvalid = errors.New("authentication is invalid")
)

// SimpleHook sends log message via client to a specified url.
type SimpleHook struct {
	authHeader string
	url        string
	mu         sync.Mutex // protects enabled
	enabled    bool
	client     *httplib.Client
	errChan    chan error
}

// NewSimpleHook returns a new SimpleHook.
func NewSimpleHook(eCh chan error, authHeader, url string) *SimpleHook {
	return &SimpleHook{
		authHeader: authHeader,
		url:        url,
		enabled:    true,
		client:     httplib.NewClient(),
		errChan:    eCh,
	}
}

// SimpleMessage represents a message to be sent to the remote host.
type SimpleMessage struct {
	Level   string                 `json:"level"`
	Content map[string]interface{} `json:"content"`
}

// NewSimpleMessage returns a new message created from logrus.Entry.
func NewSimpleMessage(entry *logrus.Entry) *SimpleMessage {
	level := entry.Level.String()

	m := &SimpleMessage{
		Level:   level,
		Content: make(map[string]interface{}),
	}

	for k, v := range entry.Data {
		m.Content[k] = v
	}

	m.Content["msg"] = entry.Message
	m.Content["level"] = level
	m.Content["time"] = entry.Time

	return m
}

// SetAuthHeader sets authHeader for a SimpleHook.
func (h *SimpleHook) SetAuthHeader(ah string) {
	h.authHeader = ah
}

// Enable enables a SimpleHook.
func (h *SimpleHook) Enable() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enabled = true
}

// Disable disables a SimpleHook.
func (h *SimpleHook) Disable() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.enabled = false
}

// IsEnabled returns a value of hook.enabled.
func (h *SimpleHook) IsEnabled() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.enabled
}

// Fire creates a SimpleMessage from logrus.Entry and sends it using client.
func (h *SimpleHook) Fire(entry *logrus.Entry) error {
	if !h.IsEnabled() {
		return nil
	}

	var err error

	m := NewSimpleMessage(entry)
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}

	request := httplib.NewRequest("POST", h.url, j)
	request.SetContentTypeJSON()
	request.SetAuthorization(h.authHeader)

	response, err := h.client.Send(request)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
	case 401:
		err = ErrAuthInvalid
	default:
		err = errors.Errorf("status: %d, body %s", response.StatusCode, response.Body)
	}

	if err != nil {
		h.errChan <- err
		return err
	}

	return nil
}

// Levels returns a slice of all levels supported by this hook.
func (h *SimpleHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
	}
}
