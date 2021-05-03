package wait4it

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	// ErrTimeout occurs when the maximum timeout is reached.
	ErrTimeout = errors.New("wait4it: timeout is reached")

	// ErrMaxRetry occurs when the number of retries on the checker reaches its maximum.
	ErrMaxRetry = errors.New("wait4it: maximum retry")

	// ErrCheckerFailed occurs when a checker fails during checking.
	ErrCheckerFailed = errors.New("wait4it: checker failed")
)

// Checkable is an interface that checks the health of services. Each service
// that tends to perform health checking should implement the Checkable interface.
type Checkable interface {
	// Check checks the health of the service and returns an error.
	//
	// There is a possibility that we may encounter a temporary failure while checking.
	// In such circumstances, as contract, all checkers should implement the temporary
	// interface as a sign that we can retry.
	Check(ctx context.Context) error
}

type temporary interface {
	Temporary() bool
}

// CheckFunc is an adapter type to allow the use of ordinary functions as a checker.
type CheckFunc func(context.Context) error

// Check calls the adapted ordinary function.
func (f CheckFunc) Check(ctx context.Context) error {
	return f(ctx)
}

// OptionFunc represents the function that receives the Wait4it and applies a configuration
// to that. OptionFunc is provided for use in the Wait4it construction and Apply method.
type OptionFunc func(w *Wait4it) error

// Wait4it provides functionality to health checking services with the support of retrying,
// log progressing, retry limitation, timeout limitation, and so on.
type Wait4it struct {
	checkingInterval time.Duration
	output           io.Writer
	timeout          time.Duration // zero value means there's no timing limitation.
	maxRetries       uint          // zero value means there's no retrying limitation.
}

// NewWait4it returns a new Wait4it and an error. Variadic options of OptionFunc are receivable
// during construction. NewWait4it tries to apply these options to the Wait4it and if anything
// goes wrong, an error will be returned.
func NewWait4it(opts ...OptionFunc) (*Wait4it, error) {
	w := &Wait4it{
		checkingInterval: time.Second,
		output:           io.Discard,
	}

	if err := w.Apply(opts...); err != nil {
		return nil, err
	}

	return w, nil
}

// Apply applies all the given options to the Wait4it.
//
// Even though we have the same functionality on the object constructor, in certain situations
// the object configuration may need to be changed.
func (w *Wait4it) Apply(opts ...OptionFunc) error {
	for _, opt := range opts {
		if err := opt(w); err != nil {
			return err
		}
	}

	return nil
}

// Run runs the checker and returns a nil error when everything went well.
//
// Given that every checker returns an error that's maybe an implementation of the temporary
// interface if it is temporary, Run takes advantage of this and supports the retry mechanism.
// After each failed check, if it is retriable, Run performs another check. This cycle continues
// until reaches the limitation point. Between each of these retry checks, there is also an interval
// for waiting and then perform another check. By default, this interval duration is considered one
// second.
//
// Run has a timeout and maximum retry strategies for cancellation and limitation. After each failed
// check, if each of these limitations reaches a point where they are no longer able to let continue
// retrying, a suitable error message will be returned. By default, no limitation is considered.
func (w *Wait4it) Run(ctx context.Context, checker Checkable) error {
	fmt.Fprint(w.output, "Wait4it...")

	if w.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, w.timeout)
		defer cancel()
	}

	if err := w.ticker(ctx, checker); err != nil {
		fmt.Fprintln(w.output, "failed")

		return err
	}

	fmt.Fprintln(w.output, "succeed")

	return nil
}

func (w *Wait4it) ticker(ctx context.Context, checker Checkable) error {
	var tries uint

	t := time.NewTicker(w.checkingInterval)
	defer t.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return ErrTimeout
			}

			return ctx.Err()
		case <-t.C:
			fmt.Fprint(w.output, ".")

			err := checker.Check(ctx)
			if err != nil && !isTemporary(err) {
				return err
			}
			tries++

			if err == nil {
				break loop
			}

			if w.hasRetryLimitation() && w.hasReachedItsMaximumRetry(tries) {
				return ErrMaxRetry
			}
		}
	}

	return nil
}

func isTemporary(err error) bool {
	temp, ok := err.(temporary)
	return ok && temp.Temporary()
}

func (w *Wait4it) hasRetryLimitation() bool {
	return w.maxRetries != 0
}

func (w *Wait4it) hasReachedItsMaximumRetry(tries uint) bool {
	return tries == w.maxRetries
}

// WithCheckingInterval applies the checking interval duration to the Wait4it.
func WithCheckingInterval(d time.Duration) OptionFunc {
	return func(w *Wait4it) error {
		if d <= 0 {
			return errors.New("wait4it: non-positive interval")
		}

		w.checkingInterval = d
		return nil
	}
}

// WithOutputStream applies the output stream for progressing logs to the Wait4it.
func WithOutputStream(output io.Writer) OptionFunc {
	return func(w *Wait4it) error {
		if output == nil {
			return errors.New("wait4it: nil output stream")
		}

		w.output = output
		return nil
	}
}

// WithMaxRetries applies the maximum retry limitation to the Wait4it.
func WithMaxRetries(retries uint) OptionFunc {
	return func(w *Wait4it) error {
		w.maxRetries = retries
		return nil
	}
}

// WithTimeout applies the timeout duration limitation to the Wait4it
func WithTimeout(timeout time.Duration) OptionFunc {
	return func(w *Wait4it) error {
		w.timeout = timeout
		return nil
	}
}
