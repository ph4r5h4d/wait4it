package main

import (
	"os"
	"time"
	"wait4it/checkMySQL"
	"wait4it/checkTcp"
)

func ticker(cs interface{}, ct string, t *time.Ticker, d chan bool) {
	for {
		select {
		case <-d:
			return
		case <-t.C:
			check(cs, ct)
		}
	}
}

func check(cs interface{}, ct string) {
	switch ct {
	case "tcp":
		tcpCheck(cs.(checkTcp.IP))
	case "mysql":
		mysqlCheck(cs.(checkMySQL.MySQLConnection))
	}
}

func tcpCheck(ip checkTcp.IP) {
	r, err := ip.DoesPortAcceptConnection()

	if err != nil {
		wStdErr(err.Error())
		os.Exit(2)
	}

	wStdOut(r)
}

func mysqlCheck(c checkMySQL.MySQLConnection) {
	r, err := c.IsMySQLAvailable()

	if err != nil {
		wStdErr(err.Error())
		os.Exit(2)
	}

	wStdOut(r)
}
