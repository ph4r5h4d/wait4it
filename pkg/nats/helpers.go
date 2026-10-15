package nats

import (
	"errors"
	"fmt"
	"strings"

	natsgo "github.com/nats-io/nats.go"
)

func tlsEnabled(sslMode string) bool {
	switch strings.ToLower(sslMode) {
	case "enable", "require", "true", "1":
		return true
	default:
		return false
	}
}

func (c *NatsChecker) buildURL() string {
	scheme := "nats"
	if c.TLS {
		scheme = "tls"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, c.Host, c.Port)
}

func (c *NatsChecker) connectOptions() []natsgo.Option {
	var opts []natsgo.Option
	opts = append(opts, natsgo.NoReconnect())

	if c.Token != "" {
		opts = append(opts, natsgo.Token(c.Token))
	}
	if c.Username != "" {
		opts = append(opts, natsgo.UserInfo(c.Username, c.Password))
	}
	if c.TLS {
		opts = append(opts, natsgo.Secure())
	}

	return opts
}

func isAuthError(err error) bool {
	return errors.Is(err, natsgo.ErrAuthorization) ||
		errors.Is(err, natsgo.ErrAuthExpired) ||
		errors.Is(err, natsgo.ErrAuthRevoked) ||
		errors.Is(err, natsgo.ErrAccountAuthExpired)
}
