package scheduler_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ondrej-smola/mesos-go-http/lib/client"
	. "github.com/ondrej-smola/mesos-go-http/lib/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scheduler Suite")
}

var _ = Describe("Scheduler", func() {

	var subscribe *Call

	BeforeEach(func() {
		subscribe = Subscribe(TestFrameworkInfo())
	})

	It("First call must be subscribe call", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)

		err := sched.Push(&PingMessage{}, context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("First message must be subscribe"))
		close(done)
	})

	It("Should subscribe", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			<-cl.ReqIn
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: &client.TestEmptyResponse{Sid: "1"}}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		close(done)
	})

	It("Set stream id on requests following subscribe call", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: client.NewTestChanResponse("1")}
			msg := <-cl.ReqIn
			req, err := http.NewRequest("POST", "http://localhost", gbytes.NewBuffer())
			Expect(err).To(Succeed())
			client.RequestOpts(msg.Opts).Apply(req)
			Expect(req.Header.Get(client.MESOS_STREAM_ID_HEADER)).To(Equal("1"))
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: &client.TestEmptyResponse{}}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		Expect(sched.Push(Teardown(), ctx)).To(Succeed())
		close(done)
	})

	It("Return context cancelled after subscribe failure", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()
		err := errors.New("Failed")

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn
			cl.ReqOut <- &client.TestClientResponseOrError{Err: err}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Equal(err))
		Expect(sched.Push(Teardown(), ctx)).To(Equal(context.Canceled))
		close(done)
	})

	It("Return context cancelled after response read failure", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn

			respChan := client.NewTestChanResponse("1")
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Err: errors.New("Failed")}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		_, e := sched.Pull(ctx)
		Expect(e).To(Equal(context.Canceled))
		close(done)
	})

	It("Pull remaining messages after context cancelled", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		testMsg := TestHeartbeat()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn

			respChan := client.NewTestChanResponse("1")
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Msg: testMsg}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Msg: testMsg}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Err: errors.New("Failed")}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		msg, e := sched.Pull(ctx)
		Expect(e).To(Succeed())
		Expect(msg).To(Equal(testMsg))
		msg, e = sched.Pull(ctx)
		Expect(e).To(Succeed())
		Expect(msg).To(Equal(testMsg))

		_, e = sched.Pull(ctx)
		Expect(e).To(Equal(context.Canceled))

		close(done)
	})

	It("Stop reading from response and signal buffer full when internal buffer is full", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		wait := make(chan struct{})

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := client.NewTestChanResponse("1")
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Msg: TestHeartbeat()}
			// should not pull more
			select {
			case <-respChan.ReadIn:
				Fail("Should not read more from response")
			case <-time.After(50 * time.Millisecond):
				close(wait)
			}
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		<-wait

		Expect(sched.Push(&PingMessage{}, ctx)).To(Equal(ErrBufferFull))
		close(done)
	})

	It("Return context cancelled after close", func(done Done) {
		sched := New(client.NewTestChanClient())
		ctx := context.Background()

		Expect(sched.Close()).To(Succeed())

		Expect(sched.Push(subscribe, ctx)).To(Equal(context.Canceled))
		_, e := sched.Pull(ctx)
		Expect(e).To(Equal(context.Canceled))
		close(done)
	})

	It("Close response after scheduler close", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		wait := make(chan struct{})

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := client.NewTestChanResponse("1")
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &client.TestMessageOrError{Msg: TestHeartbeat()}
			// should not pull more
			<-respChan.CloseIn
			respChan.CloseOut <- nil
			close(wait)
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		Expect(sched.Close()).To(Succeed())
		<-wait
		close(done)
	})

	It("Return on pull request context cancelled", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := client.NewTestChanResponse("1")
			cl.ReqOut <- &client.TestClientResponseOrError{Resp: respChan}
			// simulate reader blocked to allow testing of context cancelled
			<-respChan.ReadIn
		}()

		Expect(sched.Push(subscribe, ctx)).To(Succeed())
		c, cancel := context.WithCancel(ctx)
		cancel()
		_, e := sched.Pull(c)
		Expect(e).To(Equal(context.Canceled))

		close(done)
	})

	It("Return on push request context cancelled", func(done Done) {
		cl := client.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))

		go func() {
			defer GinkgoRecover()
			req := <-cl.ReqIn
			<-req.Ctx.Done()
			cl.ReqOut <- &client.TestClientResponseOrError{Err: req.Ctx.Err()}
		}()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		Expect(sched.Push(subscribe, ctx)).To(Equal(context.Canceled))
		close(done)
	})
})
