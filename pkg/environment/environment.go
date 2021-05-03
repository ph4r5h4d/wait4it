package environment

import (
	"strings"
)

// Environment is an interface for getting os environment variables.
type Environment interface {
	// GetEnv fetches the environment variable associated with the key.
	//
	// returns an empty string when noting is found.
	GetEnv(key string) string
}

type environment struct {
	vars map[string]string
}

// NewEnvironment receives a slice of strings representing the environment
// variables in the form "key=value" and tries to parse and separate them through
// the "=" char. Then it returns an implementation of the Environment interface
// along with all of the parsed environment variables that fetchable from it.
func NewEnvironment(vars []string) Environment {
	e := &environment{
		vars: make(map[string]string, len(vars)),
	}

	for _, v := range vars {
		parts := strings.SplitN(v, "=", 2)

		if len(parts) != 2 {
			continue
		}

		e.vars[parts[0]] = parts[1]
	}

	return e
}

func (e *environment) GetEnv(key string) string {
	return e.vars[key]
}
