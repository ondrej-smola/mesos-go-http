package heartbeat_test

import (
	. "github.com/ondrej-smola/mesos-go-http/scheduler/stage/heartbeat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"testing"
	"time"
)

func TestAck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Heartbeat stage suite")
}

var _ = Describe("FrameworkId stage", func() {

	It("Set deadling from subscribed event", func(done Done) {
		fwId := New(WithMaxMissedHeartbeats(1))
		sink := flow.NewTestFlow()
		fwId.Via(sink)

		go func() {
			defer GinkgoRecover()
			id := "5"
			_, pullReply := sink.ExpectPull()
			s := scheduler.Subscribed(id)
			s.Subscribed.HeartbeatIntervalSeconds = 0.005 // 5 millis
			pullReply.Message(s)

			// this pull times out
			ctx, reply := sink.ExpectPull()
			start := time.Now()
			<-ctx.Done()
			Expect(time.Now().Sub(start)).To(BeNumerically("~", 10*time.Millisecond, time.Millisecond))
			reply.Error(ctx.Err())

			// this pull does not
			_, reply = sink.ExpectPull()
			time.Sleep(5 * time.Millisecond)
			reply.Message(&scheduler.Event{})
		}()

		_, err := fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = fwId.Pull(context.Background())
		Expect(err).To(Equal(context.DeadlineExceeded))

		_, err = fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		close(done)
	})

	It("Set deadling from configuration", func(done Done) {
		fwId := New(WithMaxMissedHeartbeats(3), WithHeartbeatDeadline(10*time.Millisecond))
		sink := flow.NewTestFlow()
		fwId.Via(sink)

		go func() {
			defer GinkgoRecover()
			id := "5"
			_, pullReply := sink.ExpectPull()
			s := scheduler.Subscribed(id)
			s.Subscribed.HeartbeatIntervalSeconds = 0.005 // 5 millis - should be IGNORED
			pullReply.Message(s)

			// this pull times out
			ctx, reply := sink.ExpectPull()
			start := time.Now()
			<-ctx.Done()
			Expect(time.Now().Sub(start)).To(BeNumerically("~", 40*time.Millisecond, time.Millisecond))
			reply.Error(ctx.Err())

			// this pull does not
			_, reply = sink.ExpectPull()
			time.Sleep(30 * time.Millisecond)
			reply.Message(&scheduler.Event{})
		}()

		_, err := fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		_, err = fwId.Pull(context.Background())
		Expect(err).To(Equal(context.DeadlineExceeded))

		_, err = fwId.Pull(context.Background())
		Expect(err).To(Succeed())

		close(done)
	})
})
