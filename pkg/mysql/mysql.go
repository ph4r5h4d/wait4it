package mysql

import (
	"context"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"

	"wait4it/pkg/model"

	"github.com/go-sql-driver/mysql"
)

func (c *checker) buildContext(cx model.CheckContext) {
	c.port = cx.Port
	c.host = cx.Host
	c.username = cx.Username
	c.password = cx.Password
	c.databaseName = cx.DatabaseName
}

func (c *checker) validate() error {
	if len(c.host) == 0 || len(c.username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	dsl := c.buildConnectionString()

	db, err := sql.Open("mysql", dsl)
	if err != nil {
		return false, true, err
	}

	err = mysql.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags))
	if err != nil {
		return false, true, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		// todo: need a logger
		return false, false, nil
	}
	_ = db.Close()

	return true, true, nil
}
