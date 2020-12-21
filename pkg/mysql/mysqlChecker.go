package mysql

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"wait4it/pkg/model"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func (m *MySQLConnection) BuildContext(cx model.CheckContext) {
	m.Port = cx.Port
	m.Host = cx.Host
	m.Username = cx.Username
	m.Password = cx.Password
	m.DatabaseName = cx.DatabaseName
}

func (m *MySQLConnection) Validate() error {
	if len(m.Host) == 0 || len(m.Username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if m.Port < 1 || m.Port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}

func (m *MySQLConnection) Check() (bool, bool, error) {
	dsl := m.BuildConnectionString()

	db, err := sql.Open("mysql", dsl)

	// if there is an error opening the connection, handle it
	if err != nil {
		return false, true, err
	}

	err = mysql.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags))
	if err != nil {
		return false, true, err
	}

	err = db.Ping()
	if err != nil {
		// todo: need a logger
		return false, false, nil
	}
	_ = db.Close()

	return true, true, nil
}
