package checkers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"wait4it"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPChecker(t *testing.T) {
	tt := []struct {
		name           string
		rawurl         string
		statsu         int
		text           string
		followRedirect bool

		checker         wait4it.Checkable
		isErrorExpected bool
	}{
		{
			name:   "no error",
			rawurl: "https://example.com/",
			statsu: 200,
			checker: &httpChecker{
				url:        "https://example.com/",
				statusCode: 200,
			},
			isErrorExpected: false,
		},
		{
			name:            "with invalid url",
			rawurl:          "%invalid^URL!",
			isErrorExpected: true,
		},
		{
			name:            "with invalid host",
			rawurl:          "https://",
			isErrorExpected: true,
		},
		{
			name:            "with invalid url scheme",
			rawurl:          "htt://example.com/",
			isErrorExpected: true,
		},
		{
			name:            "with invalid status code",
			rawurl:          "https://example.com/",
			statsu:          50,
			isErrorExpected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			checker, err := NewHTTPChecker(tc.rawurl, tc.statsu, tc.text, tc.followRedirect)
			if tc.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.checker, checker)
		})
	}
}

func TestHTTPCheckerCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.URL.String(), "/some/path")
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	checker, err := NewHTTPChecker(server.URL+"/some/path", http.StatusOK, "OK", false)
	require.NoError(t, err)

	checker.Check(context.Background())
}
