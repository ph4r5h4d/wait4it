package flagParser

import (
	"flag"
	"wait4it/model"
)

func Parse(ctx *model.CheckContext) *model.CheckContext {
	checkType := flag.String("type", "tcp", "define the type of check")
	timeout := flag.Int("t", 30, "Timeout, amount of time wait4it waits for the port in seconds")
	host := flag.String("h", "127.0.0.1", "IP of the host you want to test against")
	port := flag.Int("p", 0, "Port")
	username := flag.String("u", "", "Username of the service")
	password := flag.String("P", "", "Password of the service")
	databaseName := flag.String("n", "", "Name of the database")
	sslMode := flag.String("ssl", "disable", "Enable or Disable ssl mode (for some database or services)")
	operationMode := flag.String("operation-mode", "standalone", "choose operation mode (for some database or services)")
	statusCode := flag.Int("status-code", 200, "Status code to be expected from http call")
	text := flag.String("http-text", "", "Text to check inside http response")

	flag.Parse()

	if ctx.Config.CheckType == "" {
		ctx.Config.CheckType = *checkType
	}

	if ctx.Config.Timeout == 0 {
		ctx.Config.Timeout = *timeout
	}

	if ctx.Host == "" {
		ctx.Host = *host
	}

	if ctx.Port == 0 {
		ctx.Port = *port
	}

	if ctx.Username == "" {
		ctx.Username = *username
	}

	if ctx.Password == "" {
		ctx.Password = *password
	}

	if ctx.DatabaseName == "" {
		ctx.DatabaseName = *databaseName
	}

	if ctx.DBConf.SSLMode == "" {
		ctx.DBConf.SSLMode = *sslMode
	}

	if ctx.DBConf.OperationMode == "" {
		ctx.DBConf.OperationMode = *operationMode
	}

	if ctx.HttpConf.StatusCode == 0 {
		ctx.HttpConf.StatusCode = *statusCode
	}

	if ctx.HttpConf.Text == "" {
		ctx.HttpConf.Text = *text
	}

	return ctx
}
