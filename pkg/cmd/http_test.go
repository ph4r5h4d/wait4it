package cmd

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"wait4it"
	"wait4it/pkg/environment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPCheckerCommand(t *testing.T) {
	tt := []struct {
		name            string
		expectedOptions httpCheckerOptions
		rawEnvVars      []string
		flags           map[string]string
	}{
		{
			name: "default flag values",
			expectedOptions: httpCheckerOptions{
				url:            "https://example.com/",
				statusCode:     200,
				body:           "",
				followRedirect: false,
			},
		},
		{
			name: "with flags",
			expectedOptions: httpCheckerOptions{
				url:            "https://example.com/",
				statusCode:     201,
				body:           "http body",
				followRedirect: true,
			},
			flags: map[string]string{
				"body":            "http body",
				"code":            "201",
				"follow-redirect": "true",
			},
		},
		{
			name: "with env as flags fallback",
			expectedOptions: httpCheckerOptions{
				url:            "https://example.com/",
				statusCode:     202,
				body:           "http body env",
				followRedirect: false,
			},
			rawEnvVars: []string{
				"W4IT_HTTP_STATUS_CODE=202",
				"W4IT_HTTP_BODY=http body env",
				"W4IT_HTTP_FOLLOW_REDIRECT=false",
			},
		},
	}

	w4it, err := wait4it.NewWait4it(wait4it.WithCheckingInterval(time.Nanosecond))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			env := environment.NewEnvironment(tc.rawEnvVars)

			httpCheckerFunc := func(url string, statusCode int, body string, followRedirect bool) (wait4it.Checkable, error) {
				assert.Equal(t, tc.expectedOptions.url, url)
				assert.Equal(t, tc.expectedOptions.statusCode, statusCode)
				assert.Equal(t, tc.expectedOptions.body, body)
				assert.Equal(t, tc.expectedOptions.followRedirect, followRedirect)

				return wait4it.CheckFunc(func(context.Context) error {
					return nil
				}), nil
			}

			cmd := NewHTTPCheckerCommand(w4it, env, httpCheckerFunc)
			cmd.SetOut(io.Discard)
			cmd.SetArgs([]string{tc.expectedOptions.url})

			for k, v := range tc.flags {
				err := cmd.Flags().Set(k, v)
				require.NoError(t, err)
			}

			assert.NoError(t, cmd.Execute())
		})
	}
}

func TestNewHTTPCheckerCommandErrors(t *testing.T) {
	tt := []struct {
		name            string
		args            []string
		flags           map[string]string
		httpCheckerFunc NewHTTPCheckerFunc
		expectedErr     string
	}{
		{
			name:        "without required args",
			expectedErr: "accepts 1 arg(s), received 0",
		},
		{
			name: "validation error",
			args: []string{"https://example.com"},
			httpCheckerFunc: func(url string, status int, body string, followRedirect bool) (wait4it.Checkable, error) {
				return nil, errors.New("validation error")
			},
			expectedErr: "validation error",
		},
	}

	w4it, err := wait4it.NewWait4it(wait4it.WithCheckingInterval(time.Nanosecond))
	require.NoError(t, err)

	env := environment.NewEnvironment(nil)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewHTTPCheckerCommand(w4it, env, tc.httpCheckerFunc)
			cmd.SetOut(io.Discard)
			cmd.SetArgs(tc.args)

			for k, v := range tc.flags {
				cmd.Flags().Set(k, v)
			}

			err := cmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}
