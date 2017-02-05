package ack

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"github.com/pkg/errors"
	"sync"
)

type (
	Opt func(c *Acks)

	Acks struct {
		via             flow.Flow
		frameworkId     string
		failOnFailedAck bool
		log             log.Logger
		sync.RWMutex
	}
)

func WithLogger(l log.Logger) Opt {
	return func(c *Acks) {
		c.log = l
	}
}

func WithFailOnFailedAck() Opt {
	return func(c *Acks) {
		c.failOnFailedAck = true
	}
}

func Blueprint(opts ...Opt) flow.StageBlueprint {
	return func(matOpts ...flow.MatOpt) flow.Stage {
		cfg := flow.MatOpts(matOpts).Config()
		if cfg.Log != nil {
			opts = append(opts, WithLogger(log.NewContext(cfg.Log).With("stage", "heartbeats")))
		}
		return New(opts...)
	}
}

// Automatically acknowledges all pulled update requests.
// FrameworkId is set from subscribe call during pull.
// By default only logs failed acks and returns pulled update message
func New(opts ...Opt) *Acks {
	a := &Acks{
		failOnFailedAck: false,
		log:             log.NewNopLogger(),
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

var _ = flow.Stage(&Acks{})

func (i *Acks) Push(ev flow.Message, ctx context.Context) error {
	return i.via.Push(ev, ctx)
}

func (i *Acks) Pull(ctx context.Context) (flow.Message, error) {
	ev, err := i.via.Pull(ctx)

	if err != nil {
		return nil, err
	}

	switch e := ev.(type) {
	case *scheduler.Event:
		if scheduler.IsSubscribed(e) {
			i.Lock()
			i.frameworkId = e.Subscribed.ID.Value
			i.Unlock()
		} else if scheduler.IsUpdate(e) {
			i.RLock()
			fwId := i.frameworkId
			i.RUnlock()

			if fwId == "" {
				return nil, errors.New("Update message received but not yet subscribed")
			}

			state := e.Update.Status
			if len(state.UUID) > 0 {
				call := scheduler.Acknowledge(state.AgentID.Value, state.TaskID.Value, state.UUID).
					With(scheduler.FrameworkId(mesos.FrameworkID{Value: fwId}))

				if err := i.via.Push(call, ctx); err != nil {
					if i.failOnFailedAck {
						return nil, errors.Errorf("Failed to send implicit ack for %v cause: %v", e, err)
					} else {
						i.log.Log(
							"event", "ack_failed",
							"task", state.TaskID.Value,
							"agent", state.AgentID.Value,
							"state", state.State.Enum(),
							"err", err,
						)
					}
				}
				return e, nil
			}
		}
	}

	return ev, nil
}

func (i *Acks) Via(f flow.Flow) {
	i.via = f
}

func (i *Acks) Name() string {
	return "implicit_ack"
}

func (i *Acks) Close() error {
	return i.via.Close()
}
