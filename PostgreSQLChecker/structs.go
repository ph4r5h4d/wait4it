package PostgreSQLChecker

type PostgresSQLConnection struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}
