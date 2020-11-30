package AerospikeChecker

import (
	"errors"
	"github.com/aerospike/aerospike-client-go"
	"wait4it/model"
)

func (m *AerospikeConnection) BuildContext(cx model.CheckContext) {
	m.Host = cx.Host
	m.Port = cx.Port
}

func (m *AerospikeConnection) Validate() error {
	if len(m.Host) == 0 {
		return errors.New("host can't be empty")
	}

	if m.Port < 1 || m.Port > 65535 {
		return errors.New("invalid port range for Memcached")
	}

	return nil
}

func (m *AerospikeConnection) Check() (bool, bool, error) {
	c, err := aerospike.NewClient(m.Host, m.Port)

	if err != nil {
		return false, false, err
	}

	if !c.IsConnected() {
		return false, false, errors.New("client is not connected")
	}

	return true, true, nil
}
