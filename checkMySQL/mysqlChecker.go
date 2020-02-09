package checkMySQL

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func (m MySQLConnection) IsMySQLAvailable() (bool, error) {
	err := m.validateStruct()

	if err != nil {
		return false, err
	}

	dsl := m.BuildConnectionString()

	db, err := sql.Open("mysql", dsl)

	// if there is an error opening the connection, handle it
	if err != nil {
		return false, err
	}

	err = db.Ping()
	if err != nil {
		// todo: need a logger
		return false, nil
	}

	_ = db.Close()
	return true, nil
}
