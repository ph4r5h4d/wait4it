package mysql

import (
	"strconv"
)

func (m MySQLConnection) BuildConnectionString() string {
	dsl := ""

	if len(m.Password) == 0 {
		dsl = dsl + m.Username
	} else {
		dsl = dsl + m.Username + ":" + m.Password
	}

	dsl = dsl + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/"

	if len(m.DatabaseName) > 0 {
		dsl = dsl + m.DatabaseName
	}

	return dsl
}
