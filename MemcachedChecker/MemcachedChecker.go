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

func (m *MemcachedConnection) Validate() (bool, error) {
	if len(m.Host) == 0 {
		return false, errors.New("Host can't be empty")
	}

	if m.Port < 0 || m.Port > 65535 {
		return false, errors.New("Invalid port range for Memcached")
	}

	return true, nil
}

func (m *MemcachedConnection) Check() (bool, bool, error) {
	mc := memcache.New(m.BuildConnectionString())

	// Since the current version of the Memcached Golang client doesn't create
	// any connection to the Memcached server (or at least the client that I'm using)
	// at the New method and at here we don't have any way for checking the correctness
	// of the given server address, we must use another way for checking connection. As I
	// understand, it seems the connection will establish at setting and getting items time,
	// and I think this is a good point for checking connection.
	err := mc.Set(&memcache.Item{Key: "foo", Value: []byte("bar")})
	if err != nil {
		return false, true, err
	}
	mc.Delete("foo")

	return true, true, nil
}
