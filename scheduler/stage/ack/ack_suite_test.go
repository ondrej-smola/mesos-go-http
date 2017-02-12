package ack_test

import (
	. "github.com/ondrej-smola/mesos-go-http/scheduler/stage/ack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestAck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ImplicitAck stage suite")
}

var _ = Describe("ImplicitAck state", func() {

	statusUpdate := &scheduler.Event{
		Type: scheduler.Event_UPDATE,
		Update: &scheduler.Event_Update{
			Status: mesos.TaskStatus{
				TaskID:  mesos.TaskID{Value: "1"},
				AgentID: &mesos.AgentID{Value: "agent"},
				UUID:    []byte{1, 2, 3, 4},
				State:   mesos.TaskState_TASK_RUNNING.Enum(),
			},
		},
	}

	It("Ack status update - without agent id", func(done Done) {
		acks := New()
		sink := flow.NewTestFlow()
		acks.Via(sink)

		go func() {
			defer GinkgoRecover()
			_, reply := sink.ExpectPull()
			reply.Message(scheduler.Subscribed("1"))
			_, reply = sink.ExpectPull()
			reply.Message(statusUpdate)
			msg, _, pReply := sink.ExpectPush()
			c, ok := msg.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.Type).To(Equal(scheduler.Call_ACKNOWLEDGE))
			Expect(c.Acknowledge.TaskID.Value).To(Equal(statusUpdate.Update.Status.TaskID.Value))
			Expect(c.Acknowledge.AgentID.Value).To(Equal(statusUpdate.Update.Status.AgentID.Value))
			pReply.OK()
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = acks.Pull(context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Do not ack status update when UUID is not set", func(done Done) {
		acks := New()
		sink := flow.NewTestFlow()
		acks.Via(sink)

		go func() {
			defer GinkgoRecover()
			_, reply := sink.ExpectPull()
			reply.Message(scheduler.Subscribed("1"))
			_, reply = sink.ExpectPull()
			upd := proto.Clone(statusUpdate).(*scheduler.Event)
			upd.Update.Status.UUID = nil
			reply.Message(upd)
			sink.ExpectNoPush(50 * time.Millisecond)
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = acks.Pull(context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Ignore failed status update by default", func(done Done) {
		acks := New()
		sink := flow.NewTestFlow()
		acks.Via(sink)

		go func() {
			defer GinkgoRecover()
			_, reply := sink.ExpectPull()
			reply.Message(scheduler.Subscribed("1"))
			_, reply = sink.ExpectPull()
			reply.Message(statusUpdate)
			_, _, pReply := sink.ExpectPush()
			pReply.Error(errors.New("boom"))
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = acks.Pull(context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Fail on failed ack when configured to do so", func(done Done) {
		acks := New(WithFailOnFailedAck())
		sink := flow.NewTestFlow()
		acks.Via(sink)

		go func() {
			defer GinkgoRecover()
			_, reply := sink.ExpectPull()
			reply.Message(scheduler.Subscribed("1"))
			_, reply = sink.ExpectPull()
			reply.Message(statusUpdate)
			_, _, pReply := sink.ExpectPush()
			pReply.Error(errors.New("boom"))
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = acks.Pull(context.Background())
		Expect(err).To(HaveOccurred())

		close(done)
	})

	It("Fail when status update received before subscribed call", func(done Done) {
		acks := New()
		sink := flow.NewTestFlow()
		acks.Via(sink)

		go func() {
			defer GinkgoRecover()
			_, reply := sink.ExpectPull()
			reply.Message(statusUpdate)
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(HaveOccurred())

		close(done)
	})
})
