package http

import "net/url"

func (c *checker) validateUrl() bool {
	_, err := url.ParseRequestURI(c.url)
	if err != nil {
		return false
	}

	u, err := url.Parse(c.url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func (c *checker) validateStatusCode() bool {
	// check against common status code
	if c.status < 100 || c.status > 599 {
		return false
	}
	return true
}
