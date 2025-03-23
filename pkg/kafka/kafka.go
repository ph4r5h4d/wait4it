package kafka

import (
	"context"
	"errors"
	KafkaGo "github.com/segmentio/kafka-go"
	"net"
	"strconv"
	"wait4it/pkg/model"
)

const WaitTimeOutSeconds = 2

func (c *KafkaConnection) BuildContext(cx model.CheckContext) {
	c.Port = cx.Port
	c.Host = cx.Host
	c.ConnectionType = cx.KafkaConf.ConnectionType
	c.ConnectToLeaderViaNonLeader = cx.KafkaConf.ConnectToLeaderViaNonLeader
}

func (c *KafkaConnection) Validate() error {
	if len(c.Host) == 0 {
		return errors.New("host can't be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}

func (c *KafkaConnection) Check(ctx context.Context) (bool, bool, error) {
	conn, err := KafkaGo.Dial(c.ConnectionType, net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
	if err != nil {
		return false, false, err
	}
	defer conn.Close()

	// check if the operation is successful, then we can assume this is a valid connection
	_, err = conn.ReadPartitions()
	if err != nil {
		return false, false, err
	}

	// this part should work, but I have not yet fully tested it.
	if c.ConnectToLeaderViaNonLeader {
		controller, err := conn.Controller()
		if err != nil {
			return false, false, err
		}
		var connLeader *KafkaGo.Conn
		connLeader, err = KafkaGo.Dial(c.ConnectionType, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			return false, true, err
		}
		defer connLeader.Close()
	}

	return true, true, nil
}
