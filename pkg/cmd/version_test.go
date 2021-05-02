package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionCommand(t *testing.T) {
	b := new(bytes.Buffer)
	cmd := NewVersionCommand("8.0.0", "Sun-May--2-21:52:01-UTC-2021")
	cmd.SetOut(b)
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, "Wait4it version 8.0.0 (Sun-May--2-21:52:01-UTC-2021)\n", b.String())
}
