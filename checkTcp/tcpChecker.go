package checkTcp

import (
	"errors"
	"fmt"
	"net"
)

func (i IP) DoesPortAcceptConnection() (bool, error) {
	if !i.isPortInValidRange() {
		return false, errors.New("invalid port range")
	}

	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", i.Addr, i.Port))
	if err != nil {
		return false, nil
	}
	_ = c.Close()

	return true, nil
}
