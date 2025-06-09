package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"wait4it/pkg/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitChecker is the struct to check for a RabbitMQ instance
type RabbitChecker struct {
	Host     string
	Port     int
	Username string
	Password string

	conString string
}

func (rc *RabbitChecker) BuildContext(cx model.CheckContext) {
	rc.Host = cx.Host
	rc.Port = cx.Port
	rc.Username = cx.Username
	rc.Password = cx.Password()
}

func (rc *RabbitChecker) Validate() error {
	if rc.Host == "" {
		return errors.New("Host should not be empty")
	}

	if rc.Username == "" {
		return errors.New("Username should not be empty")
	}

	if rc.Port == 0 {
		return errors.New("Port should not be empty")
	}

	rc.conString = fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.Username, rc.Password, rc.Host, rc.Port)

	return nil
}

func (rc *RabbitChecker) Check(ctx context.Context) (bool, bool, error) {
	con, err := amqp.DialConfig(rc.conString, amqp.Config{
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

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &RabbitChecker{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
