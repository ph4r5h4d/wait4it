package redis

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"wait4it/pkg/model"

	"github.com/go-redis/redis/v8"
)

const (
	cluster    = "cluster"
	standalone = "standalone"
)

type checker struct {
	host          string
	port          int
	password      string
	database      int
	operationMode string
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
	c.password = cx.Password

	d, err := strconv.Atoi(cx.DatabaseName)
	if err != nil {
		d = 0
	}
	c.database = d

	switch cx.DBConf.OperationMode {
	case cluster:
		c.operationMode = cluster
	case standalone:
		c.operationMode = standalone
	default:
		c.operationMode = standalone
	}
}

func (c *checker) validate() error {
	if len(c.host) == 0 {
		return errors.New("host or username can't be empty")
	}

	if c.operationMode != cluster && c.operationMode != standalone {
		return errors.New("invalid operation mode")
	}

	if c.port < 1 || c.port > 65535 {
		return errors.New("invalid port range for redis")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	switch c.operationMode {
	case standalone:
		return c.checkStandAlone(ctx)
	case cluster:
		return c.checkCluster(ctx)
	default:
		return false, false, nil
	}
}

func (c *checker) checkStandAlone(ctx context.Context) (bool, bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.buildConnectionString(),
		Password: c.password, // no password set
		DB:       c.database, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return false, false, nil
	}

	_ = rdb.Close()

	return true, true, nil
}

func (c *checker) checkCluster(ctx context.Context) (bool, bool, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{c.buildConnectionString()}, // todo: add support for multiple hosts
		Password: c.password,
	})
	defer rdb.Close()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return false, false, nil
	}

	result, err := rdb.ClusterInfo(ctx).Result()
	if err != nil {
		return false, false, nil
	}

	if result != "" {
		if !strings.Contains(result, "cluster_state:ok") {
			return false, false, errors.New("cluster is not healthy")
		}
	}

	return true, true, nil
}

func (c *checker) buildConnectionString() string {
	return c.host + ":" + strconv.Itoa(c.port)
}
