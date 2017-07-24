package scheduler

import (
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
)

func (e *Call) Name() string {
	return e.Type.String()
}

var _ = flow.Message(&Event{})

func (e *Event) Name() string {
	return e.Type.String()
}

var _ = flow.Message(&Event{})

// Custom message useful for sending something as message
type PingMessage struct{}

func (p *PingMessage) Name() string {
	return "ping"
}

var _ = flow.Message(&PingMessage{})
