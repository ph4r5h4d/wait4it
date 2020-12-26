package main

import (
	"flag"
	"os"
	"strconv"
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

func main() {
	ctx := &model.CheckContext{}
	flag.StringVar(&ctx.Config.CheckType, "type", defaultEnv("W4IT_TYPE", "tcp"), "define the type of check")
	flag.IntVar(&ctx.Config.Timeout, "t", defaultEnvInt("W4IT_TIMEOUT", 30), "Timeout, amount of time wait4it waits for the port in seconds")
	flag.StringVar(&ctx.Host, "h", defaultEnv("W4IT_HOST", "127.0.0.1"), "IP of the host you want to test against")
	flag.IntVar(&ctx.Port, "p", defaultEnvInt("W4IT_PORT", 0), "Port")
	flag.StringVar(&ctx.Username, "u", defaultEnv("W4IT_USERNAME", ""), "Username of the service")
	flag.StringVar(&ctx.Password, "P", "", "Password of the service, it picks the W4IT_PASSWORD env if it is empty")
	flag.StringVar(&ctx.DatabaseName, "n", defaultEnv("W4IT_DBNAME", ""), "Name of the database")
	flag.StringVar(&ctx.DBConf.SSLMode, "ssl", defaultEnv("W4IT_DBNAME", "disable"), "Enable or Disable ssl mode (for some database or services)")
	flag.StringVar(&ctx.DBConf.OperationMode, "operation-mode", defaultEnv("W4IT_SSL_MODE", "standalone"), "choose operation mode (for some database or services)")
	flag.IntVar(&ctx.HttpConf.StatusCode, "status-code", defaultEnvInt("W4IT_OPERATION_MODE", 200), "Status code to be expected from http call")
	flag.StringVar(&ctx.HttpConf.Text, "http-text", defaultEnv("W4IT_HTTP_TEXT", ""), "Text to check inside http response")

	flag.Parse()
	// We don't want to show password in help message
	if ctx.Password == "" {
		defaultEnv("W4IT_PASSWORD", "")
	}
	check.RunCheck(*ctx)
}
