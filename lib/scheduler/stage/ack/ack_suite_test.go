package ack_test

import (
	. "github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/ack"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/lib"
	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
	"github.com/pkg/errors"
)

func TestAck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ImplicitAck stage suite")
}

var _ = Describe("ImplicitAck state", func() {

	statusUpdate := &scheduler.Event{
		Type: scheduler.Event_UPDATE.Enum(),
		Update: &scheduler.Event_Update{
			Status: &mesos.TaskStatus{
				TaskId:  &mesos.TaskID{Value: mesos.Strp("1")},
				AgentId: &mesos.AgentID{Value: mesos.Strp("agent")},
				Uuid:    []byte{1, 2, 3, 4},
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
			reply := sink.ExpectPull()
			reply.Message(scheduler.TestSubscribed("1"))
			reply = sink.ExpectPull()
			reply.Message(statusUpdate)
			push := sink.ExpectPush()
			c, ok := push.Msg.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.GetType()).To(Equal(scheduler.Call_ACKNOWLEDGE))
			Expect(c.Acknowledge.TaskId.Value).To(Equal(statusUpdate.Update.Status.TaskId.Value))
			Expect(c.Acknowledge.AgentId.Value).To(Equal(statusUpdate.Update.Status.AgentId.Value))
			push.OK()
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
			pull := sink.ExpectPull()
			pull.Message(scheduler.TestSubscribed("1"))
			pull = sink.ExpectPull()
			upd := proto.Clone(statusUpdate).(*scheduler.Event)
			upd.Update.Status.Uuid = nil
			pull.Message(upd)
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
			pull := sink.ExpectPull()
			pull.Message(scheduler.TestSubscribed("1"))
			pull = sink.ExpectPull()
			pull.Message(statusUpdate)
			push := sink.ExpectPush()
			push.Error(errors.New("boom"))
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
			pull := sink.ExpectPull()
			pull.Message(scheduler.TestSubscribed("1"))
			pull = sink.ExpectPull()
			pull.Message(statusUpdate)
			push := sink.ExpectPush()
			push.Error(errors.New("boom"))
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
			pull := sink.ExpectPull()
			pull.Message(statusUpdate)
		}()

		_, err := acks.Pull(context.Background())
		Expect(err).To(HaveOccurred())

		close(done)
	})
})
