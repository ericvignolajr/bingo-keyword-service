package rest

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
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
	cj, _ := cookiejar.New(nil)
	c := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     cj,
	}
	testURL, _ := url.Parse("http://localhost")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.Jar.SetCookies(testURL, tc.cookies)
			f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			h := validateGcsrf(f)
			r := httptest.NewRequest(tc.method, testURL.String(), nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, r)

			assert.Equal(t, http.StatusForbidden, rr.Result().StatusCode)
			c.Jar, _ = cookiejar.New(nil) // reset cookie jar for next test
		})
	}
}
