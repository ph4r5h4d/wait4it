package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"wait4it/tcpcheck"
)

func main() {
	ip := flag.String("h", "127.0.0.1", "IP of the host you want to test against")
	port := flag.Int("p", 0, "Port")
	timeout := flag.Int("t", 30, "Timeout, amount of time wait4it waits for the port in seconds")
	flag.Parse()

	i := tcpcheck.IP{
		Addr: *ip,
		Port: *port,
	}
	fmt.Println(fmt.Sprintf("Starting to check on %s on port %d with %d as timeout", *ip, *port, *timeout))
	fmt.Print("Wait4it...")

	t := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go ticker(i, t, done)

	time.Sleep(time.Duration(*timeout) * time.Second)
	done <- true

	fmt.Print("failed")
	os.Exit(1)
}

func ticker(ip tcpcheck.IP, t *time.Ticker, d chan bool) {
	for {
		select {
		case <-d:
			return
		case <-t.C:
			check(ip)
		}
	}
}

func check(ip tcpcheck.IP) {
	r, err := ip.DoesPortAcceptConnection()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if r {
		fmt.Println("succeed")
		os.Exit(0)
	} else {
		fmt.Print(".")
	}
}
