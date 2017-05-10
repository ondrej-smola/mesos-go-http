package monitor

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"time"
)

type (
	metrics struct {
		monit scheduler.Monitor
		via   flow.Flow
	}
)

func Blueprint(b scheduler.Monitor) flow.StageBlueprint {
	return flow.StageBlueprintFunc(func(...flow.MatOpt) flow.Stage {
		return New(b)
	})
}

func New(monitor scheduler.Monitor) flow.Stage {
	return &metrics{monit: monitor}
}

func (m *metrics) Push(ev flow.Message, ctx context.Context) error {
	name := ""
	n, ok := ev.(scheduler.MonitoredMessage)
	if ok {
		name = n.Name()
	}
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
		name := ""
		n, ok := msg.(scheduler.MonitoredMessage)
		if ok {
			name = n.Name()
		}
		m.monit.Pull(name)
		if ok, o := scheduler.IsResourceOffers(msg); ok {
			offers := o.Offers.Offers
			m.monit.OffersReceived(uint32(len(offers)))
			for _, o := range offers {
				for _, res := range o.Resources {
					// only scalar resources are supported
					if res.Scalar != nil {
						m.monit.ResourceOffered(res.GetName(), res.Scalar.Value)
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
