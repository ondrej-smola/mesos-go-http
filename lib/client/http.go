package client

import (
	"context"
	"io"
	"net/http"

	"github.com/gogo/protobuf/proto"
)

const MESOS_STREAM_ID_HEADER = "Mesos-Stream-Id"

type (
	Provider interface {
		New(endpoint string, opts ...DefaultClientOpt) Client
	}

	ProviderFunc func(endpoint string, opts ...DefaultClientOpt) Client

	Response interface {
		Read(m proto.Message) error
		StreamId() string
		io.Closer
	}

	DoFunc func(proto.Message, context.Context, ...RequestOpt) (resp Response, err error)

	Client interface {
		Do(proto.Message, context.Context, ...RequestOpt) (resp Response, err error)
	}
)

func (cf DoFunc) Do(m proto.Message, ctx context.Context, opts ...RequestOpt) (Response, error) {
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

// With header returns an RequestOpt that adds a header value to an HTTP requests's header.
func WithHeader(k, v string) RequestOpt {
	return func(r *http.Request) {
		r.Header.Add(k, v)
	}
}

// WithAuthorization returns an RequestOpt that sets authorization header value
func WithAuthorization(auth string) RequestOpt {
	return func(r *http.Request) {
		r.Header.Set("Authorization", auth)
	}
}

// Close returns a RequestOpt that determines whether to close the underlying connection after sending the request.
func WithClose(b bool) RequestOpt {
	return func(r *http.Request) {
		r.Close = b
	}
}
