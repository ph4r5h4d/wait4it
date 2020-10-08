package rabbitmq

import (
	"errors"
	"fmt"
	"wait4it/model"

	"github.com/streadway/amqp"
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
	rc.Password = cx.Password
}

func (rc *RabbitChecker) Validate() (bool, error) {
	if rc.Host == "" {
		return false, errors.New("Host should not be empty")
	}

	if rc.Username == "" {
		return false, errors.New("Username should not be empty")
	}

	if rc.Port == 0 {
		return false, errors.New("Port should not be empty")
	}

	rc.conString = fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.Username, rc.Password, rc.Host, rc.Port)

	return true, nil
}

func (rc *RabbitChecker) Check() (bool, bool, error) {
	con, err := amqp.Dial(rc.conString)

	if err != nil {
		return false, false, err
	}

	ch, err := con.Channel()
	if err != nil {
		return false, false, err
	}
	defer ch.Close()

	return true, false, nil
}
