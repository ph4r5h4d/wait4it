package PostgreSQLChecker

import (
	"fmt"
)

func (pq PostgresSQLConnection) BuildConnectionString() string {
	dsl := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=%s ",
		pq.Host, pq.Port, pq.Username, pq.Password, pq.SSLMode, pq.DatabaseName)

	return dsl
}
