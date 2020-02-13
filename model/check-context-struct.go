package model

type CheckContext struct {
	Config       ConfigurationContext
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

type ConfigurationContext struct {
	CheckType string
	Timeout   int
}
