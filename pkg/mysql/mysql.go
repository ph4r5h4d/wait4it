package mysql

import (
	"context"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"strconv"

	"wait4it/pkg/model"

	"github.com/go-sql-driver/mysql"
)

type checker struct {
	host         string
	port         int
	username     string
	password     string
	databaseName string
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
