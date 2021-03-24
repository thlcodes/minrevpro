package minrevpro_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/thlcodes/minrevpro"
)

func TestReverseProxy(t *testing.T) {
	type args struct {
		target string
		opts   []minrevpro.OptionFunc
	}
	type mocks struct {
		code int
		body string
	}
	type wants struct {
		code int
		body string
	}

	tests := []struct {
		name  string
		args  args
		mocks mocks
		wants wants
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := startServer(tt.mocks.code, tt.mocks.body)
			defer server.Close()
			targetURL, _ := url.Parse(tt.args.target)
			rp := minrevpro.NewReverseProxy(targetURL)
			assertNotNil(t, rp, "NewReverseProxy returned nil")
		})
	}
}

func startServer(code int, body string) (server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		_, _ = w.Write([]byte(body))
	}))
	return
}
