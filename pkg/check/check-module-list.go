package check

import (
	"wait4it/pkg/aerospike"
	"wait4it/pkg/elasticsearch"
	"wait4it/pkg/http"
	"wait4it/pkg/memcached"
	"wait4it/pkg/model"
	"wait4it/pkg/mongodb"
	"wait4it/pkg/mysql"
	"wait4it/pkg/postgresql"
	"wait4it/pkg/rabbitmq"
	"wait4it/pkg/redis"
	"wait4it/pkg/tcp"
)

var cm = map[string]model.CheckInterface{
	"tcp":           &tcp.Tcp{},
	"mysql":         &mysql.MySQLConnection{},
	"postgres":      &postgresql.PostgresSQLConnection{},
	"http":          &http.HttpCheck{},
	"mongo":         &mongodb.MongoDbConnection{},
	"redis":         &redis.RedisConnection{},
	"rabbitmq":      &rabbitmq.RabbitChecker{},
	"memcached":     &memcached.MemcachedConnection{},
	"elasticsearch": &elasticsearch.ElasticSearchChecker{},
	"aerospike":     &aerospike.AerospikeConnection{},
}
