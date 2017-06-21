package fwid

import (
	"context"
	"sync"

	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
)

type (
	Opt func(c *FwId)

	FwId struct {
		via         flow.Flow
		frameworkId string
		sync.RWMutex
	}
)

func WithFrameworkId(id string) Opt {
	return func(c *FwId) {
		c.frameworkId = id
	}
}

func Blueprint(opts ...Opt) flow.StageBlueprint {
	return flow.StageBlueprintFunc(func(matOpts ...flow.MatOpt) flow.Stage {
		return New(opts...)
	})
}

// Sets framework id from subscribe call on all following calls
func New(opts ...Opt) flow.Stage {
	cfg := &FwId{}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

func (h *FwId) Push(ev flow.Message, ctx context.Context) error {
	h.RLock()
	fwId := h.frameworkId
	h.RUnlock()

	if fwId != "" {
		id := &mesos.FrameworkID{Value: mesos.Strp(fwId)}
		switch e := ev.(type) {
		case *scheduler.Call:
			e.FrameworkId = id
			switch e.GetType() {
			case scheduler.Call_ACCEPT:
				ops := e.Accept.Operations
				for i, op := range ops {
					if op.GetType() == mesos.Offer_Operation_LAUNCH_GROUP {
						ops[i].LaunchGroup.Executor.FrameworkId = id
					}
				}
			case scheduler.Call_SUBSCRIBE:
				e.Subscribe.FrameworkInfo.Id = id
			}
		}
	}

	return h.via.Push(ev, ctx)
}

func (h *FwId) Pull(ctx context.Context) (flow.Message, error) {
	if msg, err := h.via.Pull(ctx); err != nil {
		return msg, err
	} else {
		is, e := scheduler.IsSubscribedMessage(msg)
		if is {
			h.Lock()
			h.frameworkId = e.Subscribed.FrameworkId.GetValue()
			h.Unlock()
		}
		return msg, err
	}
}

func (h *FwId) Via(f flow.Flow) {
	h.via = f
}

func (h *FwId) Close() error {
	return h.via.Close()
}
