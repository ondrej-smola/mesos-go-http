package mesos

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"io"
)

type TestEmptyResponse struct {
	Sid string
}

func (n *TestEmptyResponse) Read(m proto.Message) error {
	return io.EOF
}

func (n *TestEmptyResponse) Close() error {
	return nil
}

func (t *TestEmptyResponse) StreamId() string {
	return t.Sid
}

type TestMessageOrError struct {
	Msg proto.Message
	Err error
}

type TestChanResponse struct {
	Sid      string
	CloseIn  chan bool
	CloseOut chan error

	ReadIn  chan bool
	ReadOut chan *TestMessageOrError
}

func NewTestChanResponse(optionalStreamId ...string) *TestChanResponse {
	r := &TestChanResponse{
		CloseIn:  make(chan bool),
		CloseOut: make(chan error),

		ReadIn:  make(chan bool),
		ReadOut: make(chan *TestMessageOrError),
	}

	if len(optionalStreamId) > 0 {
		r.Sid = optionalStreamId[0]
	}

	return r
}

func (t *TestChanResponse) Read(m proto.Message) error {
	t.ReadIn <- true
	resp := <-t.ReadOut
	if resp.Err != nil {
		return resp.Err
	} else {
		// copy received message to m
		body, err := proto.Marshal(resp.Msg)
		if err != nil {
			return err
		}
		return proto.Unmarshal(body, m)
	}
}

func (t *TestChanResponse) Close() error {
	t.CloseIn <- true
	return <-t.CloseOut
}

func (t *TestChanResponse) StreamId() string {
	return t.Sid
}

type TestClientMessageWithContext struct {
	Msg  proto.Message
	Ctx  context.Context
	Opts []RequestOpt
}

type TestClientResponseOrError struct {
	Resp Response
	Err  error
}

type TestChanClient struct {
	ReqIn  chan *TestClientMessageWithContext
	ReqOut chan *TestClientResponseOrError
}

func NewTestChanClient() *TestChanClient {
	return &TestChanClient{
		ReqIn:  make(chan *TestClientMessageWithContext),
		ReqOut: make(chan *TestClientResponseOrError),
	}
}

func (t *TestChanClient) Do(m proto.Message, ctx context.Context, opts ...RequestOpt) (Response, error) {
	t.ReqIn <- &TestClientMessageWithContext{Msg: m, Ctx: ctx, Opts: opts}
	resp := <-t.ReqOut
	return resp.Resp, resp.Err
}
