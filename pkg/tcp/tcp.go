package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func (tcp *Check) Validate() error {
	if !tcp.IsPortInValidRange() {
		return errors.New("invalid port range")
	}
	return nil
}

func (tcp *Check) Check(ctx context.Context) (bool, bool, error) {
	dialer := tcp.dialer
	if dialer == nil {
		dialer = &net.Dialer{}
	}
	c, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", tcp.Addr, tcp.Port))
	if err != nil {
		return false, false, err
	}
	_ = c.Close()

	return true, false, nil
}
