package mongodb

import "strconv"

func (c *checker) buildConnectionString() string {
	if len(c.username) > 0 {
		return "mongodb://" + c.username + ":" + c.password + "@" + c.host + ":" + strconv.Itoa(c.port)
	}

	return "mongodb://" + c.host + ":" + strconv.Itoa(c.port)
}
