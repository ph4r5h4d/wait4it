package tcp

func (tcp *check) isPortInValidRange() bool {
	if tcp.port < minPort || tcp.port > maxPort {
		return false
	}
	return true
}
