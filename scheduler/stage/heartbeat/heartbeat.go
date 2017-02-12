package heartbeat

import (
	"context"
	"fmt"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/log"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"time"
)

type (
	Opt func(c *Heartbeats)

	Heartbeats struct {
		maxMissed         int64
		heartbeatDeadline time.Duration

		via flow.Flow
		log log.Logger
	}
)

const INITIAL_DEADLINE = 15 * time.Second
const DEFAULT_MAX_MISSED = 1

func WithMaxMissedHeartbeats(max uint64) Opt {
	return func(c *Heartbeats) {
		c.maxMissed = int64(max)
	}
}

func WithLogger(l log.Logger) Opt {
	return func(c *Heartbeats) {
		c.log = l
	}
}

func WithHeartbeatDeadline(d time.Duration) Opt {
	if d <= 0 {
		panic(fmt.Sprintf("Deadline must be > 0, is %v", d))
	}

	return func(c *Heartbeats) {
		c.heartbeatDeadline = d
	}
}

func Blueprint(opts ...Opt) flow.StageBlueprint {
	return flow.StageBlueprintFunc(func(matOpts ...flow.MatOpt) flow.Stage {
		cfg := flow.MatOpts(matOpts).Config()

		if cfg.Log != nil {
			opts = append(opts, WithLogger(log.NewContext(cfg.Log).With("stage", "heartbeats")))
		}
		return New(opts...)
	})
}

// Set deadline for pull request based on configuration.
// Also sets initial deadline for subscribe call.
// When no deadline is configured - it is set from subscribed event.
func New(opts ...Opt) *Heartbeats {
	h := &Heartbeats{
		maxMissed: DEFAULT_MAX_MISSED,
		log:       log.NewNopLogger(),
	}

	for _, o := range opts {
		o(h)
	}

	h.heartbeatDeadline = h.heartbeatDeadline * time.Duration(h.maxMissed+1)

	return h
}

var _ = flow.Stage(&Heartbeats{})

func (h *Heartbeats) Push(ev flow.Message, ctx context.Context) error {
	return h.via.Push(ev, ctx)
}

func (h *Heartbeats) Pull(ctx context.Context) (flow.Message, error) {
	deadlineSet := true
	deadline := h.heartbeatDeadline

	if deadline == 0 {
		deadlineSet = false
		deadline = INITIAL_DEADLINE
	}

	deadlinedCtx, cancel := context.WithTimeout(ctx, deadline)
	defer cancel()

	ev, err := h.via.Pull(deadlinedCtx)

	if err == nil {
		switch e := ev.(type) {
		case *scheduler.Event:
			if scheduler.IsSubscribed(e) && !deadlineSet {
				// use precision up to milliseconds
				tmp := int64(e.Subscribed.HeartbeatIntervalSeconds*1000) * (h.maxMissed + 1)
				deadline := time.Duration(tmp) * time.Millisecond
				h.log.Log("event", "heartbeat_set", "deadline", deadline)
				h.heartbeatDeadline = deadline
			}
		}
	}

	return ev, err
}

func (h *Heartbeats) Via(f flow.Flow) {
	h.via = f
}

func (h *Heartbeats) Close() error {
	return h.via.Close()
}
