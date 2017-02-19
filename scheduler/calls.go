package scheduler

import (
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
)

func Filters(fo ...mesos.FilterOpt) CallOpt {
	return func(c *Call) {
		switch c.Type {
		case Call_ACCEPT:
			c.Accept.Filters = mesos.OptionalFilters(fo...)
		case Call_ACCEPT_INVERSE_OFFERS:
			c.AcceptInverseOffers.Filters = mesos.OptionalFilters(fo...)
		case Call_DECLINE:
			c.Decline.Filters = mesos.OptionalFilters(fo...)
		case Call_DECLINE_INVERSE_OFFERS:
			c.DeclineInverseOffers.Filters = mesos.OptionalFilters(fo...)
		}
	}
}

func KillPolicy(ko ...mesos.KillOpt) CallOpt {
	return func(c *Call) {
		switch c.Type {
		case Call_KILL:
			c.Kill.KillPolicy = mesos.OptionalKillPolicy(ko...)
		}
	}
}

func FrameworkId(id mesos.FrameworkID) CallOpt {
	return func(c *Call) {
		c.FrameworkID = &id
	}
}

type acceptBuilder struct {
	offerIDs   map[mesos.OfferID]struct{}
	operations []mesos.Offer_Operation
	filters    *mesos.Filters
}

type AcceptOpt func(*acceptBuilder)

type OfferOperations []mesos.Offer_Operation

// WithOffers allows a client to pair some set of OfferOperations with multiple resource offers.
// Example: calls.Accept(calls.OfferOperations{calls.OpLaunch(tasks...)}.WithOffers(offers...))
func (ob OfferOperations) WithOffers(ids ...mesos.OfferID) AcceptOpt {
	return func(ab *acceptBuilder) {
		for i := range ids {
			ab.offerIDs[ids[i]] = struct{}{}
		}
		ab.operations = append(ab.operations, ob...)
	}
}

// Accept returns an accept call with the given parameters.
// Callers are expected to fill in the FrameworkID and Filters.
func Accept(ops ...AcceptOpt) *Call {
	ab := &acceptBuilder{
		offerIDs: make(map[mesos.OfferID]struct{}, len(ops)),
	}
	for _, op := range ops {
		op(ab)
	}
	offerIDs := make([]mesos.OfferID, 0, len(ab.offerIDs))
	for id := range ab.offerIDs {
		offerIDs = append(offerIDs, id)
	}
	return &Call{
		Type: Call_ACCEPT,
		Accept: &Call_Accept{
			OfferIDs:   offerIDs,
			Operations: ab.operations,
			Filters:    ab.filters,
		},
	}
}

// Subscribe returns a subscribe call with the given parameters.
// The call's FrameworkID is automatically filled in from the info specification.
func Subscribe(info mesos.FrameworkInfo) *Call {
	return &Call{
		Type:        Call_SUBSCRIBE,
		FrameworkID: info.GetID(),
		Subscribe:   &Call_Subscribe{FrameworkInfo: info},
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
		return r.Type == Event_OFFERS, r
	}

	return false, nil
}

func IsOfferDecline(e flow.Message) (bool, *Call) {
	switch r := e.(type) {
	case *Call:
		return r.Type == Call_DECLINE, r
	}

	return false, nil
}

func IsSubscribed(e *Event) bool {
	return e.Type == Event_SUBSCRIBED
}

func IsSubscribeMessage(ev flow.Message) (bool, *Call) {
	switch r := ev.(type) {
	case *Call:
		return r.Type == Call_SUBSCRIBE, r
	}
	return false, nil
}

func IsUpdate(e *Event) bool {
	return e.Type == Event_UPDATE
}

func Teardown() *Call {
	return &Call{
		Type: Call_TEARDOWN,
	}
}

// AcceptInverseOffers returns an accept-inverse-offers call for the given offer IDs.
// Callers are expected to fill in the FrameworkID and Filters.
func AcceptInverseOffers(offerIDs ...mesos.OfferID) *Call {
	return &Call{
		Type: Call_ACCEPT_INVERSE_OFFERS,
		AcceptInverseOffers: &Call_AcceptInverseOffers{
			InverseOfferIDs: offerIDs,
		},
	}
}

// DeclineInverseOffers returns a decline-inverse-offers call for the given offer IDs.
// Callers are expected to fill in the FrameworkID and Filters.
func DeclineInverseOffers(offerIDs ...mesos.OfferID) *Call {
	return &Call{
		Type: Call_DECLINE_INVERSE_OFFERS,
		DeclineInverseOffers: &Call_DeclineInverseOffers{
			InverseOfferIDs: offerIDs,
		},
	}
}

// OpLaunch returns a launch operation builder for the given tasks
func OpLaunch(ti ...mesos.TaskInfo) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_LAUNCH,
		Launch: &mesos.Offer_Operation_Launch{
			TaskInfos: ti,
		},
	}
}

func OpLaunchGroup(exec mesos.ExecutorInfo, ti ...mesos.TaskInfo) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_LAUNCH_GROUP,
		LaunchGroup: &mesos.Offer_Operation_LaunchGroup{
			Executor: exec,
			TaskGroup: mesos.TaskGroupInfo{
				Tasks: ti,
			},
		},
	}
}

func OpReserve(rs ...mesos.Resource) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_RESERVE,
		Reserve: &mesos.Offer_Operation_Reserve{
			Resources: rs,
		},
	}
}

func OpUnreserve(rs ...mesos.Resource) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_UNRESERVE,
		Unreserve: &mesos.Offer_Operation_Unreserve{
			Resources: rs,
		},
	}
}

func OpCreate(rs ...mesos.Resource) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_CREATE,
		Create: &mesos.Offer_Operation_Create{
			Volumes: rs,
		},
	}
}

func OpDestroy(rs ...mesos.Resource) mesos.Offer_Operation {
	return mesos.Offer_Operation{
		Type: mesos.Offer_Operation_DESTROY,
		Destroy: &mesos.Offer_Operation_Destroy{
			Volumes: rs,
		},
	}
}

// Revive returns a revive call.
// Callers are expected to fill in the FrameworkID.
func Revive() *Call {
	return &Call{
		Type: Call_REVIVE,
	}
}

// Suppress returns a suppress call.
// Callers are expected to fill in the FrameworkID.
func Suppress() *Call {
	return &Call{
		Type: Call_SUPPRESS,
	}
}

// Decline returns a decline call with the given parameters.
// Callers are expected to fill in the FrameworkID and Filters.
func Decline(offerIDs ...mesos.OfferID) *Call {
	return &Call{
		Type: Call_DECLINE,
		Decline: &Call_Decline{
			OfferIDs: offerIDs,
		},
	}
}

// Kill returns a kill call with the given parameters.
// Callers are expected to fill in the FrameworkID.
func Kill(taskID, agentID string) *Call {
	return &Call{
		Type: Call_KILL,
		Kill: &Call_Kill{
			TaskID:  mesos.TaskID{Value: taskID},
			AgentID: optionalAgentID(agentID),
		},
	}
}

// Shutdown returns a shutdown call with the given parameters.
// Callers are expected to fill in the FrameworkID.
func Shutdown(executorID, agentID string) *Call {
	return &Call{
		Type: Call_SHUTDOWN,
		Shutdown: &Call_Shutdown{
			ExecutorID: mesos.ExecutorID{Value: executorID},
			AgentID:    mesos.AgentID{Value: agentID},
		},
	}
}

// Acknowledge returns an acknowledge call with the given parameters.
// Callers are expected to fill in the FrameworkID.
func Acknowledge(agentID, taskID string, uuid []byte) *Call {
	return &Call{
		Type: Call_ACKNOWLEDGE,
		Acknowledge: &Call_Acknowledge{
			AgentID: mesos.AgentID{Value: agentID},
			TaskID:  mesos.TaskID{Value: taskID},
			UUID:    uuid,
		},
	}
}

// ReconcileTasks constructs a reconcile scheduler call from the given mappings:
//     map[string]string{taskID:agentID}
// Map keys (taskID's) are required to be non-empty, but values (agentID's) *may* be empty.
func ReconcileTasks(tasks map[string]string) *Call {
	rec := []Call_Reconcile_Task{}

	for tid, aid := range tasks {
		rec = append(rec, Call_Reconcile_Task{
			TaskID:  mesos.TaskID{Value: tid},
			AgentID: optionalAgentID(aid),
		})
	}

	return &Call{
		Type: Call_RECONCILE,
		Reconcile: &Call_Reconcile{
			Tasks: rec,
		},
	}
}

// Message returns a message call with the given parameters.
// Callers are expected to fill in the FrameworkID.
func Message(agentID, executorID string, data []byte) *Call {
	return &Call{
		Type: Call_MESSAGE,
		Message: &Call_Message{
			AgentID:    mesos.AgentID{Value: agentID},
			ExecutorID: mesos.ExecutorID{Value: executorID},
			Data:       data,
		},
	}
}

// Request returns a resource request call with the given parameters.
// Callers are expected to fill in the FrameworkID.
func Request(requests ...mesos.Request) *Call {
	return &Call{
		Type: Call_REQUEST,
		Request: &Call_Request{
			Requests: requests,
		},
	}
}

func optionalAgentID(agentID string) *mesos.AgentID {
	if agentID == "" {
		return nil
	}
	return &mesos.AgentID{Value: agentID}
}
