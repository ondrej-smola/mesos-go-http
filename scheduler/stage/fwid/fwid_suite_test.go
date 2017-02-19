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
			pull := sink.ExpectPull()
			pull.Message(scheduler.TestSubscribed(id))

			push := sink.ExpectPush()
			c, ok := push.Msg.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.FrameworkID.Value).To(Equal(id))
			push.OK()
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

			push := sink.ExpectPush()
			c, ok := push.Msg.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.FrameworkID.Value).To(Equal(id))
			push.OK()
		}()

		err := fwId.Push(&scheduler.Call{Type: scheduler.Call_KILL}, context.Background())
		Expect(err).To(Succeed())

		close(done)
	})
})
