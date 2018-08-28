// +build linux darwin

package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// Daemon represents a service to run as a daemon.
//
// It must be initialized with NewDaemon before use.
type Daemon struct {
	ctx           context.Context
	cancel        context.CancelFunc
	graceful      bool
	process       Runner
	processErrors chan error
	signals       chan os.Signal
	timeout       time.Duration
}

// NewDaemon returns an instance of a Daemon.
func NewDaemon() *Daemon {
	ctx, cancel := context.WithCancel(context.Background())
	return &Daemon{
		ctx:           ctx,
		cancel:        cancel,
		processErrors: make(chan error, 1),
		signals:       make(chan os.Signal, 1),
	}
}

// SetProcess sets the process to be run.
func (d *Daemon) SetProcess(r Runner) {
	// There is no need to check type of r before setting it.
	// If it does not implement the Runner it would not compile.
	d.process = r
}

// GetCtx returns a new context derived from the daemon.
func (d *Daemon) GetCtx() (context.Context, func()) {
	return context.WithCancel(d.ctx)
}

// SetGraceful sets graceful shutdown mode.
func (d *Daemon) SetGraceful(mode bool) {
	d.graceful = mode
}

// SetTimeout sets timeout for graceful shutdown.
func (d *Daemon) SetTimeout(t time.Duration) {
	d.timeout = t
}

// Start starts Daemon and listens for signals to stop.
func (d *Daemon) Start() error {
	if d.process == nil {
		return errors.Wrap(ErrProcessNotFound, "daemon Start")
	}

	// Subscribe for signals from the system.
	signal.Notify(d.signals, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// Run the service.
	go runProcess(d.processErrors, d.process)

	// Listen channels for events.
	select {
	case err := <-d.processErrors:
		d.cancel()
		return err
	case <-d.signals:
		d.cancel()
		if d.graceful && d.timeout > 0 {
			gCtx, gCancel := context.WithTimeout(context.Background(), d.timeout)
			defer gCancel()
			select {
			case <-gCtx.Done():
			case err := <-d.processErrors:
				return err
			}
		}
		// switch srv := d.process.(type) {
		// case Service:
		// }
	}

	return nil
}
