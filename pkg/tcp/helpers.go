package tcp

func (c *checker) isPortInValidRange() bool {
	if c.port < minPort || c.port > maxPort {
		return false
	}
	return true
}
