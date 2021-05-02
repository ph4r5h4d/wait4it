package cmd

import (
	"wait4it"
	"wait4it/pkg/environment"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

type httpCheckerOptions struct {
	url            string
	statusCode     int
	body           string
	followRedirect bool
}

// NewHTTPCheckerFunc represents the function that returns the HTTP checker.
type NewHTTPCheckerFunc func(url string, status int, body string, followRedirect bool) (wait4it.Checkable, error)

// NewHTTPCheckerCommand creates and returns the `http` command.
func NewHTTPCheckerCommand(w4it *wait4it.Wait4it, env environment.Environment, fn NewHTTPCheckerFunc) *cobra.Command {
	var options httpCheckerOptions

	cmd := &cobra.Command{
		Use:   "http <URL>",
		Short: "Check http service whther it's ready",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			$ wait4it http https://farshad.nematdoust.com --body "Software Engineer" --code 200 --follow-redirect
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			options.url = args[0]

			checker, err := fn(options.url, options.statusCode, options.body, options.followRedirect)
			if err != nil {
				return err
			}

			return w4it.Run(cmd.Context(), checker)
		},
	}

	cmd.Flags().StringVarP(&options.body, "body", "b", defaultEnvString(env, "W4IT_HTTP_BODY", ""), "body to check inside http response")
	cmd.Flags().IntVarP(&options.statusCode, "code", "c", defaultEnvInt(env, "W4IT_HTTP_STATUS_CODE", 200), "status code to be expected from http call")
	cmd.Flags().BoolVar(&options.followRedirect, "follow-redirect", defaultEnvBool(env, "W4IT_HTTP_FOLLOW_REDIRECT", false), "whether to follow the redirect while doing the HTTP check")

	return cmd
}
