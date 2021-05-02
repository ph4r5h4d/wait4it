package cmd

import (
	"testing"
	"time"

	"wait4it/pkg/environment"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEnvString(t *testing.T) {
	e := environment.NewEnvironment([]string{"key=value"})

	assert.Equal(t, "value", defaultEnvString(e, "key", ""))
	assert.Equal(t, "default", defaultEnvString(e, "k", "default"))
}

func TestDefaultEnvInt(t *testing.T) {
	e := environment.NewEnvironment([]string{"key=123"})

	assert.Equal(t, 123, defaultEnvInt(e, "key", 0))
	assert.Equal(t, 999, defaultEnvInt(e, "k", 999))
}

func TestDefaultEnvUint(t *testing.T) {
	e := environment.NewEnvironment([]string{"key=456"})

	assert.Equal(t, uint(456), defaultEnvUint(e, "key", 0))
	assert.Equal(t, uint(999), defaultEnvUint(e, "k", 999))
}

func TestDefaultEnvDuration(t *testing.T) {
	e := environment.NewEnvironment([]string{"key=1s"})

	assert.Equal(t, 1*time.Second, defaultEnvDuration(e, "key", 0))
	assert.Equal(t, 2*time.Second, defaultEnvDuration(e, "k", 2*time.Second))
}

func TestDefaultEnvBool(t *testing.T) {
	e := environment.NewEnvironment([]string{"key=true", "invalid=fals"})

	assert.Equal(t, true, defaultEnvBool(e, "key", false))
	assert.Equal(t, false, defaultEnvBool(e, "k", false))
}
