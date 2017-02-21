package scheduler

import (
	"time"
)

type (
	MonitoredMessage interface {
		Name() string
	}

	Monitor interface {
		// called for every push, name is empty when message does not implement MonitoredMessage interface
		Push(name string)
		// called for every pull, name is empty when message does not implement MonitoredMessage interface
		Pull(name string)

		// called when push failed
		PushErr(err error)
		// called when pull failed
		PullErr(err error)

		PushLatency(time time.Duration)
		PullLatency(time time.Duration)

		OffersReceived(count uint32)
		OffersDeclined(count uint32)

		// called when pulled message that contains resource offers, only for scalar resources
		ResourceOffered(name string, value float64)
	}
)
