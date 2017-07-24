package scheduler

import (
	"time"

	"github.com/ondrej-smola/mesos-go-http/lib"
)

type CallOpt func(*Call)

func (c *Call) With(opts ...CallOpt) *Call {
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func RefuseSeconds(dur time.Duration) CallOpt {
	filter := &mesos.Filters{
		RefuseSeconds: mesos.F64p(dur.Seconds()),
	}

	return func(c *Call) {
		switch c.GetType() {
		case Call_ACCEPT:
			c.Accept.Filters = filter.Clone()
		case Call_ACCEPT_INVERSE_OFFERS:
			c.AcceptInverseOffers.Filters = filter.Clone()
		case Call_DECLINE:
			c.Decline.Filters = filter.Clone()
		case Call_DECLINE_INVERSE_OFFERS:
			c.DeclineInverseOffers.Filters = filter.Clone()
		}
	}
}

func KillPolicy(dur time.Duration) CallOpt {
	policy := &mesos.KillPolicy{
		GracePeriod: &mesos.DurationInfo{Nanoseconds: mesos.I64p(dur.Nanoseconds())},
	}

	return func(c *Call) {
		switch c.GetType() {
		case Call_KILL:
			c.Kill.KillPolicy = policy.Clone()
		}
	}
}
