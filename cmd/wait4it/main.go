package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"wait4it/pkg/check"
	"wait4it/pkg/model"

	"github.com/fzerorubigd/clictx"
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

func main() {
	cfg := &model.CheckContext{}
	flag.StringVar(&cfg.Config.CheckType, "type", defaultEnv("W4IT_TYPE", "tcp"), "define the type of check")
	flag.IntVar(&cfg.Config.Timeout, "t", defaultEnvInt("W4IT_TIMEOUT", 30), "Timeout, amount of time wait4it waits for the port in seconds")
	flag.StringVar(&cfg.Host, "h", defaultEnv("W4IT_HOST", "127.0.0.1"), "IP of the host you want to test against")
	flag.IntVar(&cfg.Port, "p", defaultEnvInt("W4IT_PORT", 0), "Port")
	flag.StringVar(&cfg.Username, "u", defaultEnv("W4IT_USERNAME", ""), "Username of the service")
	flag.StringVar(&cfg.Password, "P", "", "Password of the service, it picks the W4IT_PASSWORD env if it is empty")
	flag.StringVar(&cfg.DatabaseName, "n", defaultEnv("W4IT_DBNAME", ""), "Name of the database")
	flag.StringVar(&cfg.DBConf.SSLMode, "ssl", defaultEnv("W4IT_SSL_MODE", "disable"), "Enable or Disable ssl mode (for some database or services)")
	flag.StringVar(&cfg.DBConf.OperationMode, "operation-mode", defaultEnv("W4IT_OPERATION_MODE", "standalone"), "choose operation mode (for some database or services)")
	flag.IntVar(&cfg.HttpConf.StatusCode, "status-code", defaultEnvInt("W4IT_HTTP_STATUS_CODE", 200), "Status code to be expected from http call")
	flag.StringVar(&cfg.HttpConf.Text, "http-text", defaultEnv("W4IT_HTTP_TEXT", ""), "Text to check inside http response")

	flag.Parse()
	// We don't want to show password in help message
	if cfg.Password == "" {
		defaultEnv("W4IT_PASSWORD", "")
	}
	cfg.Progress = func(s string) {
		fmt.Print(s)
	}
	if err := check.RunCheck(clictx.DefaultContext(), cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Success!")
}
