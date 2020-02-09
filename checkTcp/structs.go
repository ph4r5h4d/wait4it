package checkTcp

const (
	minPort = 1
	maxPort = 65535
)

type IP struct {
	Addr string
	Port int
}
