package model

import (
	"log"
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
		password, err := os.ReadFile(cc.PasswordFile)
		if err != nil {
			log.Println("Failed to load password from file", "file", cc.PasswordFile, "err", err)
		}
		return strings.TrimSpace(string(password))
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
