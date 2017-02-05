package flow

import "context"

// TODO(os) needs testing
// Purpose of TestFlow is to allow easy testing of Flow related components and as utility for library users to mock
// mock this library
type (
	msg struct {
		m   Message
		ctx context.Context
		err error
	}

	TestFlow struct {
		push chan chan *msg
		pull chan chan *msg
	}

	TestPushReply struct {
		to chan *msg
	}

	TestPullReply struct {
		to chan *msg
	}
)

func NewTestFlow() *TestFlow {
	return &TestFlow{
		push: make(chan chan *msg),
		pull: make(chan chan *msg),
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

// caller must reply to every receive
func (t *TestFlow) AcceptPush() (Message, context.Context, *TestPushReply) {
	req := <-t.push
	msg := <-req
	return msg.m, msg.ctx, &TestPushReply{to: req}
}

// caller must reply to every receive
func (t *TestFlow) AcceptPull() (context.Context, *TestPullReply) {
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
