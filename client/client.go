package client

import (
	"bytes"
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"github.com/ondrej-smola/mesos-go-http/codec/framing"
	"github.com/ondrej-smola/mesos-go-http/codec/framing/recordio"
	"github.com/ondrej-smola/mesos-go-http/codec/framing/single"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"sync"
)

type (
	Provider func(endpoint string, opts ...Opt) mesos.Client

	RequestBuilder func(proto.Message, context.Context, ...mesos.RequestOpt) (*http.Request, error)

	DoFunc func(*http.Request) (*http.Response, error)

	// Arguments are response and error from DoFunc and cancel function for canceling request context
	ResponseHandler func(*http.Response, error, context.CancelFunc) (mesos.Response, error)

	Opt func(c *Client)

	// client for sending requests to Mesos masters and agents
	Client struct {
		codec       *codec.Codec
		endpoint    string
		requestOpts []mesos.RequestOpt

		errorMapper ErrorMapperFunc
		response    ResponseHandler
		request     RequestBuilder
		framing     framing.Provider
		do          DoFunc
	}

	stream struct {
		readCancel context.CancelFunc
		respCloser io.Closer
		decoder    *codec.Decoder
	}

	readerState struct {
		readCancel context.CancelFunc
		respCloser io.Closer
		decoder    *codec.Decoder
		sync.Mutex
	}

	response struct {
		head   http.Header
		dec    *codec.Decoder
		body   io.Closer
		cancel context.CancelFunc
	}
)

func (r *response) Close() error {
	r.cancel()
	return r.body.Close()
}

// implements scheduler.ResponseWithMesosStreamId
func (r *response) StreamId() string {
	return r.head.Get(mesos.MESOS_STREAM_ID_HEADER)
}

func (r *response) Read(m proto.Message) error {
	if r.dec == nil {
		return EmptyResponseErr
	}

	return r.dec.Decode(m)
}

func (p Provider) New(endpoint string, opts ...Opt) mesos.Client {
	return p(endpoint, opts...)
}

// HandleResponse returns a functional config option to set the HTTP response handler of the client.
func WithResponseHandler(f ResponseHandler) Opt {
	return func(c *Client) {
		c.response = f
	}
}

func WithErrorMapper(e ErrorMapperFunc) Opt {
	return func(c *Client) {
		c.errorMapper = e
	}
}

func WithDoFunc(do DoFunc) Opt {
	return func(c *Client) {
		c.do = do
	}
}

func WithRequestBuilder(r RequestBuilder) Opt {
	return func(c *Client) {
		c.request = r
	}
}

func WithRequestOpts(opts ...mesos.RequestOpt) Opt {
	return func(c *Client) {
		c.requestOpts = opts
	}
}

func WithWrappedDoer(f func(DoFunc) DoFunc) Opt {
	return func(c *Client) {
		f(c.do)
	}
}

func WithCodec(codec *codec.Codec) Opt {
	return func(c *Client) {
		c.codec = codec
	}
}

func WithRecordIOFraming() Opt {
	return func(c *Client) {
		c.framing = recordio.NewProvider()
	}
}

func WithSingleFraming() Opt {
	return func(c *Client) {
		c.framing = single.NewProvider()
	}
}

func WithFramingProvider(p framing.Provider) Opt {
	return func(c *Client) {
		c.framing = p
	}
}

func NewProvider(opts ...Opt) Provider {
	return func(endpoint string, addOpts ...Opt) mesos.Client {
		res := append(opts, addOpts...)
		return New(endpoint, res...)
	}
}

func New(endpoint string, opts ...Opt) *Client {
	if endpoint == "" {
		panic("Endpoint cannot be blank")
	}

	client := &Client{
		codec:       codec.ProtobufCodec,
		errorMapper: DefaultErrorMapper,
		do:          NewDoFunc(),
		endpoint:    endpoint,
		framing:     single.NewProvider(),
	}

	client.request = client.buildRequest
	client.response = client.handleResponse

	for _, o := range opts {
		o(client)
	}

	return client
}

func (c *Client) buildRequest(m proto.Message, ctx context.Context, opts ...mesos.RequestOpt) (*http.Request, error) {
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

func (c *Client) handleResponse(resp *http.Response, err error, cancel context.CancelFunc) (mesos.Response, error) {
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
		ct := resp.Header.Get("Content-Type")
		if ct != c.codec.DecoderContentType {
			resp.Body.Close()
			return nil, errors.Errorf("Unexpected context-type  %v - should be %v", ct, c.codec.DecoderContentType)
		}
		dec = c.codec.NewDecoder(c.framing(resp.Body))
	}

	return &response{
		head:   resp.Header,
		dec:    dec,
		body:   resp.Body,
		cancel: cancel,
	}, nil
}

func (c *Client) Endpoint() string {
	return c.endpoint
}

func (c *Client) Do(m proto.Message, ctx context.Context, opts ...mesos.RequestOpt) (mesos.Response, error) {
	var req *http.Request

	reqOpts := append(c.requestOpts, opts...)

	ctx, cancel := context.WithCancel(ctx)
	req, err := c.request(m, ctx, reqOpts...)
	if err != nil {
		cancel()
		return nil, errors.Wrap(err, "Request build failed")
	}

	res, err := c.do(req)
	return c.response(res, err, cancel)
}
