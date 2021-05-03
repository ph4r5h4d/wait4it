package cmd

import (
	"github.com/spf13/cobra"
)

// NewVersionCommand creates and returns the `version` command.
func NewVersionCommand(version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Wait4it version %s (%s)\n", version, buildDate)
		},
	}

	return cmd
}
