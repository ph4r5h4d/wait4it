package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"wait4it/pkg/model"

	"github.com/streadway/amqp"
)

type checker struct {
	host     string
	port     int
	username string
	password string

	conString string
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	checker := &checker{}
	checker.buildContext(*c)
	if err := checker.validate(); err != nil {
		return nil, err
	}

	return checker, nil
}

func (c *checker) buildContext(cx model.CheckContext) {
	c.host = cx.Host
	c.port = cx.Port
	c.username = cx.Username
	c.password = cx.Password
}

func (c *checker) validate() error {
	if c.host == "" {
		return errors.New("Host should not be empty")
	}

	if c.username == "" {
		return errors.New("Username should not be empty")
	}

	if c.port == 0 {
		return errors.New("Port should not be empty")
	}

	c.conString = fmt.Sprintf("amqp://%s:%s@%s:%d/", c.username, c.password, c.host, c.port)

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	con, err := amqp.DialConfig(c.conString, amqp.Config{
		Heartbeat: time.Second * 10,
		Locale:    "en_US",
		Dial: func(network, addr string) (net.Conn, error) {
			var d net.Dialer
			conn, err := d.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			// Heartbeating hasn't started yet, don't stall forever on a dead server.
			// A deadline is set for TLS and AMQP handshaking. After AMQP is established,
			// the deadline is cleared in openComplete.
			if err := conn.SetDeadline(time.Now().Add(time.Second * 30)); err != nil {
				return nil, err
			}

			return conn, nil
		},
	})

	if err != nil {
		return false, false, err
	}
	defer con.Close()

	ch, err := con.Channel()
	if err != nil {
		return false, false, err
	}
	defer ch.Close()

	return true, false, nil
}
