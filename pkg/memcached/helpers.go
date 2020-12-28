package memcached

import (
	"strconv"
)

func (c *checker) buildConnectionString() string {
	return c.host + ":" + strconv.Itoa(c.port)
}
