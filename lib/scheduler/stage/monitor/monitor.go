package monitor

import (
	"context"
	"time"

	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
)

type (
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
		ResourceOffered(name string, role string, value float64)
	}

	metrics struct {
		monit Monitor
		via   flow.Flow
	}
)

func Blueprint(b Monitor) flow.StageBlueprint {
	return flow.StageBlueprintFunc(func(...flow.MatOpt) flow.Stage {
		return New(b)
	})
}

func New(monitor Monitor) flow.Stage {
	return &metrics{monit: monitor}
}

func (m *metrics) Push(ev flow.Message, ctx context.Context) error {
	name := ev.Name()
	m.monit.Push(name)

	if ok, _ := scheduler.IsOfferDecline(ev); ok {
		m.monit.OffersDeclined(1)
	}

	start := time.Now()
	err := m.via.Push(ev, ctx)
	m.monit.PushLatency(time.Now().Sub(start))

	if err != nil {
		m.monit.PushErr(err)
	}

	return err
}

func (m *metrics) Pull(ctx context.Context) (flow.Message, error) {
	start := time.Now()
	msg, err := m.via.Pull(ctx)
	m.monit.PullLatency(time.Now().Sub(start))

	if err == nil {
		m.monit.Pull(msg.Name())
		if ok, o := scheduler.IsResourceOffers(msg); ok {
			offers := o.Offers.Offers
			m.monit.OffersReceived(uint32(len(offers)))
			for _, o := range offers {
				for _, res := range o.Resources {
					// only scalar resources are supported
					if res.Scalar != nil {
						m.monit.ResourceOffered(res.GetName(), res.GetRole(), res.Scalar.GetValue())
					}
				}
			}
		}
	} else {
		m.monit.PullErr(err)
	}

	return msg, err
}

func (m *metrics) Via(f flow.Flow) {
	m.via = f
}

func (m *metrics) Close() error {
	return m.via.Close()
}
