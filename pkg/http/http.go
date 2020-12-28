package http

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"wait4it/pkg/model"
)

func (c *checker) buildContext(cx model.CheckContext) {
	c.url = cx.Host
	c.status = cx.HttpConf.StatusCode

	if len(cx.HttpConf.Text) > 0 {
		c.text = cx.HttpConf.Text
	}
}

func (c *checker) validate() error {
	if !c.validateUrl() {
		return errors.New("invalid URL provided")
	}

	if !c.validateStatusCode() {
		return errors.New("invalid status code provided")
	}

	return nil
}

func (c *checker) Check(ctx context.Context) (bool, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return false, false, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, true, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, true, err
	}

	if resp.StatusCode != c.status {
		return false, false, errors.New("invalid status code")
	}

	if len(c.text) > 0 {
		if !strings.Contains(string(body), c.text) {
			return false, false, errors.New("can't find substring in response")
		}
	}

	return true, false, nil
}
