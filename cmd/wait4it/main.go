package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wait4it"
	"wait4it/pkg/checkers"
	"wait4it/pkg/cmd"
	"wait4it/pkg/environment"
)

// set by ldflags when you build.
var (
	version   = "<not set>"
	buildDate = "<not set>"
)

func main() {
	w4it, err := wait4it.NewWait4it(
		wait4it.WithOutputStream(os.Stdout),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGKILL,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	defer cancel()

	env := environment.NewEnvironment(os.Environ())

	root := cmd.NewRootCommand(w4it, os.Stdout, os.Stdin, env)
	root.AddCommand(cmd.NewVersionCommand(version, buildDate))
	root.AddCommand(cmd.NewHTTPCheckerCommand(w4it, env, checkers.NewHTTPChecker))

	root.SetArgs(os.Args[1:])
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
