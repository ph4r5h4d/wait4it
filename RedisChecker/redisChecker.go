package RedisChecker

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
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
