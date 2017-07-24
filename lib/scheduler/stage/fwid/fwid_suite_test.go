package fwid_test

import (
	. "github.com/ondrej-smola/mesos-go-http/lib/scheduler/stage/fwid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"testing"

	"github.com/ondrej-smola/mesos-go-http/lib/flow"
	"github.com/ondrej-smola/mesos-go-http/lib/scheduler"
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
			Expect(c.FrameworkId.GetValue()).To(Equal(id))
			push.OK()
		}()

		_, err := fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		err = fwId.Push(&scheduler.Call{Type: scheduler.Call_KILL.Enum()}, context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Set frameworkId from configuration", func(done Done) {
		id := "5"
		fwId := New(WithFrameworkId(id))
		sink := flow.NewTestFlow()
		fwId.Via(sink)

		go func() {
			defer GinkgoRecover()

			push := sink.ExpectPush()
			c, ok := push.Msg.(*scheduler.Call)
			Expect(ok).To(BeTrue())
			Expect(c.FrameworkId.GetValue()).To(Equal(id))
			push.OK()
		}()

		err := fwId.Push(&scheduler.Call{Type: scheduler.Call_KILL.Enum()}, context.Background())
		Expect(err).To(Succeed())

		close(done)
	})
})
