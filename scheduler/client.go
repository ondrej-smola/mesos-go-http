package scheduler

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/log"
	"github.com/pkg/errors"
)

var ErrBufferFull = errors.New("Scheduler: buffer is full")

type (
	Opt func(c *Client)

	event struct {
		event flow.Message
		ctx   context.Context
		err   error
		reqId uint64
	}

	// Handles thread-safe communication with Mesos client
	// Supports multiple concurrent read and write requests
	// First push request must be subscribe call
	Client struct {
		bufferSize int

		client mesos.Client

		// client is subscribed when not empty
		streamId string

		// buffer for events from Mesos and for loopback message (not mesos Call)
		buffer chan flow.Message
		// write request queue
		write chan chan *event

		ctx    context.Context
		cancel context.CancelFunc

		log log.Logger
	}

	// Marker interface to get streamId from mesos response
	ResponseWithMesosStreamId interface {
		StreamId() string
	}
)

func WithLogger(l log.Logger) Opt {
	return func(c *Client) {
		c.log = l
	}
}

func WithBufferSize(size int) Opt {
	return func(c *Client) {
		c.bufferSize = size
	}
}

var _ = flow.Flow(&Client{})

func Blueprint(client mesos.Client, opts ...Opt) flow.SinkBlueprint {
	return flow.SinkBlueprintFunc(func(matOpts ...flow.MatOpt) flow.Sink {
		cfg := flow.MatOpts(matOpts).Config()

		if cfg.Log != nil {
			opts = append(opts, WithLogger(log.NewContext(cfg.Log).With("src", "scheduler_client")))
		}

		return New(client, opts...)
	})
}

// Provided mesos.Client must return response implementing ResponseWithMesosStreamId interface
func New(client mesos.Client, opts ...Opt) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		bufferSize: 16,
		ctx:        ctx,
		cancel:     cancel,
		client:     client,
		write:      make(chan chan *event),
		log:        log.NewNopLogger(),
	}

	for _, o := range opts {
		o(c)
	}

	c.buffer = make(chan flow.Message, c.bufferSize)

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
	// read messages from buffer to allow pull after context cancelled
	select {
	case ev := <-c.buffer:
		return ev, nil
	default:
	}

	select {
	case ev := <-c.buffer:
		return ev, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	}
}

// Pull/Push will return context.ContextCancelled after close
func (c *Client) Close() error {
	c.cancel()
	c.log.Log("event", "closed")
	return nil
}

// implements flow.Sink interface
func (c *Client) IsSink() {}

func (c *Client) schedulerLoop() {
	for {
		select {
		case <-c.ctx.Done():
			c.log.Log("event", "main_loop_cancelled")
			return
		case w := <-c.write:
			// write requests are dispatched asynchronously except first call (subscribe)
			if c.streamId == "" {
				resp, err := c.subscribe(<-w)
				if err != nil {
					c.cancel()
				} else {
					c.log.Log("event", "subscribed", "stream-id", c.streamId)
					go c.readerLoop(resp)
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

func (c *Client) subscribe(ev *event) (mesos.Response, error) {
	isSubscribe, subscribe := IsSubscribeMessage(ev.event)

	if !isSubscribe {
		return nil, errors.Errorf("First message must be subscribe - but is %T", ev.event)
	}

	// set Close(true) because subscribe call cannot reuse connections
	if res, err := c.client.Do(subscribe, c.ctx, mesos.WithClose(true)); err != nil {
		return nil, err
	} else {

		s, ok := res.(ResponseWithMesosStreamId)
		if !ok {
			return nil, errors.Errorf("Client response does not implement MesosStreamId interface, got type %T", res)
		}

		sid := s.StreamId()
		if sid == "" {
			return nil, errors.Errorf("Mesos is expected to set %v but it is empty", mesos.MESOS_STREAM_ID_HEADER)
		}
		c.streamId = sid
		return res, nil
	}
}

func (c *Client) readerLoop(resp mesos.Response) {
	defer c.log.Log("event", "reader_loop_cancelled")
	defer resp.Close()

	for {
		msg := &Event{}
		err := resp.Read(msg)
		if err != nil {
			c.log.Log("event", "reader_loop", "err", err)
			c.cancel()
			return
		} else {
			select {
			case c.buffer <- flow.Message(msg):
			case <-c.ctx.Done():
				return
			}
		}
	}
}
