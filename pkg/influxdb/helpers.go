package influxdb

// BuildConnectionURL returns the base URL (with scheme) for the InfluxDB instance.
// Kept for consistency with other checkers and potential future use (e.g. https support).
func (i InfluxDBConnection) BuildConnectionURL() string {
	return i.buildURL()
}
