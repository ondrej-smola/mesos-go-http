package monitor_test

import (
	. "github.com/ondrej-smola/mesos-go-http/scheduler/stage/monitor"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/flow"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestAck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Monitor stage suite")
}

var _ = Describe("Monitor", func() {
	It("Push", func(done Done) {
		ctx := context.Background()
		sink := flow.NewTestFlow()
		testMonit := scheduler.NewTestMonitor()
		m := New(testMonit)
		m.Via(sink)

		go func() {
			push := sink.ExpectPush()
			time.Sleep(5 * time.Millisecond)
			push.Error(errors.New("boom!"))

			push = sink.ExpectPush()
			time.Sleep(5 * time.Millisecond)
			push.OK()

			push = sink.ExpectPush()
			time.Sleep(5 * time.Millisecond)
			push.OK()
		}()

		ping := &scheduler.PingMessage{}
		decline := &scheduler.Call{
			Type: scheduler.Call_DECLINE,
		}
		m.Push(ping, ctx)
		m.Push(decline, ctx)
		m.Push(ping, ctx)

		testMonit.Lock()
		Expect(testMonit.PushC[ping.Name()]).To(BeEquivalentTo(2))
		Expect(testMonit.PushC[decline.Name()]).To(BeEquivalentTo(1))
		Expect(testMonit.OffersDeclinedC).To(BeEquivalentTo(1))
		Expect(testMonit.OffersDeclinedC).To(BeEquivalentTo(1))
		Expect(testMonit.PushErrC).To(BeEquivalentTo(1))
		Expect(testMonit.PushLatencyC).To(BeNumerically(">=", 15*time.Millisecond))
		testMonit.Unlock()

		close(done)
	})

	It("Pull", func(done Done) {
		ctx := context.Background()
		sink := flow.NewTestFlow()
		testMonit := scheduler.NewTestMonitor()
		m := New(testMonit)
		m.Via(sink)

		cpus := mesos.BuildResource().Name("cpus").Scalar(1).Build()
		mem := mesos.BuildResource().Name("mem").Scalar(512).Build()

		offers := &scheduler.Event{
			Type: scheduler.Event_OFFERS,
			Offers: &scheduler.Event_Offers{
				Offers: []mesos.Offer{
					{Resources: []mesos.Resource{cpus, mem}},
					{Resources: []mesos.Resource{mem}},
				},
			},
		}

		ping := &scheduler.PingMessage{}

		go func() {
			pull := sink.ExpectPull()
			time.Sleep(5 * time.Millisecond)
			pull.Error(errors.New("boom!"))

			pull = sink.ExpectPull()
			time.Sleep(5 * time.Millisecond)
			pull.Message(offers)

			pull = sink.ExpectPull()
			time.Sleep(5 * time.Millisecond)
			pull.Message(ping)
		}()

		m.Pull(ctx)
		m.Pull(ctx)
		m.Pull(ctx)

		testMonit.Lock()
		Expect(testMonit.PullC[ping.Name()]).To(BeEquivalentTo(1))
		Expect(testMonit.PullC[offers.Name()]).To(BeEquivalentTo(1))
		Expect(testMonit.OffersC).To(BeEquivalentTo(2))
		Expect(testMonit.PullErrC).To(BeEquivalentTo(1))
		Expect(testMonit.PullLatencyC).To(BeNumerically(">=", 15*time.Millisecond))
		Expect(testMonit.ResourcesC[cpus.Name]).To(BeEquivalentTo(1))
		Expect(testMonit.ResourcesC[mem.Name]).To(BeEquivalentTo(1024))
		testMonit.Unlock()

		close(done)
	})
})
