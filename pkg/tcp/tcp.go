package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func (tcp *check) validate() error {
	if !tcp.isPortInValidRange() {
		return errors.New("invalid port range")
	}
	return nil
}

func (tcp *check) Check(ctx context.Context) (bool, bool, error) {
	var d net.Dialer
	c, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", tcp.addr, tcp.port))
	if err != nil {
		return false, false, err
	}
	_ = c.Close()

	return true, false, nil
}
