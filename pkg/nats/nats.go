package nats

import (
	"context"
	"errors"
	"net"
	"strings"

	natsgo "github.com/nats-io/nats.go"
	"wait4it/pkg/model"
)

type contextDialer struct {
	ctx context.Context
}

func (d *contextDialer) Dial(network, address string) (net.Conn, error) {
	var nd net.Dialer
	return nd.DialContext(d.ctx, network, address)
}

func (c *NatsChecker) BuildContext(cx model.CheckContext) {
	c.Host = cx.Host
	c.Port = cx.Port
	c.Username = cx.Username
	c.Password = cx.Password()
	c.Token = cx.NatsConf.Token
	c.TLS = tlsEnabled(cx.DBConf.SSLMode)
}

func (c *NatsChecker) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return errors.New("host should not be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return errors.New("invalid port range for nats")
	}

	return nil
}

func (c *NatsChecker) Check(ctx context.Context) (bool, bool, error) {
	opts := c.connectOptions()
	opts = append(opts, natsgo.SetCustomDialer(&contextDialer{ctx: ctx}))

	nc, err := natsgo.Connect(c.buildURL(), opts...)
	if err != nil {
		endOnError := isAuthError(err)
		return false, endOnError, err
	}
	defer nc.Close()

	if err := nc.Flush(); err != nil {
		endOnError := isAuthError(err)
		return false, endOnError, err
	}

	if !nc.IsConnected() {
		return false, false, errors.New("nats: not connected after connect")
	}

	return true, true, nil
}

// NewChecker creates a new NATS checker.
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &NatsChecker{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
