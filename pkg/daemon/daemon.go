package daemon

import (
	"github.com/pkg/errors"
)

var (
	// ErrProcessNotFound returned when process for Daemon is nil.
	ErrProcessNotFound = errors.New("process must be set")
)

// Runner represents a service with Run.
type Runner interface {
	Run() error
}

// Service is the easiest implementation of Runner.
// It can be used to wrap something not satisfying the requirements of Runner.
type Service struct {
	Cmd func() error
}

// Run wraps Cmd.
func (s Service) Run() error {
	return s.Cmd()
}

// runProcess just calls Run of the given Runner.
func runProcess(errChan chan<- error, srv Runner) {
	errChan <- srv.Run()
}
