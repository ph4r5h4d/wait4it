package influxdb

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"wait4it/pkg/model"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

func (i *InfluxDBConnection) BuildContext(cx model.CheckContext) {
	i.Host = cx.Host
	i.Port = cx.Port
	i.Token = cx.InfluxConf.Token
	// Org/Bucket intentionally not copied; not required for health check.
}

func (i *InfluxDBConnection) Validate() error {
	if len(i.Host) == 0 {
		return errors.New("host can't be empty")
	}

	if i.Port < 1 || i.Port > 65535 {
		return errors.New("invalid port range for influxdb")
	}

	return nil
}

func (i *InfluxDBConnection) Check(ctx context.Context) (bool, bool, error) {
	url := i.buildURL()

	// Token can be empty for basic health checks against the public /health endpoint
	client := influxdb2.NewClient(url, i.Token)
	defer client.Close()

	health, err := client.Health(ctx)
	if err != nil || health.Status != domain.HealthCheckStatusPass {
		// Transient error (or non-pass status) during startup or network hiccup -> retry
		return false, false, nil
	}

	return true, true, nil
}

func (i *InfluxDBConnection) buildURL() string {
	// Support scheme in Host (e.g. https://host) for TLS, consistent with
	// elasticsearch checker. Falls back to http if no scheme provided.
	h := i.Host
	if !strings.HasPrefix(strings.ToLower(h), "http://") && !strings.HasPrefix(strings.ToLower(h), "https://") {
		h = "http://" + h
	}
	return h + ":" + strconv.Itoa(i.Port)
}
