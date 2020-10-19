package TcpChecker

import (
	"errors"
	"fmt"
	"net"
	"wait4it/model"
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

func (tcp *Tcp) Check() (bool, bool, error) {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", tcp.Addr, tcp.Port))
	if err != nil {
		return false, false, err
	}
	_ = c.Close()

	return true, true, nil
}
