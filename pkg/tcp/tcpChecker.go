package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"wait4it/pkg/model"
)

func (tcp *Tcp) BuildContext(cx model.CheckContext) {
	tcp.Addr = cx.Host
	tcp.Port = cx.Port
}

func (tcp *Tcp) Validate() error {
	if !tcp.isPortInValidRange() {
		return errors.New("invalid port range")
	}
	return nil
}

func (tcp *Tcp) Check(ctx context.Context) (bool, bool, error) {
	var d net.Dialer
	c, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", tcp.Addr, tcp.Port))
	if err != nil {
		return false, false, err
	}
	_ = c.Close()

	return true, true, nil
}
