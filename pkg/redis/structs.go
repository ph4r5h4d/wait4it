package redis

type RedisConnection struct {
	Host          string
	Port          int
	Password      string
	Database      int
	OperationMode string
}
