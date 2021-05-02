package wait4it

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockChecker struct {
	checkFunc func(ctx context.Context) error
}

func NewMockChecker() *MockChecker {
	return &MockChecker{
		checkFunc: func(ctx context.Context) error {
			return nil
		},
	}
}

func (m *MockChecker) Check(ctx context.Context) error {
	return m.checkFunc(ctx)
}

type temporaryError string

func (e temporaryError) Error() string   { return string(e) }
func (e temporaryError) Temporary() bool { return true }

func TestCheckFunc(t *testing.T) {
	checker := CheckFunc(func(c context.Context) error { return nil })
	assert.NoError(t, checker.Check(context.Background()))
}

func TestNewWait4it(t *testing.T) {
	tt := []struct {
		name            string
		in              []OptionFunc
		out             *Wait4it
		isErrorExpected bool
	}{
		{
			name: "wait4it constructor without options",
			in:   []OptionFunc{},
			out: &Wait4it{
				checkingInterval: time.Second,
				output:           io.Discard,
				timeout:          0,
				maxRetries:       0,
			},
			isErrorExpected: false,
		},
		{
			name: "wait4it constructor with all options",
			in: []OptionFunc{
				WithCheckingInterval(2 * time.Second),
				WithOutputStream(&bytes.Buffer{}),
				WithTimeout(30 * time.Second),
				WithMaxRetries(10),
			},
			out: &Wait4it{
				checkingInterval: 2 * time.Second,
				output:           &bytes.Buffer{},
				timeout:          30 * time.Second,
				maxRetries:       10,
			},
			isErrorExpected: false,
		},
		{
			name: "wait4it constructor with nil output stream",
			in: []OptionFunc{
				WithOutputStream(nil),
			},
			out:             nil,
			isErrorExpected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w4it, err := NewWait4it(tc.in...)
			if tc.isErrorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, w4it, tc.out)
		})
	}
}

func TestWait4itRun(t *testing.T) {
	w4it, err := NewWait4it(WithCheckingInterval(time.Nanosecond))
	require.NoError(t, err)

	err = w4it.Run(context.Background(), NewMockChecker())
	require.NoError(t, err)
}

func TestWait4itRunCancellationThroughContext(t *testing.T) {
	w4it, err := NewWait4it(WithCheckingInterval(time.Nanosecond))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = w4it.Run(ctx, NewMockChecker())
	require.Error(t, err)
	require.IsType(t, err, context.Canceled)
}

func TestWait4itRunTimeoutError(t *testing.T) {
	w4it, err := NewWait4it(WithTimeout(-1))
	require.NoError(t, err)

	err = w4it.Run(context.Background(), NewMockChecker())
	require.Error(t, err)
	require.IsType(t, err, ErrTimeout)
}

func TestWait4itRunMaxRetry(t *testing.T) {
	maxRetries := 2
	w4it, err := NewWait4it(
		WithCheckingInterval(time.Nanosecond),
		WithMaxRetries(uint(maxRetries)),
	)
	require.NoError(t, err)

	var tries int
	mock := NewMockChecker()
	mock.checkFunc = func(ctx context.Context) error {
		tries++
		return temporaryError("error")
	}

	err = w4it.Run(context.Background(), mock)
	require.Error(t, err)
	require.IsType(t, err, ErrMaxRetry)

	assert.Equal(t, maxRetries, tries)
}

func TestWait4itRunUnretrievableCheckFailer(t *testing.T) {
	w4it, err := NewWait4it(
		WithCheckingInterval(time.Nanosecond),
	)
	require.NoError(t, err)

	mock := NewMockChecker()
	mock.checkFunc = func(ctx context.Context) error {
		return errors.New("error")
	}

	err = w4it.Run(context.Background(), mock)
	require.Error(t, err)
}

func TestWait4itRunOutputStreamOnSuccessAfterFewRetries(t *testing.T) {
	b := new(bytes.Buffer)

	maxRetries := 5
	w4it, err := NewWait4it(
		WithCheckingInterval(time.Nanosecond),
		WithOutputStream(b),
		WithMaxRetries(uint(maxRetries)+1),
	)
	require.NoError(t, err)

	var tries int
	mock := NewMockChecker()
	mock.checkFunc = func(ctx context.Context) error {
		tries++
		if tries == maxRetries {
			return nil
		}

		return temporaryError("error")
	}

	err = w4it.Run(context.Background(), mock)
	require.NoError(t, err)

	expected := "Wait4it..." + strings.Repeat(".", maxRetries) + "succeed\n"
	assert.Equal(t, expected, b.String())
}

func TestWait4itRunOutputStreamOnFailer(t *testing.T) {
	b := new(bytes.Buffer)

	maxRetries := 5
	w4it, err := NewWait4it(
		WithCheckingInterval(time.Nanosecond),
		WithOutputStream(b),
		WithMaxRetries(uint(maxRetries)),
	)
	require.NoError(t, err)

	mock := NewMockChecker()
	mock.checkFunc = func(ctx context.Context) error {
		return temporaryError("error")
	}

	err = w4it.Run(context.Background(), mock)
	require.EqualError(t, err, ErrMaxRetry.Error())

	expected := "Wait4it..." + strings.Repeat(".", maxRetries) + "failed\n"
	assert.Equal(t, expected, b.String())
}

func TestWithCheckingInterval(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	err = WithCheckingInterval(10 * time.Second)(w4it)
	require.NoError(t, err)

	assert.Equal(t, 10*time.Second, w4it.checkingInterval)
}

func TestWithCheckingIntervalNonePositiveDuration(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	err = WithCheckingInterval(-1)(w4it)
	require.Error(t, err)
}

func TestWithOutputStream(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	b := new(bytes.Buffer)
	err = WithOutputStream(b)(w4it)
	require.NoError(t, err)

	assert.Equal(t, b, w4it.output)
}

func TestWithOutputStreamNilStream(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	err = WithOutputStream(nil)(w4it)
	require.Error(t, err)
}

func TestWithMaxRetries(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	err = WithMaxRetries(5)(w4it)
	require.NoError(t, err)

	assert.Equal(t, uint(5), w4it.maxRetries)
}

func TestWithTimeout(t *testing.T) {
	w4it, err := NewWait4it()
	require.NoError(t, err)

	err = WithTimeout(5 * time.Second)(w4it)
	require.NoError(t, err)

	assert.Equal(t, 5*time.Second, w4it.timeout)
}
