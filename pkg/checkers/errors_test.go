package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemporaryError(t *testing.T) {
	errorMessage := "error"
	err := temporaryError(errorMessage)
	assert.Equal(t, errorMessage, err.Error())
	assert.Equal(t, true, err.Temporary())
}
