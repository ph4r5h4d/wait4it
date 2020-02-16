package cmd

import (
	"wait4it/MySQLChecker"
	"wait4it/PostgreSQLChecker"
	"wait4it/TcpChecker"
)

var cm = map[string]interface{}{
	"tcp":      &TcpChecker.Tcp{},
	"mysql":    &MySQLChecker.MySQLConnection{},
	"postgres": &PostgreSQLChecker.PostgresSQLConnection{},
}
