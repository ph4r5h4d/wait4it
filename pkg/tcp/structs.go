package tcp

import (
	"context"
	"net"
	"wait4it/pkg/model"
)

const (
	minPort = 1
	maxPort = 65535
)

// Dialer is an interface for dialing TCP connections.
// This allows injection of custom dialers in tests.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type Check struct {
	Addr  string
	Port  int
	dialer Dialer
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &Check{
		Addr: c.Host,
		Port: c.Port,
	}
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}

func (c *Check) BuildContext(cx model.CheckContext) {
	c.Addr = cx.Host
	c.Port = cx.Port
}
