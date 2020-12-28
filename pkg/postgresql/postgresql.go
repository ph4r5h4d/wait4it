package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"wait4it/pkg/model"

	_ "github.com/lib/pq"
)

func (c *checker) buildContext(cx model.CheckContext) {
	c.port = cx.Port
	c.host = cx.Host
	c.username = cx.Username
	c.password = cx.Password
	c.databaseName = cx.DatabaseName

	if len(cx.DBConf.SSLMode) < 1 {
		c.sslMode = "disable"
	} else {
		c.sslMode = cx.DBConf.SSLMode
	}
}

func (c *checker) validate() error {
	if len(c.host) == 0 || len(c.username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("invalid port range for PostgresSQL")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	dsl := c.buildConnectionString()

	db, err := sql.Open("postgres", dsl)
	if err != nil {
		return false, true, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return false, false, nil
	}
	_ = db.Close()

	return true, true, nil
}
