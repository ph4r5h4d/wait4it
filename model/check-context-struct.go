package model

type CheckContext struct {
	Config       ConfigurationContext
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	DBConf       DatabaseSpecificConf
}

type ConfigurationContext struct {
	CheckType string
	Timeout   int
}

type DatabaseSpecificConf struct {
	SSLMode string
}
