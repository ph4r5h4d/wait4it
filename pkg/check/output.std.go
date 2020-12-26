package check

import (
	"fmt"
	"os"
)

func wStdOut(r bool) {
	if r {
		_, _ = fmt.Fprintln(os.Stdout, "succeed")
		os.Exit(0)
	} else {
		_, _ = fmt.Fprint(os.Stdout, ".")
	}
}
