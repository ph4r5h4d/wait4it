// +build mage

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build build the wait4it binary. This build is static and is not using CGO
func Build() error {
	gocmd := mg.GoCmd()
	os.Setenv("CGO_ENABLED", "0")
	return sh.RunV(gocmd, "build", "-a", "-installsuffix", "cgo", "-ldflags", ldflags(), "-o", "wait4it", "./cmd/wait4it")
}

// Docker build the docker file.
func Docker() error {
	return sh.RunV("docker", "build", ".", "--file", "Dockerfile", "--tag", "wait4it")
}

func ldflags() string {
	return fmt.Sprintf(`-X "main.version=%s" -X "main.buildDate=%s"`, version(), buildDate())
}

func version() string {
	s, _ := sh.Output("git", "describe", "--tags")
	parts := strings.SplitN(s, "-", 2)
	return parts[0]
}

func buildDate() string {
	return time.Now().Format(time.RFC822)
}
