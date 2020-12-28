package memcached

import (
	"context"
	"errors"

	"wait4it/pkg/model"

	"github.com/bradfitz/gomemcache/memcache"
)

func (c *checker) buildContext(cx model.CheckContext) {
	c.host = cx.Host
	c.port = cx.Port
}

func (c *checker) validate() error {
	if len(c.host) == 0 {
		return errors.New("Host can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("Invalid port range for Memcached")
	}

	return nil
}

func (c *checker) Check(context.Context) (bool, bool, error) {
	// TODO: is it possible to handle ping using context?
	mc := memcache.New(c.buildConnectionString())

	if err := mc.Ping(); err != nil {
		return false, false, err
	}

	return true, true, nil
}
