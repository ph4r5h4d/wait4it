package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"wait4it/internal/banner"

	"wait4it/pkg/check"
	"wait4it/pkg/model"
)

func defaultEnv(env string, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}

	return v
}

func defaultEnvInt(env string, def int) int {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	i, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		return def
	}
	return int(i)
}

func defaultEnvBool(env string, def bool) bool {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	if v == "true" || v == "1" {
		return true
	}
	return false
}

func main() {
	cfg := &model.CheckContext{}
	flag.StringVar(&cfg.Config.CheckType, "type", defaultEnv("W4IT_TYPE", "tcp"), "define the type of check")
	flag.IntVar(&cfg.Config.Timeout, "t", defaultEnvInt("W4IT_TIMEOUT", 30), "Timeout, amount of time wait4it waits for the port in seconds")
	flag.StringVar(&cfg.Host, "h", defaultEnv("W4IT_HOST", "127.0.0.1"), "IP of the host you want to test against")
	flag.IntVar(&cfg.Port, "p", defaultEnvInt("W4IT_PORT", 0), "Port")
	flag.StringVar(&cfg.Username, "u", defaultEnv("W4IT_USERNAME", ""), "Username of the service")
	flag.StringVar(&cfg.Password, "P", defaultEnv("W4IT_PASSWORD", ""), "Password of the service, it picks the W4IT_PASSWORD env if it is empty")
	flag.StringVar(&cfg.DatabaseName, "n", defaultEnv("W4IT_DBNAME", ""), "Name of the database")
	flag.StringVar(&cfg.DBConf.SSLMode, "ssl", defaultEnv("W4IT_SSL_MODE", "disable"), "Enable or Disable ssl mode (for some database or services)")
	flag.StringVar(&cfg.DBConf.OperationMode, "operation-mode", defaultEnv("W4IT_OPERATION_MODE", "standalone"), "choose operation mode (for some database or services)")
	flag.IntVar(&cfg.HttpConf.StatusCode, "status-code", defaultEnvInt("W4IT_HTTP_STATUS_CODE", 200), "Status code to be expected from http call")
	flag.StringVar(&cfg.HttpConf.Text, "http-text", defaultEnv("W4IT_HTTP_TEXT", ""), "Text to check inside http response")
	flag.BoolVar(&cfg.HttpConf.FollowRedirect, "http-follow-redirect", defaultEnvBool("W4IT_HTTP_FOLLOW_REDIRECT", true), "Whether to follow the redirect while doing the HTTP check")
	flag.StringVar(&cfg.KafkaConf.ConnectionType, "kafka-connection-type", defaultEnv("W4IT_KAFKA_CONNECTION_TYPE", "tcp"), "Kafka Connection Type, default is tcp")
	flag.BoolVar(&cfg.KafkaConf.ConnectToLeaderViaNonLeader, "kafka-connect-to-leader-via-non-leader", defaultEnvBool("W4IT_KAFKA_CONNECT_TO_LEADER_VIA_NON_LEADER", false), "Whether to connect to leader via non-leader, default is false")

	flag.Parse()
	// We don't want to show password in help message
	if cfg.Password == "" {
		defaultEnv("W4IT_PASSWORD", "")
	}
	cfg.Progress = func(s string) {
		fmt.Print(s)
	}
	banner.Banner()
	if err := check.RunCheck(context.Background(), cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Success!")
}
