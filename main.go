package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	ct := flag.String("type", "tcp", "define the type of check, currently [tcp,mysql] are supported")
	timeout := flag.Int("t", 30, "Timeout, amount of time wait4it waits for the port in seconds")
	host := flag.String("h", "127.0.0.1", "IP of the host you want to test against")
	port := flag.Int("p", 0, "Port")
	username := flag.String("u", "", "Username of the service")
	password := flag.String("P", "", "Password of the service")
	databaseName := flag.String("n", "", "Name of the database")

	flag.Parse()

	c := checkContext{
		config: configurationContext{
			checkType: *ct,
			timeout:   *timeout,
		},
		Host:         *host,
		Port:         *port,
		Username:     *username,
		Password:     *password,
		DatabaseName: *databaseName,
	}
	cs := c.getStructByCheckType()

	fmt.Print("Wait4it...")

	t := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go ticker(cs, c.config.checkType, t, done)

	time.Sleep(time.Duration(*timeout) * time.Second)
	done <- true

	fmt.Print("failed")
	os.Exit(1)
}
