package mesos

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"io"
	"net/http"
)

const MESOS_STREAM_ID_HEADER = "Mesos-Stream-Id"

type (

	// transform host:port to full URL
	EndpointFunc func(string) string

	Response interface {
		Read(m codec.Message) error
		// return assigned Mesos-Stream-Id if set
		StreamId() string
		io.Closer
	}

	DoFunc func(codec.Message, context.Context, ...RequestOpt) (resp Response, err error)

	Client interface {
		Do(codec.Message, context.Context, ...RequestOpt) (resp Response, err error)
	}
)

func (cf DoFunc) Do(m codec.Message, ctx context.Context, opts ...RequestOpt) (Response, error) {
	return cf(m, ctx, opts...)
}

// RequestOpt defines a functional option for an http.Request.
type RequestOpt func(*http.Request)

// RequestOpts is a convenience type
type RequestOpts []RequestOpt

func (opts RequestOpts) Apply(req *http.Request) {
	// apply per-request options
	for _, o := range opts {
		if o != nil {
			o(req)
		}
	}
}

// Set mesos stream id request header
func WithMesosStreamId(id string) RequestOpt {
	return WithHeader(MESOS_STREAM_ID_HEADER, id)
}

// Header returns an RequestOpt that adds a header value to an HTTP requests's header.
func WithHeader(k, v string) RequestOpt {
	return func(r *http.Request) {
		r.Header.Add(k, v)
	}
}

// Close returns a RequestOpt that determines whether to close the underlying connection after sending the request.
func WithClose(b bool) RequestOpt {
	return func(r *http.Request) {
		r.Close = b
	}
}
