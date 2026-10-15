package nats

// NatsChecker waits for a NATS server to accept client connections.
type NatsChecker struct {
	Host     string
	Port     int
	Username string
	Password string
	Token    string
	TLS      bool
}
