package scheduler

import "github.com/ondrej-smola/mesos-go-http"

func Subscribed(fwId string) *Event {
	return &Event{
		Type: Event_SUBSCRIBED,
		Subscribed: &Event_Subscribed{
			ID: mesos.FrameworkID{Value: fwId},
			HeartbeatIntervalSeconds: 15,
			MasterInfo: &mesos.MasterInfo{
				ID:       "test_masters",
				IP:       2130706433, //127.0.0.1
				Hostname: "localhost",
				Version:  "1.1.0",
			},
		},
	}

}
