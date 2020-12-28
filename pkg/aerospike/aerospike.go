package aerospike

import (
	"context"
	"errors"

	"wait4it/pkg/model"

	"github.com/aerospike/aerospike-client-go"
)

func (c *checker) buildContext(cx model.CheckContext) {
	c.host = cx.Host
	c.port = cx.Port
}

func (c *checker) validate() error {
	if len(c.host) == 0 {
		return errors.New("host can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("invalid port range for Memcached")
	}

	return nil
}

func (c *checker) Check(context.Context) (bool, bool, error) {
	// TODO: is it possible to handle ping using context?
	client, err := aerospike.NewClient(c.host, c.port)

	if err != nil {
		return false, false, err
	}

	if !client.IsConnected() {
		return false, false, errors.New("client is not connected")
	}

	return true, true, nil
}
