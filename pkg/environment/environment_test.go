package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment([]string{
		"key1=value1",
		"key2=value2=value2",
		"key3 value",
	})

	assert.Equal(t, "value1", env.GetEnv("key1"))
	assert.Equal(t, "value2=value2", env.GetEnv("key2"))
	assert.Equal(t, "", env.GetEnv("key3 value"))
	assert.Equal(t, "", env.GetEnv("some_rand_key"))
}
