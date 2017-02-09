package scheduler

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/pkg/errors"
)

var ErrBufferFull = errors.New("Scheduler: buffer is full")

const (
	DEFAULT_EVENT_BUFFER_SIZE = 64
)

type (
	Opt func(c *Client)

	event struct {
		event flow.Message
		ctx   context.Context
		err   error
		reqId uint64
	}

	// Handles thread-safe communication with Mesos leader
	// Supports multiple concurrent read and write requests
	// First request must be subscribe call
	Client struct {
		client mesos.Client

		// client is subscribed when not empty
		streamId string

		// buffer for events from Mesos and for loopback message (not mesos Call)
		buffer chan flow.Message
		// write request queue
		write chan chan *event
		// read request queue
		read chan chan *event

		ctx    context.Context
		cancel context.CancelFunc

		log log.Logger
	}
)

func V1SchedulerAPIEndpointFunc(hostPort string) string {
	return fmt.Sprintf("http://%v/api/v1/scheduler", hostPort)
}

func WithLogger(l log.Logger) Opt {
	return func(c *Client) {
		c.log = l
	}
}

var _ = flow.Flow(&Client{})

func Blueprint(client mesos.Client, opts ...Opt) flow.SinkBlueprint {
	return flow.SinkBlueprintFunc(func(matOpts ...flow.MatOpt) flow.Sink {
		cfg := flow.MatOpts(matOpts).Config()

		if cfg.Log != nil {
			opts = append(opts, WithLogger(log.NewContext(cfg.Log).With("sink", "scheduler")))
		}
		return New(client, opts...)
	})
}

func New(client mesos.Client, opts ...Opt) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		ctx:    ctx,
		cancel: cancel,
		client: client,
		write:  make(chan chan *event),
		read:   make(chan chan *event),
		buffer: make(chan flow.Message, DEFAULT_EVENT_BUFFER_SIZE),
		log:    log.NewNopLogger(),
	}

	for _, o := range opts {
		o(c)
	}

	go c.schedulerLoop()

	return c
}

func (c *Client) Push(ev flow.Message, ctx context.Context) error {
	req := make(chan *event)

	select {
	case c.write <- req:
		req <- &event{event: ev, ctx: ctx}
	case <-ctx.Done():
		return ctx.Err()
	case <-c.ctx.Done():
		return c.ctx.Err()
	}

	res := <-req

	if res != nil && res.err != nil {
		return res.err
	} else {
		return nil
	}
}

func (c *Client) Pull(ctx context.Context) (flow.Message, error) {
	req := make(chan *event)

	select {
	case c.read <- req:
		req <- &event{ctx: ctx}
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	}

	res := <-req
	// sanity check
	if res == nil || (res.event == nil && res.err == nil) || (res.event != nil && res.err != nil) {
		panic(fmt.Sprintf("Invalid pull response %v", res))
	}

	return res.event, res.err
}

func (c *Client) Close() error {
	c.cancel()
	close(c.read)
	close(c.write)
	c.log.Log("event", "closed")
	return nil
}

// implements flow.Sink interface
func (c *Client) IsSink() {}

// main client loop
func (c *Client) schedulerLoop() {
	l := log.NewContext(c.log).With("internal", "scheduler_loop")

	for {
		select {
		case <-c.ctx.Done():
			l.Log("event", "cancelled")
			return
		case r := <-c.read:
			// read requests are dispatched asynchronously
			go c.doRead(r)
		case w := <-c.write:
			// write requests are dispatched asynchronously except first call (subscribe)
			if c.streamId == "" {
				resp, err := c.subscribe(<-w)
				if err == nil {
					l.Log("event", "subscribed", "stream-id", c.streamId)
					go c.readerLoop(resp, c.ctx)
				}
				w <- &event{err: err}
			} else {
				go c.doWrite(w, c.streamId)
			}
		}
	}
}

func (c *Client) doWrite(w chan *event, streamId string) {
	defer close(w)
	ev := <-w

	switch r := ev.event.(type) {
	case *Call:
		resp, err := c.client.Do(r, ev.ctx, mesos.WithMesosStreamId(streamId))
		if err == nil {
			resp.Close()
		} else {
			w <- &event{err: err}
		}
	default:
		select {
		case c.buffer <- ev.event:
		default:
			w <- &event{err: ErrBufferFull}
		}
	}
}

func (c *Client) doRead(read chan *event) {
	defer close(read)
	req := <-read

	select {
	case ev := <-c.buffer:
		read <- &event{event: ev}
	case <-c.ctx.Done():
		read <- &event{err: c.ctx.Err()}
	case <-req.ctx.Done():
		read <- &event{err: req.ctx.Err()}
	}
}

// returns mesos-stream-id or error
func (c *Client) subscribe(ev *event) (mesos.Response, error) {
	isSubscribe, subscribe := IsSubscribeMessage(ev.event)

	if !isSubscribe {
		return nil, errors.Errorf("First message must be subscribe - but is %T", ev.event)
	}

	// set Close(true) because subscribe call cannot reuse connections
	if res, err := c.client.Do(subscribe, c.ctx, mesos.WithClose(true)); err != nil {
		return nil, err
	} else {
		sid := res.StreamId()
		if sid == "" {
			return nil, errors.Errorf("Mesos is expected to set %v but it is empty", mesos.MESOS_STREAM_ID_HEADER)
		}
		c.streamId = sid
		return res, nil
	}
}

func (c *Client) readerLoop(resp mesos.Response, ctx context.Context) {
	defer c.log.Log("event", "reader_loop", "state", "cancelled")
	defer resp.Close()

	for {
		msg := &Event{}
		err := resp.Read(msg)
		if err != nil {
			c.cancel()
			return
		} else {
			select {
			case c.buffer <- flow.Message(msg):
			case <-ctx.Done():
				return
			}
		}
	}
}
