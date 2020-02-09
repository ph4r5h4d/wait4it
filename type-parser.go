package main

import (
	"fmt"
	"wait4it/checkMySQL"
	"wait4it/checkTcp"
)

func (c checkContext) getStructByCheckType() interface{} {
	switch c.config.checkType {
	case "tcp":
		return buildIpCheckStruct(c)
	case "mysql":
		return buildMySQLCheckStruct(c)
	}

	return nil
}

func buildIpCheckStruct(c checkContext) checkTcp.IP {
	fmt.Println("here")
	return checkTcp.IP{
		Addr: c.Host,
		Port: c.Port,
	}
}

func buildMySQLCheckStruct(c checkContext) checkMySQL.MySQLConnection {
	return checkMySQL.MySQLConnection{
		Host:         c.Host,
		Port:         c.Port,
		Username:     c.Username,
		Password:     c.Password,
		DatabaseName: c.DatabaseName,
	}
}
