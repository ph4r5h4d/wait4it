package mysql

type MySQLConnection struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}
