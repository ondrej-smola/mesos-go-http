package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type DoFuncConfig struct {
	client    *http.Client
	dialer    *net.Dialer
	transport *http.Transport
}

type DoFuncOpt func(*DoFuncConfig)

// Returns a DoFunc that executes HTTP round-trips.
// The default implementation provides reasonable defaults for timeouts:
// keep-alive, connection, request/response read/write, and TLS handshake.
// Callers can customize configuration by specifying one or more DoFunc's.
func NewDoFunc(opt ...DoFuncOpt) HttpDoFunc {
	var (
		dialer = &net.Dialer{
			LocalAddr: &net.TCPAddr{IP: net.IPv4zero},
			KeepAlive: 30 * time.Second,
			Timeout:   5 * time.Second,
		}
		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial:  dialer.Dial,
			ResponseHeaderTimeout: 5 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
			TLSHandshakeTimeout:   5 * time.Second,
		}
		config = &DoFuncConfig{
			dialer:    dialer,
			transport: transport,
			client:    &http.Client{Transport: transport},
		}
	)
	for _, o := range opt {
		if o != nil {
			o(config)
		}
	}
	return config.client.Do
}

// Timeout returns an ConfigOpt that sets a Config's timeout and keep-alive timeout.
func FuncTimeout(d time.Duration) DoFuncOpt {
	return func(c *DoFuncConfig) {
		c.transport.ResponseHeaderTimeout = d
		c.transport.TLSHandshakeTimeout = d
		c.dialer.Timeout = d
	}
}

// RoundTripper returns a ConfigOpt that sets a Config's round-tripper.
func FuncRoundTripper(rt http.RoundTripper) DoFuncOpt {
	return func(c *DoFuncConfig) {
		c.client.Transport = rt
	}
}

// TLSConfig returns a ConfigOpt that sets a Config's TLS configuration.
func TLSConfig(tc *tls.Config) DoFuncOpt {
	return func(c *DoFuncConfig) {
		c.transport.TLSClientConfig = tc
	}
}

// Transport returns a ConfigOpt that allows tweaks of the default Config's http.Transport
func Transport(modifyTransport func(*http.Transport)) DoFuncOpt {
	return func(c *DoFuncConfig) {
		if modifyTransport != nil {
			modifyTransport(c.transport)
		}
	}
}

// WrapRoundTripper allows a caller to customize a configuration's HTTP exchanger. Useful
// for authentication protocols that operate over stock HTTP.
func WrapRoundTripper(f func(http.RoundTripper) http.RoundTripper) DoFuncOpt {
	return func(c *DoFuncConfig) {
		if f != nil {
			if rt := f(c.client.Transport); rt != nil {
				c.client.Transport = rt
			}
		}
	}
}

// HTTPRequestHelper wraps an http.Request and provides utility funcs to simplify code elsewhere
type HTTPRequestHelper struct {
	*http.Request
}

func (r *HTTPRequestHelper) withOptions(optsets ...RequestOpts) *HTTPRequestHelper {
	for _, opts := range optsets {
		opts.Apply(r.Request)
	}
	return r
}

func (r *HTTPRequestHelper) withHeaders(hh http.Header) *HTTPRequestHelper {
	for k, v := range hh {
		r.Header[k] = v
	}
	return r
}

func (r *HTTPRequestHelper) withHeader(key, value string) *HTTPRequestHelper {
	r.Header.Set(key, value)
	return r
}

func (r *HTTPRequestHelper) Unwrap() *http.Request {
	return r.Request
}
