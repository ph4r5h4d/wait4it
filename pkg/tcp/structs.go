package tcp

import (
	"wait4it/pkg/model"
)

const (
	minPort = 1
	maxPort = 65535
)

type checker struct {
	addr string
	port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	checker := &checker{
		addr: c.Host,
		port: c.Port,
	}
	if err := checker.validate(); err != nil {
		return nil, err
	}

	return checker, nil
}
