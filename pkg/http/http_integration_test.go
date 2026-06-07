package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"wait4it/pkg/model"
)

func newTestCheck(t *testing.T, url string, status int, text string, followRedirect bool) *HttpCheck {
	t.Helper()
	h := &HttpCheck{
		Url:            url,
		Status:         status,
		Text:           text,
		FollowRedirect: followRedirect,
	}
	if err := h.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
	return h
}

func TestCheck_StatusMatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := newTestCheck(t, ts.URL, http.StatusOK, "", false)
	h.client = ts.Client()

	ok, endOnError, err := h.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false")
	}
	if err != nil {
		t.Errorf("Check() returned unexpected error: %v", err)
	}
}

func TestCheck_StatusMismatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	h := newTestCheck(t, ts.URL, http.StatusOK, "", false)
	h.client = ts.Client()

	ok, endOnError, err := h.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for status mismatch")
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false for status mismatch")
	}
	if err == nil {
		t.Error("Check() should return error for status mismatch")
	}
}

func TestCheck_BodyTextMatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	h := newTestCheck(t, ts.URL, http.StatusOK, "hello", false)
	h.client = ts.Client()

	ok, _, err := h.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
	if err != nil {
		t.Errorf("Check() returned unexpected error: %v", err)
	}
}

func TestCheck_BodyTextMismatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	h := newTestCheck(t, ts.URL, http.StatusOK, "goodbye", false)
	h.client = ts.Client()

	ok, _, err := h.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for text mismatch")
	}
	if err == nil {
		t.Error("Check() should return error for text mismatch")
	}
}

func TestCheck_FollowRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := newTestCheck(t, ts.URL+"/redirect", http.StatusOK, "", true)
	h.client = ts.Client()

	ok, _, err := h.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true (should follow redirect); err = %v", err)
	}
}

func TestCheck_NoFollowRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// When not following redirects, the client receives the 302 status code
	h := newTestCheck(t, ts.URL+"/redirect", http.StatusFound, "", false)
	// Use a client that does not follow redirects (matching the production behavior)
	h.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: ts.Client().Transport,
	}

	ok, _, err := h.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true (should not follow redirect); err = %v", err)
	}
}

func TestCheck_ContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		select {}
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	h := newTestCheck(t, ts.URL, http.StatusOK, "", false)
	h.client = ts.Client()

	ok, _, err := h.Check(ctx)
	if ok {
		t.Error("Check() ok = true, want false for cancelled context")
	}
	if err == nil {
		t.Error("Check() should return error for cancelled context")
	}
}

func TestCheck_ConnectionRefused(t *testing.T) {
	// Use a URL that will refuse connections — no server listening here
	h := &HttpCheck{
		Url:    "http://127.0.0.1:1",
		Status: http.StatusOK,
	}
	if err := h.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, endOnError, err := h.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for connection refused")
	}
	if !endOnError {
		t.Error("Check() endOnError = false, want true for connection refused")
	}
	if err == nil {
		t.Error("Check() should return error for connection refused")
	}
}

func TestNewChecker_WithContext(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cx := &model.CheckContext{
		Host: ts.URL,
		HttpConf: model.HttpSpecificConf{
			StatusCode:     http.StatusOK,
			Text:           "",
			FollowRedirect: false,
		},
	}

	checker, err := NewChecker(cx)
	if err != nil {
		t.Fatalf("NewChecker() returned unexpected error: %v", err)
	}

	// Inject the test client after construction
	h := checker.(*HttpCheck)
	h.client = ts.Client()

	ok, _, err := checker.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}