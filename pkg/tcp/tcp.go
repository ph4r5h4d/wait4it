package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"

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

func (c *checker) validate() error {
	if !c.isPortInValidRange() {
		return errors.New("invalid port range")
	}
	return nil
}

func (c *checker) isPortInValidRange() bool {
	if c.port < minPort || c.port > maxPort {
		return false
	}
	return true
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	var d net.Dialer
	client, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", c.addr, c.port))
	if err != nil {
		return false, false, err
	}
	_ = client.Close()

	return true, false, nil
}
