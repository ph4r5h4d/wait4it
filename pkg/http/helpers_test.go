package http

import "testing"

func TestValidateUrl_ValidUrls(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"http URL", "http://localhost:8080/health"},
		{"https URL", "https://example.com/api/status"},
		{"URL with path", "http://localhost/health/check"},
		{"URL with port", "http://localhost:9200"},
		{"URL with query", "http://localhost:8080/check?verbose=true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpCheck{Url: tt.url}
			if !h.validateUrl() {
				t.Errorf("validateUrl(%q) = false, want true", tt.url)
			}
		})
	}
}

func TestValidateUrl_InvalidUrls(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty string", ""},
		{"no scheme", "localhost:8080"},
		{"no host", "http://"},
		{"just path", "/health"},
		{"random string", "not-a-url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpCheck{Url: tt.url}
			if h.validateUrl() {
				t.Errorf("validateUrl(%q) = true, want false", tt.url)
			}
		})
	}
}

func TestValidateStatusCode_ValidCodes(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"100 Continue", 100},
		{"200 OK", 200},
		{"301 Moved", 301},
		{"404 Not Found", 404},
		{"500 Internal Server Error", 500},
		{"599 max client error", 599},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpCheck{Status: tt.status}
			if !h.validateStatusCode() {
				t.Errorf("validateStatusCode(%d) = false, want true", tt.status)
			}
		})
	}
}

func TestValidateStatusCode_InvalidCodes(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"zero", 0},
		{"negative", -1},
		{"99 too low", 99},
		{"600 too high", 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HttpCheck{Status: tt.status}
			if h.validateStatusCode() {
				t.Errorf("validateStatusCode(%d) = true, want false", tt.status)
			}
		})
	}
}
