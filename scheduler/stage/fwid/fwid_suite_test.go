package fwid_test

import (
	. "github.com/ondrej-smola/mesos-go-http/scheduler/stage/fwid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"testing"
)

func TestAck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FrameworkId stage suite")
}

var _ = Describe("FrameworkId stage", func() {

	It("Set frameworkId from subscribe call", func(done Done) {
		fwId := New()
		sink := flow.NewTestFlow()
		fwId.Via(sink)

		go func() {
			defer GinkgoRecover()
			id := "5"
			_, pullReply := sink.ExpectPull()
			pullReply.Message(scheduler.Subscribed(id))

			m, _, reply := sink.ExpectPush()
			c, ok := m.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.FrameworkID.Value).To(Equal(id))
			reply.OK()
		}()

		_, err := fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		err = fwId.Push(&scheduler.Call{Type: scheduler.Call_KILL}, context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Set frameworkId from configuration", func(done Done) {
		id := "5"
		fwId := New(WithFrameworkId(mesos.FrameworkID{Value: id}))
		sink := flow.NewTestFlow()
		fwId.Via(sink)

		go func() {
			defer GinkgoRecover()

			m, _, reply := sink.ExpectPush()
			c, ok := m.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.FrameworkID.Value).To(Equal(id))
			reply.OK()
		}()

		err := fwId.Push(&scheduler.Call{Type: scheduler.Call_KILL}, context.Background())
		Expect(err).To(Succeed())

		close(done)
	})
})
