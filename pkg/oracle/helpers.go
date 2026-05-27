package oracle

import (
	go_ora "github.com/sijms/go-ora/v2"
)

func (o OracleConnection) BuildConnectionString() string {
	return go_ora.BuildUrl(o.Host, o.Port, o.DatabaseName, o.Username, o.Password, nil)
}
