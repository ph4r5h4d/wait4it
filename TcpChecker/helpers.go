package TcpChecker

func (tcp *Tcp) isPortInValidRange() bool {
	if tcp.Port < minPort || tcp.Port > maxPort {
		return false
	}
	return true
}
