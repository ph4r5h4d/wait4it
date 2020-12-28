package elasticsearch

import (
	"context"
	"errors"
	"strconv"
	"wait4it/pkg/model"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

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
	c.host = cx.Host
	c.port = cx.Port
	c.username = cx.Username
	c.password = cx.Password
}

func (c *checker) validate() error {
	if len(c.host) == 0 {
		return errors.New("Host can't be empty")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("Invalid port range for ElasticSearch")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			c.buildConnectionString(),
		},
		Username: c.username,
		Password: c.password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return false, true, err
	}

	if _, err := es.Ping(es.Ping.WithContext(ctx)); err != nil {
		return false, false, err
	}

	return true, true, nil
}

func (c *checker) buildConnectionString() string {
	return c.host + ":" + strconv.Itoa(c.port)
}
