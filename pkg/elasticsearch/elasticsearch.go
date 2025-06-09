package elasticsearch

import (
	"context"
	"errors"
	"strconv"
	"wait4it/pkg/model"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearchChecker struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (esc *ElasticSearchChecker) BuildContext(cx model.CheckContext) {
	esc.Host = cx.Host
	esc.Port = cx.Port
	esc.Username = cx.Username
	esc.Password = cx.Password()
}

func (esc *ElasticSearchChecker) Validate() error {
	if len(esc.Host) == 0 {
		return errors.New("Host can't be empty")
	}

	if esc.Port < 1 || esc.Port > 65535 {
		return errors.New("Invalid port range for ElasticSearch")
	}

	return nil
}

func (esc *ElasticSearchChecker) Check(ctx context.Context) (bool, bool, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			esc.BuildConnectionString(),
		},
		Username: esc.Username,
		Password: esc.Password,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return false, true, err
	}

	if _, err := es.Ping(es.Ping.WithContext(ctx)); err != nil {
		return false, false, err
	}

	return true, true, nil
}

func (esc *ElasticSearchChecker) BuildConnectionString() string {
	return esc.Host + ":" + strconv.Itoa(esc.Port)
}

// NewChecker creates a new checker
func NewChecker(c *model.CheckContext) (model.CheckInterface, error) {
	check := &ElasticSearchChecker{}
	check.BuildContext(*c)
	if err := check.Validate(); err != nil {
		return nil, err
	}

	return check, nil
}
