package flow

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// TODO(os) experimental
// Purpose of TestFlow is to allow easy testing of Flow related components and as utility for library users to
// mock this library
type (
	msg struct {
		m   Message
		ctx context.Context
		err error
	}

	TestFlow struct {
		push     chan chan *msg
		pull     chan chan *msg
		closeIn  chan bool
		closeOut chan error
	}

	TestFlowBlueprint struct {
		NextFlowIn  chan bool
		NextFlowOut chan Flow
	}

	TestPushReply struct {
		to chan *msg
	}

	TestPullReply struct {
		to chan *msg
	}

	CloseReply struct {
		to chan error
	}
)

func NewTestFlowBlueprint() *TestFlowBlueprint {
	return &TestFlowBlueprint{
		NextFlowIn:  make(chan bool),
		NextFlowOut: make(chan Flow),
	}
}

func (t *TestFlowBlueprint) Mat(opts ...MatOpt) Flow {
	t.NextFlowIn <- true
	return <-t.NextFlowOut
}

func NewTestFlow() *TestFlow {
	return &TestFlow{
		push:     make(chan chan *msg),
		pull:     make(chan chan *msg),
		closeIn:  make(chan bool),
		closeOut: make(chan error),
	}
}

func (r *TestPushReply) OK() {
	r.to <- &msg{}
	close(r.to)
}

func (r *TestPushReply) Error(err error) {
	r.to <- &msg{err: err}
	close(r.to)
}

func (r *TestPullReply) Message(m Message) {
	r.to <- &msg{m: m}
	close(r.to)
}

func (r *TestPullReply) Error(err error) {
	r.to <- &msg{err: err}
	close(r.to)
}

func (r *CloseReply) OK() {
	close(r.to)
}

func (r *CloseReply) Error(err error) {
	r.to <- err
	close(r.to)
}

func (t *TestFlow) ExpectNoPush(wait time.Duration) error {
	select {
	case p := <-t.push:
		return errors.Errorf("Unexpected push: %v", <-p)
	case <-time.After(wait):
		return nil
	}
}

func (t *TestFlow) ExpectNoPull(wait time.Duration) error {
	select {
	case <-t.pull:
		return errors.New("Unexpected pull")
	case <-time.After(wait):
		return nil
	}
}

// caller must reply to every receive
func (t *TestFlow) ExpectPush() (Message, context.Context, *TestPushReply) {
	req := <-t.push
	msg := <-req
	return msg.m, msg.ctx, &TestPushReply{to: req}
}

// caller must reply to every receive
func (t *TestFlow) ExpectPull() (context.Context, *TestPullReply) {
	req := <-t.pull
	msg := <-req
	return msg.ctx, &TestPullReply{to: req}
}

func (t *TestFlow) Pull(ctx context.Context) (Message, error) {
	req := make(chan *msg)
	t.pull <- req
	req <- &msg{ctx: ctx}
	res := <-req
	return res.m, res.err

}
func (t *TestFlow) Push(m Message, ctx context.Context) error {
	req := make(chan *msg)
	t.push <- req
	req <- &msg{m: m, ctx: ctx}
	res := <-req
	return res.err
}

func (t *TestFlow) Close() error {
	t.closeIn <- true
	return <-t.closeOut
}

// can be called only once
func (t *TestFlow) AcceptClose() *CloseReply {
	<-t.closeIn
	return &CloseReply{to: t.closeOut}
}
