package scheduler

import (
	"github.com/ondrej-smola/mesos-go-http/flow"
)

func (e *Call) IsFlowMessage() {}
func (e *Call) Name() string   { return e.Type.String() }

var _ = MonitoredMessage(&Event{})
var _ = flow.Message(&Event{})

func (e *Event) IsFlowMessage() {}
func (e *Event) Name() string   { return e.Type.String() }

var _ = MonitoredMessage(&Event{})
var _ = flow.Message(&Event{})

// Custom message useful for sending something as message
type PingMessage struct{}

func (p *PingMessage) IsFlowMessage() {}
func (p *PingMessage) Name() string {
	return "ping"
}
