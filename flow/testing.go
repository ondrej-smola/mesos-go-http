package flow

import (
	"context"
	"fmt"
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

		maxWait time.Duration
	}

	TestFlowBlueprint struct {
		NextFlowIn  chan bool
		NextFlowOut chan Flow
	}

	TestPushContext struct {
		Msg Message
		Ctx context.Context

		to chan *msg
	}

	TestPullContext struct {
		Ctx context.Context
		to  chan *msg
	}

	CloseReply struct {
		to chan error
	}

	TestFlowOpt func(t *TestFlow)
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

func WithDefaultMaxWait(w time.Duration) TestFlowOpt {
	return func(t *TestFlow) {
		t.maxWait = w
	}
}

func NewTestFlow() *TestFlow {
	return &TestFlow{
		push:     make(chan chan *msg),
		pull:     make(chan chan *msg),
		closeIn:  make(chan bool),
		closeOut: make(chan error),
		maxWait:  time.Second,
	}
}

func (r *TestPushContext) OK() {
	r.to <- &msg{}
	close(r.to)
}

func (r *TestPushContext) Error(err error) {
	r.to <- &msg{err: err}
	close(r.to)
}

func (r *TestPullContext) Message(m Message) {
	r.to <- &msg{m: m}
	close(r.to)
}

func (r *TestPullContext) Error(err error) {
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

func (t *TestFlow) ExpectNoPush(wait ...time.Duration) {
	w := t.maxWait
	if len(wait) == 1 {
		w = wait[0]
	}

	select {
	case p := <-t.push:
		panic(fmt.Sprintf("Unexpected push: %v", <-p))
	case <-time.After(w):
	}
}

func (t *TestFlow) ExpectNoPull(wait ...time.Duration) {
	w := t.maxWait
	if len(wait) == 1 {
		w = wait[0]
	}

	select {
	case <-t.pull:
		panic("Unexpected pull")
	case <-time.After(w):
	}
}

// caller must reply to every receive
func (t *TestFlow) ExpectPush(wait ...time.Duration) *TestPushContext {
	w := t.maxWait
	if len(wait) == 1 {
		w = wait[0]
	}

	select {
	case req := <-t.push:
		msg := <-req
		return &TestPushContext{Msg: msg.m, Ctx: msg.ctx, to: req}
	case <-time.After(w):
		panic("Push timeout")
	}
}

// caller must reply to every receive
func (t *TestFlow) ExpectPull(wait ...time.Duration) *TestPullContext {
	w := t.maxWait
	if len(wait) == 1 {
		w = wait[0]
	}

	select {
	case req := <-t.pull:
		msg := <-req
		return &TestPullContext{Ctx: msg.ctx, to: req}
	case <-time.After(w):
		panic("Pull timeout")
	}
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
