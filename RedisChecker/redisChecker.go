package RedisChecker

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"wait4it/model"
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

	if cx.DBConf.ClusterMode == "enable" {
		m.IsClustered = true
	} else {
		m.IsClustered = false
	}
}

func (m *RedisConnection) Validate() (bool, error) {
	if len(m.Host) == 0 {
		return false, errors.New("host or username can't be empty")
	}

	if m.Port < 0 || m.Port > 65535 {
		return false, errors.New("invalid port range for mysql")
	}

	return true, nil
}

func (m *RedisConnection) Check() (bool, bool, error) {
	ctx := context.Background()

	if !m.IsClustered {
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
	} else {
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              []string{m.BuildConnectionString()},
			Password:           m.Password,
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
}
