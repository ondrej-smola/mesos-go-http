package scheduler

import (
	"github.com/ondrej-smola/mesos-go-http/lib"
)

func TestFrameworkInfo() *mesos.FrameworkInfo {
	return &mesos.FrameworkInfo{
		Name: mesos.Strp("test"),
		User: mesos.Strp("user"),
	}
}

func TestSubscribed(fwId string) *Event {
	return &Event{
		Type: Event_SUBSCRIBED.Enum(),
		Subscribed: &Event_Subscribed{
			FrameworkId: &mesos.FrameworkID{Value: &fwId},
			MasterInfo: &mesos.MasterInfo{
				Id: mesos.Strp("test_masters"),
				Address: &mesos.Address{
					Ip:   mesos.Strp("127.0.0.1"),
					Port: mesos.I32p(5050),
				},
			},
		},
	}
}

func TestHeartbeat() *Event {
	return &Event{
		Type: Event_HEARTBEAT.Enum(),
	}
}
