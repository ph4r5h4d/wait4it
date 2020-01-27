package tcpcheck

func (i IP) isPortInValidRange() bool {
	if i.Port < minPort || i.Port > maxPort {
		return false
	}
	return true
}
