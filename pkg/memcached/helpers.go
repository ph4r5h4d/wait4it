package memcached

import (
	"strconv"
)

func (m *MemcachedConnection) BuildConnectionString() string {
	return m.Host + ":" + strconv.Itoa(m.Port)
}
