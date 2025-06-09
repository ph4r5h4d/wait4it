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
	Cluster    = "cluster"
	Standalone = "standalone"
)

func (m *RedisConnection) BuildContext(cx model.CheckContext) {
	m.Host = cx.Host
	m.Port = cx.Port
	m.Password = cx.Password()

	d, err := strconv.Atoi(cx.DatabaseName)
	if err != nil {
		d = 0
	}
	m.Database = d

	switch cx.DBConf.OperationMode {
	case Cluster:
		m.OperationMode = Cluster
	case Standalone:
		m.OperationMode = Standalone
	default:
		m.OperationMode = Standalone
	}
}

func (m *RedisConnection) Validate() error {
	if len(m.Host) == 0 {
		return errors.New("host or username can't be empty")
	}

	if m.OperationMode != Cluster && m.OperationMode != Standalone {
		return errors.New("invalid operation mode")
	}

	if m.Port < 1 || m.Port > 65535 {
		return errors.New("invalid port range for redis")
	}

	return nil
}

func (m *RedisConnection) Check(ctx context.Context) (bool, bool, error) {
	switch m.OperationMode {
	case Standalone:
		return m.checkStandAlone(ctx)
	case Cluster:
		return m.checkCluster(ctx)
	default:
		return false, false, nil
	}
}

func (m *RedisConnection) checkStandAlone(ctx context.Context) (bool, bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     m.BuildConnectionString(),
		Password: m.Password, // no password set
		DB:       m.Database, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return false, false, nil
	}

	_ = rdb.Close()

	return true, true, nil
}

func (m *RedisConnection) checkCluster(ctx context.Context) (bool, bool, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{m.BuildConnectionString()}, //todo: add support for multiple hosts
		Password: m.Password,
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
