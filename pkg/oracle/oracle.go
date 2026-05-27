package oracle

import (
	"context"
	"database/sql"
	"errors"
	"wait4it/pkg/model"

	_ "github.com/sijms/go-ora/v2"
)

func (o *OracleConnection) BuildContext(cx model.CheckContext) {
	o.Port = cx.Port
	o.Host = cx.Host
	o.Username = cx.Username
	o.Password = cx.Password()
	o.DatabaseName = cx.DatabaseName
}

func (o *OracleConnection) Validate() error {
	if len(o.Host) == 0 || len(o.Username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if o.Port < 1 || o.Port > 65535 {
		return errors.New("invalid port range for oracle")
	}

	return nil
}

func (o *OracleConnection) Check(ctx context.Context) (bool, bool, error) {
	dsl := o.BuildConnectionString()

	db, err := sql.Open("oracle", dsl)
	if err != nil {
		return false, true, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return false, false, nil
	}
	_ = db.Close()

	return true, true, nil
}
