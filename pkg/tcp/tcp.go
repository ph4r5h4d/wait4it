package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
)

func (c *checker) validate() error {
	if !c.isPortInValidRange() {
		return errors.New("invalid port range")
	}
	return nil
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
