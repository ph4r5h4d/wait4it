package RedisChecker

import (
	"strconv"
)

func (m RedisConnection) BuildConnectionString() string {
	return m.Host + ":" + strconv.Itoa(m.Port)
}
