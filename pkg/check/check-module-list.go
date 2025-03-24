package check

import (
	"wait4it/pkg/aerospike"
	"wait4it/pkg/elasticsearch"
	"wait4it/pkg/http"
	"wait4it/pkg/kafka"
	"wait4it/pkg/memcached"
	"wait4it/pkg/model"
	"wait4it/pkg/mongodb"
	"wait4it/pkg/mysql"
	"wait4it/pkg/postgresql"
	"wait4it/pkg/rabbitmq"
	"wait4it/pkg/redis"
	"wait4it/pkg/tcp"
)

var cm = map[string]func(c *model.CheckContext) (model.CheckInterface, error){
	"tcp":           tcp.NewChecker,
	"mysql":         mysql.NewChecker,
	"postgres":      postgresql.NewChecker,
	"http":          http.NewChecker,
	"mongo":         mongodb.NewChecker,
	"redis":         redis.NewChecker,
	"rabbitmq":      rabbitmq.NewChecker,
	"memcached":     memcached.NewChecker,
	"elasticsearch": elasticsearch.NewChecker,
	"aerospike":     aerospike.NewChecker,
	"kafka":         kafka.NewChecker,
}
