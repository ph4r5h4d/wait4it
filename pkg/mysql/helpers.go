package mysql

import (
	"strconv"
)

func (c *checker) buildConnectionString() string {
	dsl := ""

	if len(c.password) == 0 {
		dsl = dsl + c.username
	} else {
		dsl = dsl + c.username + ":" + c.password
	}

	dsl = dsl + "@tcp(" + c.host + ":" + strconv.Itoa(c.port) + ")/"

	if len(c.databaseName) > 0 {
		dsl = dsl + c.databaseName
	}

	return dsl
}
