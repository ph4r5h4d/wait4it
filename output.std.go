package main

import (
	"fmt"
	"os"
)

func wStdOut(r bool) {
	if r {
		_, _ = fmt.Fprint(os.Stdout, "succeed")
		os.Exit(0)
	} else {
		fmt.Fprint(os.Stdout, ".")
	}
}

func wStdErr(a ...interface{}){
	fmt.Fprint(os.Stderr,a)
}