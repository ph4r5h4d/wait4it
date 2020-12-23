// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build build the wait4it binary. This build is static and is not using CGO
func Build() error {
	gocmd := mg.GoCmd()
	os.Setenv("CGO_ENABLED", "0")
	return sh.RunV(gocmd, "build", "-a", "-installsuffix", "cgo", "-o", "w4it", "./cmd/w4it")
}

// Docker build the docker file.
func Docker() error {
	return sh.RunV("docker", "build", ".", "--file", "Dockerfile", "--tag", "wait4it")
}
