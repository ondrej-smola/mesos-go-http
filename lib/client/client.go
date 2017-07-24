package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/lib/codec"
	"github.com/ondrej-smola/mesos-go-http/lib/codec/framing"
	"github.com/ondrej-smola/mesos-go-http/lib/codec/framing/recordio"
	"github.com/ondrej-smola/mesos-go-http/lib/codec/framing/single"
	"github.com/pkg/errors"
)

type (
	HttpDoFunc func(*http.Request) (*http.Response, error)

	DefaultClientOpt func(c *DefaultClient)

	// client for sending requests to Mesos masters and agents
	DefaultClient struct {
		codec       *codec.Codec
		endpoint    string
		requestOpts []RequestOpt

		errorMapper ErrorMapperFunc
		framing     framing.Provider
		do          HttpDoFunc
	}

	response struct {
		head   http.Header
		dec    *codec.Decoder
		body   io.Closer
		cancel context.CancelFunc
	}
)

func (p ProviderFunc) New(endpoint string, opts ...DefaultClientOpt) Client {
	return p(endpoint, opts...)
}

func (r *response) Close() error {
	r.cancel()
	return r.body.Close()
}

func (r *response) StreamId() string {
	return r.head.Get(MESOS_STREAM_ID_HEADER)
}

func (r *response) Read(m proto.Message) error {
	if r.dec == nil {
		return ReadOnEmptyResponse
	}

	return r.dec.Decode(m)
}

func WithErrorMapper(e ErrorMapperFunc) DefaultClientOpt {
	return func(c *DefaultClient) {
		c.errorMapper = e
	}
}

func WithHttpDoFunc(do HttpDoFunc) DefaultClientOpt {
	return func(c *DefaultClient) {
		c.do = do
	}
}

func WithRequestOpts(opts ...RequestOpt) DefaultClientOpt {
	return func(c *DefaultClient) {
		c.requestOpts = opts
	}
}

func WithCodec(codec *codec.Codec) DefaultClientOpt {
	return func(c *DefaultClient) {
		c.codec = codec
	}
}

func WithRecordIOFraming() DefaultClientOpt {
	return func(c *DefaultClient) {
		c.framing = recordio.NewProvider()
	}
}

func WithSingleFraming() DefaultClientOpt {
	return func(c *DefaultClient) {
		c.framing = single.NewProvider()
	}
}

func WithFramingProvider(p framing.Provider) DefaultClientOpt {
	return func(c *DefaultClient) {
		c.framing = p
	}
}

func NewProvider(opts ...DefaultClientOpt) Provider {
	return ProviderFunc(func(endpoint string, addOpts ...DefaultClientOpt) Client {
		res := append(opts, addOpts...)
		return New(endpoint, res...)
	})
}

var defaultDoFunc = NewDoFunc()

func New(endpoint string, opts ...DefaultClientOpt) *DefaultClient {
	if endpoint == "" {
		panic("Endpoint cannot be blank")
	}

	client := &DefaultClient{
		codec:       codec.ProtobufCodec,
		errorMapper: DefaultErrorMapper,
		do:          defaultDoFunc,
		endpoint:    endpoint,
		framing:     single.NewProvider(),
	}

	for _, o := range opts {
		o(client)
	}

	return client
}

func (c *DefaultClient) createRequest(m proto.Message, ctx context.Context, opts ...RequestOpt) (*http.Request, error) {
	cod := c.codec
	buf := &bytes.Buffer{}
	err := cod.NewEncoder(buf).Encode(m)
	if err != nil {
		return nil, errors.Wrap(err, "Message encoding failed")
	}

	req, err := http.NewRequest("POST", c.endpoint, buf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create request entity")
	}

	req.WithContext(ctx)
	helper := HTTPRequestHelper{req}

	return helper.
		withOptions(c.requestOpts, opts).
		withHeader("Content-Type", string(cod.EncoderContentType)).
		withHeader("Accept", string(cod.DecoderContentType)).Unwrap(), nil
}

func (c *DefaultClient) handleResponse(resp *http.Response, err error, cancel context.CancelFunc) (Response, error) {
	if err == nil {
		err = c.errorMapper(resp)
	}

	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	var dec *codec.Decoder
	switch resp.StatusCode {
	case http.StatusOK:
		// If response is not empty (-1 is set when unknown/chunked)
		if resp.ContentLength != 0 {
			ct := resp.Header.Get("Content-Type")
			if ct != c.codec.DecoderContentType {
				resp.Body.Close()
				return nil, errors.Errorf(
					"Unexpected context-type %v - should be %v", ct, c.codec.DecoderContentType,
				)
			}
			dec = c.codec.NewDecoder(c.framing(resp.Body))
		}
	}

	return &response{
		head:   resp.Header,
		dec:    dec,
		body:   resp.Body,
		cancel: cancel,
	}, nil
}

func (c *DefaultClient) Endpoint() string {
	return c.endpoint
}

func (c *DefaultClient) Do(m proto.Message, ctx context.Context, opts ...RequestOpt) (Response, error) {
	var req *http.Request

	reqOpts := append(c.requestOpts, opts...)

	ctx, cancel := context.WithCancel(ctx)
	req, err := c.createRequest(m, ctx, reqOpts...)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "Request build failed")
	}

	res, err := c.do(req)
	return c.handleResponse(res, err, cancel)
}
