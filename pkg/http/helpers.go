package http

import "net/url"

func (h *HttpCheck) validateUrl() bool {
	_, err := url.ParseRequestURI(h.Url)
	if err != nil {
		return false
	}

	u, err := url.Parse(h.Url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func (h *HttpCheck) validateStatusCode() bool {
	// check against common status code
	if h.Status < 100 || h.Status > 599 {
		return false
	}
	return true
}
