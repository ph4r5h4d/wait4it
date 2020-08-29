package envParser

import (
	"os"
	"strconv"
	"wait4it/model"
)

func Parse(ctx *model.CheckContext) *model.CheckContext {
	ctx.Config.CheckType = os.Getenv("W4IT_TYPE")
	ctx.Config.Timeout, _ = strconv.Atoi(os.Getenv("W4IT_TIMEOUT"))
	ctx.Host = os.Getenv("W4IT_HOST")
	ctx.Port, _ = strconv.Atoi(os.Getenv("W4IT_PORT"))
	ctx.Username = os.Getenv("W4IT_USERNAME")
	ctx.Password = os.Getenv("W4IT_PASSWORD")
	ctx.DatabaseName = os.Getenv("W4IT_DBNAME")
	ctx.DBConf.SSLMode = os.Getenv("W4IT_SSL_MODE")
	ctx.DBConf.OperationMode = os.Getenv("W4IT_OPERATION_MODE")
	ctx.HttpConf.StatusCode, _ = strconv.Atoi(os.Getenv("W4IT_HTTP_STATUS_CODE"))
	ctx.HttpConf.Text = os.Getenv("W4IT_HTTP_TEXT")
	return ctx
}
