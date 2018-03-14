package scheduler

import (
	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
)

func Subscribe(info *mesos.FrameworkInfo) *Call {
	return &Call{
		Type:        Call_SUBSCRIBE.Enum(),
		FrameworkId: info.GetId(),
		Subscribe:   &Call_Subscribe{FrameworkInfo: info},
	}
}

func Teardown() *Call {
	return &Call{
		Type: Call_TEARDOWN.Enum(),
	}
}

func AcceptOffer(offerId *mesos.OfferID, ops ...*mesos.Offer_Operation) *Call {
	return &Call{
		Type: Call_ACCEPT.Enum(),
		Accept: &Call_Accept{
			OfferIds:   []*mesos.OfferID{offerId},
			Operations: ops,
		},
	}
}

func AcceptInverseOffers(offerIDs ...*mesos.OfferID) *Call {
	return &Call{
		Type: Call_ACCEPT_INVERSE_OFFERS.Enum(),
		AcceptInverseOffers: &Call_AcceptInverseOffers{
			InverseOfferIds: offerIDs,
		},
	}
}

func DeclineInverseOffers(offerIDs ...*mesos.OfferID) *Call {
	return &Call{
		Type: Call_DECLINE_INVERSE_OFFERS.Enum(),
		DeclineInverseOffers: &Call_DeclineInverseOffers{
			InverseOfferIds: offerIDs,
		},
	}
}

func OpLaunch(ti ...*mesos.TaskInfo) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_LAUNCH.Enum(),
		Launch: &mesos.Offer_Operation_Launch{
			TaskInfos: ti,
		},
	}
}

func OpLaunchGroup(exec *mesos.ExecutorInfo, ti ...*mesos.TaskInfo) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_LAUNCH_GROUP.Enum(),
		LaunchGroup: &mesos.Offer_Operation_LaunchGroup{
			Executor: exec,
			TaskGroup: &mesos.TaskGroupInfo{
				Tasks: ti,
			},
		},
	}
}

func OpReserve(rs ...*mesos.Resource) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_RESERVE.Enum(),
		Reserve: &mesos.Offer_Operation_Reserve{
			Resources: rs,
		},
	}
}

func OpUnreserve(rs ...*mesos.Resource) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_UNRESERVE.Enum(),
		Unreserve: &mesos.Offer_Operation_Unreserve{
			Resources: rs,
		},
	}
}

func OpCreate(rs ...*mesos.Resource) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_CREATE.Enum(),
		Create: &mesos.Offer_Operation_Create{
			Volumes: rs,
		},
	}
}

func OpDestroy(rs ...*mesos.Resource) *mesos.Offer_Operation {
	return &mesos.Offer_Operation{
		Type: mesos.Offer_Operation_DESTROY.Enum(),
		Destroy: &mesos.Offer_Operation_Destroy{
			Volumes: rs,
		},
	}
}

func Revive(roles ...string) *Call {
	return &Call{
		Type: Call_REVIVE.Enum(),
		Revive: &Call_Revive{
			Roles: roles,
		},
	}
}

func Suppress(roles ...string) *Call {
	return &Call{
		Type: Call_SUPPRESS.Enum(),
		Suppress: &Call_Suppress{
			Roles: roles,
		},
	}
}

func Decline(offerIDs ...*mesos.OfferID) *Call {
	return &Call{
		Type: Call_DECLINE.Enum(),
		Decline: &Call_Decline{
			OfferIds: offerIDs,
		},
	}
}

func Kill(taskID, agentID string) *Call {
	return &Call{
		Type: Call_KILL.Enum(),
		Kill: &Call_Kill{
			TaskId:  &mesos.TaskID{Value: &taskID},
			SlaveId: optionalSlaveID(agentID),
		},
	}
}

func optionalSlaveID(agentID string) *mesos.SlaveID {
	if agentID == "" {
		return nil
	}
	return &mesos.SlaveID{Value: &agentID}
}

func Shutdown(executorID, agentID string) *Call {
	return &Call{
		Type: Call_SHUTDOWN.Enum(),
		Shutdown: &Call_Shutdown{
			ExecutorId: &mesos.ExecutorID{Value: &executorID},
			SlaveId:    &mesos.SlaveID{Value: &agentID},
		},
	}
}

func Acknowledge(agentID, taskID string, uuid []byte) *Call {
	return &Call{
		Type: Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &Call_Acknowledge{
			SlaveId: &mesos.SlaveID{Value: &agentID},
			TaskId:  &mesos.TaskID{Value: &taskID},
			Uuid:    uuid,
		},
	}
}

// Map keys (taskID's) are required to be non-empty, but values (agentID's) *may* be empty.
func ReconcileTasks(tasks map[string]string) *Call {
	var rec []*Call_Reconcile_Task

	for tid, aid := range tasks {
		rec = append(rec, &Call_Reconcile_Task{
			TaskId:  &mesos.TaskID{Value: mesos.Strp(tid)},
			SlaveId: optionalSlaveID(aid),
		})
	}

	return &Call{
		Type: Call_RECONCILE.Enum(),
		Reconcile: &Call_Reconcile{
			Tasks: rec,
		},
	}
}

func Request(requests ...*mesos.Request) *Call {
	return &Call{
		Type: Call_REQUEST.Enum(),
		Request: &Call_Request{
			Requests: requests,
		},
	}
}

func IsSubscribedMessage(e flow.Message) (bool, *Event) {
	switch r := e.(type) {
	case *Event:
		return IsSubscribed(r), r
	}

	return false, nil
}

func IsResourceOffers(e flow.Message) (bool, *Event) {
	switch r := e.(type) {
	case *Event:
		return r.GetType() == Event_OFFERS, r
	}

	return false, nil
}

func IsOfferDecline(e flow.Message) (bool, *Call) {
	switch r := e.(type) {
	case *Call:
		return r.GetType() == Call_DECLINE, r
	}

	return false, nil
}

func IsSubscribed(e *Event) bool {
	return e.GetType() == Event_SUBSCRIBED
}

func IsSubscribeMessage(ev flow.Message) (bool, *Call) {
	switch r := ev.(type) {
	case *Call:
		return r.GetType() == Call_SUBSCRIBE, r
	}
	return false, nil
}

func IsUpdate(e *Event) bool {
	return e.GetType() == Event_UPDATE
}
