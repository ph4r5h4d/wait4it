package model

import (
	"os"
	"strings"
)

type CheckContext struct {
	Config        ConfigurationContext
	Host          string
	Port          int
	Username      string
	PasswordValue string
	PasswordFile  string
	DatabaseName  string
	DBConf        DatabaseSpecificConf
	HttpConf      HttpSpecificConf
	KafkaConf     KafkaConf

	Progress func(string)
}

func (cc CheckContext) Validate() (err error) {
	if cc.PasswordFile != "" {
		_, err = os.ReadFile(cc.PasswordFile)
	}
	return
}

func (cc CheckContext) Password() string {
	// assume higher precedence of password file over value
	if cc.PasswordFile != "" {
		// error checked during validation
		password, _ := os.ReadFile(cc.PasswordFile)
		return strings.Trim(string(password), "\n")
	}
	return cc.PasswordValue
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
	StatusCode     int
	Text           string
	FollowRedirect bool
}

type KafkaConf struct {
	ConnectionType              string
	ConnectToLeaderViaNonLeader bool
}
