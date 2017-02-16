package mesos

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"io"
)

type TestEmptyResponse struct{}

func (n *TestEmptyResponse) Read(m proto.Message) error {
	return io.EOF
}

func (n *TestEmptyResponse) Close() error {
	return nil
}

type TestMessageOrError struct {
	Msg proto.Message
	Err error
}

type TestChanResponse struct {
	CloseIn  chan bool
	CloseOut chan error

	ReadIn  chan bool
	ReadOut chan *TestMessageOrError
}

func NewTestChanResponse() *TestChanResponse {
	return &TestChanResponse{
		CloseIn:  make(chan bool),
		CloseOut: make(chan error),

		ReadIn:  make(chan bool),
		ReadOut: make(chan *TestMessageOrError),
	}
}

func (t *TestChanResponse) Read(m proto.Message) error {
	t.ReadIn <- true
	resp := <-t.ReadOut
	if resp.Err != nil {
		return nil
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

func NoopClientFunc(proto.Message, context.Context, ...RequestOpt) (Response, error) {
	return &TestEmptyResponse{}, nil
}
