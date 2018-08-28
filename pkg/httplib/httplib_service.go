package httplib

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	maxInt64 = math.MaxInt64 - 512
	// RetryInfinite specifies the infinite number of retries.
	RetryInfinite = 0
	// Retry10 specifies 10 retries.
	Retry10 = 10
)

var (
	// defaultBackoffMin represents starting interval for retry.
	defaultBackoffMin = 100 * time.Millisecond
	// defaultBackoffMax represents max interval for retry.
	defaultBackoffMax = 5 * time.Minute
	// defaultBackoffFactor represents factor for exponential growth.
	defaultBackoffFactor = 2
)

// Service represents a http service ready for incoming requests.
type Service struct {
	client    *Client
	ctx       context.Context
	ctxCancel context.CancelFunc
	log       *logrus.Entry
	request   chan *Request
	response  chan *ResponseWithErr
	bMin      time.Duration
	bMax      time.Duration
	bFactor   int
}

// ResponseWithErr holds Reponse and Err.
type ResponseWithErr struct {
	Response *Response
	Err      error
}

// NewService returns a new service.
func NewService(pCtx context.Context, client *Client, log *logrus.Entry) *Service {
	return newServiceWithSettings(pCtx, client, log, defaultBackoffMin, defaultBackoffMax, defaultBackoffFactor)
}

// NewServiceWithSettings returns a new service with low level settings for Backoff Retry mechanism.
//
// WARNING! You MUST totally understand what are you doing by passing incorrect values.
func NewServiceWithSettings(
	pCtx context.Context,
	client *Client,
	log *logrus.Entry,
	bMin, bMax time.Duration,
	bFactor int,
) *Service {
	return newServiceWithSettings(pCtx, client, log, bMin, bMax, bFactor)
}

// newServiceWithSettings is a private constructor used internally.
func newServiceWithSettings(
	pCtx context.Context,
	client *Client,
	log *logrus.Entry,
	bMin, bMax time.Duration,
	bFactor int,
) *Service {
	ctx, ctxCancel := context.WithCancel(pCtx)

	return &Service{
		client:    client,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		log:       log,
		request:   make(chan *Request),
		response:  make(chan *ResponseWithErr),
		bMin:      bMin,
		bMax:      bMax,
		bFactor:   bFactor,
	}
}

// Run starts service.
func (s *Service) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			s.log.Debugf("httpService got ctx cancel")
			return
		case req := <-s.request:
			resp, err := s.client.Send(req)
			s.response <- &ResponseWithErr{Response: resp, Err: err}
		}
	}
}

// Send makes a Request and returns ResponseWithErr.
func (s *Service) Send(req *Request) *ResponseWithErr {
	s.request <- req

	return <-s.response
}

// SendWithRetry makes a Request and returns ResponseWithErr.
func (s *Service) SendWithRetry(req *Request, attempts int) *ResponseWithErr {
	var result *ResponseWithErr
	done := false
	maxAttempts := float64(maxInt64)
	if attempts > RetryInfinite {
		maxAttempts = float64(attempts)
	}
	b := &backoff.Backoff{
		Min:    s.bMin,
		Max:    s.bMax,
		Factor: float64(s.bFactor),
	}
	for b.Attempt() < maxAttempts && !done {
		result = s.Send(req)
		err := result.Err
		if err != nil {
			s.log.Errorf("http request attempt #%.f failed. Error: %v", b.Attempt()+1, err)
			d := b.Duration()
			// Call to Duration increments Attempt.
			if b.Attempt() < maxAttempts {
				s.log.Infof("http request retry after: %.3fs", d.Seconds())
				time.Sleep(d)
			}
			continue
		}
		code := result.Response.StatusCode
		msg := result.Response.Status
		if code >= 500 {
			s.log.Errorf("http request attempt #%.f failed. Response: code %d message %s", b.Attempt()+1, code, msg)
			d := b.Duration()
			// Call to Duration increments Attempt.
			if b.Attempt() < maxAttempts {
				s.log.Infof("http request retry after: %.3fs", d.Seconds())
				time.Sleep(d)
			}
			continue
		}

		done = true
	}

	if result.Err != nil {
		if b.Attempt() >= maxAttempts {
			result.Err = errors.Wrapf(result.Err, "failed to finish request after %.f retries", maxAttempts)
		}
	}

	return result
}
