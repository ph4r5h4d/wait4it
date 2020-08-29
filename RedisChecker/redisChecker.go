package RedisChecker

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"wait4it/model"
)

const (
	ClusterMode    = "cluster"
	StandAloneMode = "standalone"
)

func (m *RedisConnection) BuildContext(cx model.CheckContext) {
	m.Host = cx.Host
	m.Port = cx.Port
	m.Password = cx.Password

	d, err := strconv.Atoi(cx.DatabaseName)
	if err != nil {
		d = 0
	}
	m.Database = d

	switch cx.DBConf.OperationMode {
	case ClusterMode:
		m.OperationMode = ClusterMode
	case StandAloneMode:
		m.OperationMode = StandAloneMode
	default:
		m.OperationMode = StandAloneMode
	}
}

func (m *RedisConnection) Validate() (bool, error) {
	if len(m.Host) == 0 {
		return false, errors.New("host or username can't be empty")
	}

	if m.OperationMode != ClusterMode && m.OperationMode != StandAloneMode {
		return false, errors.New("operation mode is invalid, 'cluster' and 'standalone' modes are valid")
	}

	if m.Port < 0 || m.Port > 65535 {
		return false, errors.New("invalid port range for redis")
	}

	return true, nil
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
		Addrs:    []string{m.BuildConnectionString()},
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
			return false, false, errors.New("can't find cluster_state:ok in response")
		}
	}

	return true, false, nil
}

func (m *RedisConnection) Check() (bool, bool, error) {
	ctx := context.Background()

	switch m.OperationMode {
	case StandAloneMode:
		return m.checkStandAlone(ctx)
	case ClusterMode:
		return m.checkCluster(ctx)
	default:
		return false, false, errors.New("operation mode is invalid")
	}
}
