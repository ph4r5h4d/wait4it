package TcpChecker

const (
	minPort = 1
	maxPort = 65535
)

type Tcp struct {
	Addr string
	Port int
}
