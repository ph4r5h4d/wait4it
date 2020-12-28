package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"wait4it/pkg/model"

	_ "github.com/lib/pq"
)

type checker struct {
	host         string
	port         int
	username     string
	password     string
	databaseName string
	sslMode      string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	checker := &checker{}
	checker.buildContext(*c)
	if err := checker.validate(); err != nil {
		return nil, err
	}

	return checker, nil
}

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

func (c *checker) buildConnectionString() string {
	dsl := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=%s ",
		c.host, c.port, c.username, c.password, c.sslMode, c.databaseName)

	return dsl
}
