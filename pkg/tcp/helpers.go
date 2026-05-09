package tcp

func (tcp *Check) IsPortInValidRange() bool {
	if tcp.Port < minPort || tcp.Port > maxPort {
		return false
	}
	return true
}
