package minrevpro

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultAddr = ":8080"
)

// Types

type ReverseProxy struct {
	addr         string
	target       *url.URL
	debug        bool
	secretHeader string
	secret       string
	mw           []Middleware
	log          *log.Logger
	basePath     string

	server *http.Server
	proxy  *httputil.ReverseProxy
}

type OptionFunc func(*ReverseProxy)

type Middleware func(next http.Handler) http.Handler

// NewReverseProxy creates a new proxy instance for given target with the
// default configuration.
func NewReverseProxy(target *url.URL, opts ...OptionFunc) *ReverseProxy {
	p := &ReverseProxy{
		target: target,
		addr:   DefaultAddr,
		log:    log.New(os.Stdout, "", log.LstdFlags),
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Options

// WithAddr sets the addr the proxy serves on.
func WithAddr(addr string) OptionFunc {
	return func(p *ReverseProxy) {
		p.addr = addr
	}
}

// WithLogger sets a custom logger
func WithLogger(log *log.Logger) OptionFunc {
	return func(p *ReverseProxy) {
		p.log = log
	}
}

// WithBasePath sets a custom logger
func WithBasePath(path string) OptionFunc {
	return func(p *ReverseProxy) {
		p.basePath = path
	}
}

// WithSecret sets the secret of the proxy.
// If set, the requests will only be forwareded if the
// given header hase the given value
func WithSecret(header, secret string) OptionFunc {
	return func(p *ReverseProxy) {
		p.secretHeader = header
		p.secret = secret
	}
}

// Debug sets the debug flag of the proxy.
// If set to true, the proxy will log every request.
func Debug(debug bool) OptionFunc {
	return func(p *ReverseProxy) {
		p.debug = debug
	}
}

// WithMiddleware allows to run custom http middleware
func WithMiddleware(mw ...Middleware) OptionFunc {
	return func(p *ReverseProxy) {
		p.mw = append(p.mw, mw...)
	}
}

// Methods
func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.basePath != "" {
		r.URL.Path = strings.Replace(r.URL.Path, p.basePath, "", 1)
		r.RequestURI = r.URL.RequestURI()
	}
	if p.secret != "" {
		secret := r.Header.Get(p.secretHeader)
		if secret == "" {
			if p.debug {
				p.log.Printf("WARNING: got unauthorized request")
			}
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Missing or empty header " + p.secretHeader))
			return
		} else if secret != p.secret {
			if p.debug {
				p.log.Printf("WARNING: got request with invalid secret")
			}
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("Invalid secret " + p.secret))
			return
		}
	}
	r.Host = p.target.Host
	wc := NewInformativeResponseWriter(w)
	start := time.Now()
	if p.debug {
		defer func() {
			p.log.Print(r.Method + " " + r.URL.String() + "(" + p.target.String() + r.URL.Path + "), HTTP " + wc.StatusCodeString() + " " + strconv.FormatUint(wc.BytesWritten(), 10) + " " + time.Since(start).String())
		}()
	}
	p.proxy.ServeHTTP(wc, r)
}

func (p *ReverseProxy) Start() (err error) {
	if p.server != nil {
		return errors.New("already running")
	}
	p.server = &http.Server{
		Addr:    p.addr,
		Handler: p,
	}
	err = p.server.ListenAndServe()
	if err != nil {
		p.server = nil
	}
	return
}

func (p *ReverseProxy) Stop() error {
	if p.server != nil {
		return p.server.Close()
	}
	return errors.New("not running")
}
