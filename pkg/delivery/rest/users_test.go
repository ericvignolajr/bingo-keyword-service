package rest

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateGcsrf(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		cookies  []*http.Cookie
		expected int
	}{
		{
			name:   "no cookie",
			method: http.MethodPost,
			cookies: []*http.Cookie{
				{
					Name:  "g_csrf_token",
					Value: "123",
				},
			},
			expected: http.StatusForbidden,
		},
	}
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := &url.URL{
		Host: os.Getenv("SERVER_URL"),
		Path: "/",
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.Jar.SetCookies(url, tc.cookies)
			f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			h := validateGcsrf(f)
			r := httptest.NewRequest(tc.method, url.String(), nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, r)

			assert.Equal(t, http.StatusForbidden, rr.Result().StatusCode)
			c.Jar, _ = cookiejar.New(nil) // reset cookie jar for next test
		})
	}
}
