package ack

import (
	"context"
	"sync"

	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/log"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
	"github.com/pkg/errors"
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
	return flow.StageBlueprintFunc(func(matOpts ...flow.MatOpt) flow.Stage {
		cfg := flow.MatOpts(matOpts).Config()
		if cfg.Log != nil {
			opts = append(opts, WithLogger(log.With(cfg.Log, "src", "implicit_ack_stage")))
		}
		return New(opts...)
	})
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

func (a *Acks) Push(ev flow.Message, ctx context.Context) error {
	return a.via.Push(ev, ctx)
}

func (a *Acks) Pull(ctx context.Context) (flow.Message, error) {
	ev, err := a.via.Pull(ctx)

	if err != nil {
		return nil, err
	}

	switch e := ev.(type) {
	case *scheduler.Event:
		if scheduler.IsSubscribed(e) {
			a.Lock()
			a.frameworkId = e.Subscribed.FrameworkId.GetValue()
			a.Unlock()
		} else if scheduler.IsUpdate(e) {
			a.RLock()
			fwId := a.frameworkId
			a.RUnlock()

			if fwId == "" {
				return nil, errors.New("Update message received but not yet subscribed")
			}

			state := e.Update.Status

			if state.SlaveId == nil {
				return nil, errors.Errorf("AgentId must be set on status update: %v", state)
			}

			if len(state.Uuid) > 0 {
				call := &scheduler.Call{
					Type: scheduler.Call_ACKNOWLEDGE.Enum(),
					Acknowledge: &scheduler.Call_Acknowledge{
						SlaveId: state.SlaveId,
						TaskId:  state.TaskId,
						Uuid:    state.Uuid,
					},
					FrameworkId: &mesos.FrameworkID{Value: &fwId},
				}

				if err := a.via.Push(call, ctx); err != nil {
					if a.failOnFailedAck {
						return nil, errors.Errorf("Failed to send implicit ack for %v cause: %v", e, err)
					} else {
						a.log.Log(
							"event", "implicit_ack_failed",
							"task", state.TaskId.Value,
							"agent", state.SlaveId.Value,
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

func (a *Acks) Via(f flow.Flow) {
	a.via = f
}

func (a *Acks) Close() error {
	return a.via.Close()
}
