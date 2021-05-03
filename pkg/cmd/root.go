package cmd

import (
	"io"
	"time"

	"wait4it"
	"wait4it/pkg/environment"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type rootOptions struct {
	checkingInterval time.Duration
	timeout          time.Duration
	retries          uint
}

// NewRootCommand creates and returns the `wait4it` command.
func NewRootCommand(w4it *wait4it.Wait4it, output, err io.Writer, env environment.Environment) *cobra.Command {
	var options rootOptions

	cmd := &cobra.Command{
		Use:   "wait4it",
		Short: "Wait4it CLI",
		Long: heredoc.Doc(`
			A simple Go application to test whether a port is ready to accept a connection or check MySQL, PostgreSQL, MongoDB,
			or Redis server is ready or not, Also you can do an HTTP call and check the response code and text in the response body.
		`),
		Example: heredoc.Doc(`
			$ wait4it http https://farshad.nematdoust.com --body "Software Engineer" --code 200 --follow-redirect
		`),
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// https://github.com/spf13/cobra/issues/340
			cmd.SilenceUsage = true

			return w4it.Apply(
				wait4it.WithCheckingInterval(options.checkingInterval),
				wait4it.WithTimeout(options.timeout),
				wait4it.WithMaxRetries(options.retries),
			)
		},
	}

	cmd.SetOut(output)
	cmd.SetErr(err)

	cmd.PersistentFlags().DurationVarP(&options.checkingInterval, "interval", "i", defaultEnvDuration(env, "W4IT_CHECKING_INTERVAL", time.Second), "the checking interval period in between each of checking.")
	cmd.PersistentFlags().DurationVarP(&options.timeout, "timeout", "t", defaultEnvDuration(env, "W4IT_TIMEOUT_DURATION", 30*time.Second), "amount of time the Wait4it waits and then fails")
	cmd.PersistentFlags().UintVarP(&options.retries, "retries", "r", defaultEnvUint(env, "W4IT_RETRIES", 10), "the maximum number of retrying to check after a failed check.")

	return cmd
}
