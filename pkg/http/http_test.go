package http

import "testing"

func TestValidate_ValidUrl(t *testing.T) {
	h := &HttpCheck{
		Url:    "http://localhost:8080/health",
		Status: 200,
	}
	if err := h.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}
}

func TestValidate_InvalidUrl(t *testing.T) {
	h := &HttpCheck{
		Url:    "not-a-url",
		Status: 200,
	}
	if err := h.Validate(); err == nil {
		t.Error("Validate() should return error for invalid URL")
	}
}

func TestValidate_InvalidStatusCode(t *testing.T) {
	h := &HttpCheck{
		Url:    "http://localhost",
		Status: 0,
	}
	if err := h.Validate(); err == nil {
		t.Error("Validate() should return error for invalid status code")
	}
}

func TestValidate_BothInvalid(t *testing.T) {
	h := &HttpCheck{
		Url:    "",
		Status: 999,
	}
	if err := h.Validate(); err == nil {
		t.Error("Validate() should return error for both invalid URL and status code")
	}
}