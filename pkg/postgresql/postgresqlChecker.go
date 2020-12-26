package postgresql

import (
	"database/sql"
	"errors"
	"wait4it/pkg/model"

	_ "github.com/lib/pq"
)

func (pq *PostgresSQLConnection) BuildContext(cx model.CheckContext) {
	pq.Port = cx.Port
	pq.Host = cx.Host
	pq.Username = cx.Username
	pq.Password = cx.Password
	pq.DatabaseName = cx.DatabaseName
	if len(cx.DBConf.SSLMode) < 1 {
		pq.SSLMode = "disable"
	} else {
		pq.SSLMode = cx.DBConf.SSLMode
	}
}

func (pq *PostgresSQLConnection) Validate() error {
	if len(pq.Host) == 0 || len(pq.Username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if pq.Port < 1 || pq.Port > 65535 {
		return errors.New("invalid port range for PostgresSQL")
	}

	return nil
}

func (pq *PostgresSQLConnection) Check() (bool, bool, error) {
	dsl := pq.BuildConnectionString()

	db, err := sql.Open("postgres", dsl)

	// if there is an error opening the connection, handle it
	if err != nil {
		return false, true, err
	}

	err = db.Ping()
	if err != nil {
		return false, false, nil
	}
	_ = db.Close()

	return true, true, nil
}
