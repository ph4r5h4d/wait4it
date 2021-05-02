package checkers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"wait4it"
)

type httpChecker struct {
	url            string
	statusCode     int
	body           string
	followRedirect bool
}

// NewHTTPChecker returns the HTTP checker with support of comparing status code, checking substring
// in the response body, following redirects, and so on. an error also will be returned if any invalid
// arguments are supplied.
func NewHTTPChecker(rawurl string, statusCode int, body string, followRedirect bool) (wait4it.Checkable, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, errors.New("empty http url host")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("http url scheme should be http/https")
	}

	if statusCode < 100 || statusCode > 599 {
		return nil, errors.New("http status code should be between 100 and 599")
	}

	return &httpChecker{
		url:            rawurl,
		statusCode:     statusCode,
		body:           body,
		followRedirect: followRedirect,
	}, nil
}

func (c *httpChecker) Check(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	if c.followRedirect {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return temporaryError(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return temporaryError(err.Error())
	}

	if resp.StatusCode != c.statusCode {
		return fmt.Errorf("unexpetced %d status code", resp.StatusCode)
	}

	if len(c.body) > 0 {
		if !strings.Contains(string(body), c.body) {
			return errors.New("can't find substring in the response body")
		}
	}

	return nil
}
