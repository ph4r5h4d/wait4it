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

	Progress func(string)
}

type ConfigurationContext struct {
	CheckType string
	Timeout   int
}

type DatabaseSpecificConf struct {
	SSLMode       string
	OperationMode string
}
type HttpSpecificConf struct {
	StatusCode int
	Text       string
}
