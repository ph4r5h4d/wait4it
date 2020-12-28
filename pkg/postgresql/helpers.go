package postgresql

import (
	"fmt"
)

func (c *checker) buildConnectionString() string {
	dsl := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=%s ",
		c.host, c.port, c.username, c.password, c.sslMode, c.databaseName)

	return dsl
}
