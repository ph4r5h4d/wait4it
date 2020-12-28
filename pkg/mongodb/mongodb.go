package mongodb

import (
	"context"
	"errors"
	"strconv"

	"wait4it/pkg/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const WaitTimeOutSeconds = 2

type checker struct {
	host     string
	port     int
	username string
	password string
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
}

func (c *checker) validate() error {
	if len(c.host) == 0 {
		return errors.New("host can't be empty")
	}

	if len(c.username) > 0 && len(c.password) == 0 {
		return errors.New("password can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(c.buildConnectionString()))
	if err != nil {
		return false, true, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return false, true, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false, false, err
	}

	return true, true, nil
}

func (c *checker) buildConnectionString() string {
	if len(c.username) > 0 {
		return "mongodb://" + c.username + ":" + c.password + "@" + c.host + ":" + strconv.Itoa(c.port)
	}

	return "mongodb://" + c.host + ":" + strconv.Itoa(c.port)
}
