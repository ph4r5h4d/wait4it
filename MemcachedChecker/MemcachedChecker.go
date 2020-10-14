package MemcachedChecker

import (
	"errors"
	"wait4it/model"

	"github.com/bradfitz/gomemcache/memcache"
)

func (m *MemcachedConnection) BuildContext(cx model.CheckContext) {
	m.Host = cx.Host
	m.Port = cx.Port
}

func (m *MemcachedConnection) Validate() error {
	if len(m.Host) == 0 {
		return errors.New("Host can't be empty")
	}

	if m.Port < 1 {
		return errors.New("Invalid port range for Memcached")
	}

	return nil
}

func (m *MemcachedConnection) Check() (bool, bool, error) {
	mc := memcache.New(m.BuildConnectionString())

	if err := mc.Ping(); err != nil {
		return false, false, err
	}

	return true, true, nil
}
