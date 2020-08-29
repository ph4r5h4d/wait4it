package model

type CheckContext struct {
	Config       ConfigurationContext
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	DBConf       DatabaseSpecificConf
	HttpConf     HttpSpecificConf
}

type ConfigurationContext struct {
	CheckType string
	Timeout   int
}

type DatabaseSpecificConf struct {
	SSLMode     string
	ClusterMode string
}
type HttpSpecificConf struct {
	StatusCode int
	Text       string
}
