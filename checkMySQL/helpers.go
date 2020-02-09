package checkMySQL

import (
	"errors"
	"strconv"
)

func (m MySQLConnection) BuildConnectionString() string {
	dsl := ""

	if len(m.Password) == 0 {
		dsl = dsl + m.Username
	}else{
		dsl = dsl + m.Username + ":" + m.Password
	}

	dsl = dsl + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/"

	if len(m.DatabaseName) > 0 {
		dsl = dsl + m.DatabaseName
	}

	return dsl
}

func (m MySQLConnection) validateStruct()  error{
	if len(m.Host) == 0 || len(m.Username) == 0 {
		return errors.New("host or username can't be empty")
	}

	if m.Port < 0 || m.Port > 65535 {
		return errors.New("invalid port range for mysql")
	}

	return nil
}


