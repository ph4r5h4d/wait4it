package main

type checkContext struct {
	config configurationContext
	Host string
	Port int
	Username string
	Password string
	DatabaseName string
}

type configurationContext struct {
	checkType string
	timeout int
}