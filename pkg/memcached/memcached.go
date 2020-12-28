package memcached

import (
	"context"
	"errors"
	"strconv"

	"wait4it/pkg/model"

	"github.com/bradfitz/gomemcache/memcache"
)

type checker struct {
	host string
	port int
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	checker := &checker{}
	checker.buildContext(*c)
	if err := checker.validate(); err != nil {
		return nil, err
	}

	return checker, nil
}

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

func (c *checker) buildConnectionString() string {
	return c.host + ":" + strconv.Itoa(c.port)
}
