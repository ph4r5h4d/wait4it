package cmd

import (
	"fmt"
	"os"
)

func wStdOut(r bool) {
	if r {
		_, _ = fmt.Fprint(os.Stdout, "succeed")
		os.Exit(0)
	} else {
		_, _ = fmt.Fprint(os.Stdout, ".")
	}
}

func wStdErr(a ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, a...)
}
