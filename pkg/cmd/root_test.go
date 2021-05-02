package cmd

import (
	"bytes"
	"testing"
	"time"

	"wait4it"
	"wait4it/pkg/environment"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCommand(t *testing.T) {
	tt := []struct {
		name            string
		expectedOptions rootOptions
		rawEnvVars      []string
		flags           map[string]string
	}{
		{
			name: "default flag values",
			expectedOptions: rootOptions{
				checkingInterval: time.Second,
				timeout:          30 * time.Second,
				retries:          10,
			},
		},
		{
			name: "with flags",
			expectedOptions: rootOptions{
				checkingInterval: 2 * time.Second,
				timeout:          15 * time.Second,
				retries:          5,
			},
			flags: map[string]string{
				"interval": "2s",
				"timeout":  "15s",
				"retries":  "5",
			},
		},
		{
			name: "with env as flags fallback",
			expectedOptions: rootOptions{
				checkingInterval: 3 * time.Second,
				timeout:          10 * time.Second,
				retries:          5,
			},
			rawEnvVars: []string{
				"W4IT_CHECKING_INTERVAL=3s",
				"W4IT_TIMEOUT_DURATION=10s",
				"W4IT_RETRIES=5",
			},
		},
	}

	w4it, err := wait4it.NewWait4it()
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			env := environment.NewEnvironment(tc.rawEnvVars)

			outputStream := new(bytes.Buffer)
			errStream := new(bytes.Buffer)

			root := NewRootCommand(w4it, outputStream, errStream, env)
			assert.Equal(t, outputStream, root.OutOrStdout())
			assert.Equal(t, errStream, root.OutOrStderr())

			for k, v := range tc.flags {
				err := root.PersistentFlags().Set(k, v)
				require.NoError(t, err)
			}

			require.NoError(t, root.Execute())

			checkingInterval, err := root.PersistentFlags().GetDuration("interval")
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOptions.checkingInterval, checkingInterval)

			timeout, err := root.PersistentFlags().GetDuration("timeout")
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOptions.timeout, timeout)

			retries, err := root.PersistentFlags().GetUint("retries")
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOptions.retries, retries)
		})
	}
}
