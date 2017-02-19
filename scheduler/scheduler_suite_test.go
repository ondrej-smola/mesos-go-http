package scheduler_test

import (
	"context"
	"github.com/ondrej-smola/mesos-go-http"
	. "github.com/ondrej-smola/mesos-go-http/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pkg/errors"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scheduler Suite")
}

var _ = Describe("Scheduler", func() {

	It("First call must be subscribe call", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)

		err := sched.Push(&PingMessage{}, context.Background())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("First message must be subscribe"))
		close(done)
	})

	It("Should subscribe", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			<-cl.ReqIn
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: &mesos.TestEmptyResponse{Sid: "1"}}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		close(done)
	})

	It("Set stream id on requests following subscribe call", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: &mesos.TestEmptyResponse{Sid: "1"}}
			msg := <-cl.ReqIn
			req, err := http.NewRequest("POST", "http://localhost", gbytes.NewBuffer())
			Expect(err).To(Succeed())
			mesos.RequestOpts(msg.Opts).Apply(req)
			Expect(req.Header.Get(mesos.MESOS_STREAM_ID_HEADER)).To(Equal("1"))
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: &mesos.TestEmptyResponse{Sid: "1"}}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		Expect(sched.Push(Teardown(), ctx)).To(Succeed())
		close(done)
	})

	It("Return context cancelled after subscribe failure", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()
		err := errors.New("Failed")

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn
			cl.ReqOut <- &mesos.TestClientResponseOrError{Err: err}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Equal(err))
		Expect(sched.Push(Teardown(), ctx)).To(Equal(context.Canceled))
		close(done)
	})

	It("Return context cancelled after response read failure", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn

			respChan := mesos.NewTestChanResponse("1")
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Err: errors.New("Failed")}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		_, e := sched.Pull(ctx)
		Expect(e).To(Equal(context.Canceled))
		close(done)
	})

	It("Pull remaining messages after context cancelled", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl)
		ctx := context.Background()

		testMsg := TestHeartbeat()

		go func() {
			defer GinkgoRecover()
			<-cl.ReqIn

			respChan := mesos.NewTestChanResponse("1")
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Msg: testMsg}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Msg: testMsg}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Err: errors.New("Failed")}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
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
		cl := mesos.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		wait := make(chan struct{})

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := mesos.NewTestChanResponse("1")
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Msg: TestHeartbeat()}
			// should not pull more
			select {
			case <-respChan.ReadIn:
				Fail("Should not read more from response")
			case <-time.After(50 * time.Millisecond):
				close(wait)
			}
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		<-wait

		Expect(sched.Push(&PingMessage{}, ctx)).To(Equal(ErrBufferFull))
		close(done)
	})

	It("Return context cancelled after close", func(done Done) {
		sched := New(mesos.NewTestChanClient())
		ctx := context.Background()

		Expect(sched.Close()).To(Succeed())

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Equal(context.Canceled))
		_, e := sched.Pull(ctx)
		Expect(e).To(Equal(context.Canceled))
		close(done)
	})

	It("Close response after scheduler close", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		wait := make(chan struct{})

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := mesos.NewTestChanResponse("1")
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: respChan}
			<-respChan.ReadIn
			respChan.ReadOut <- &mesos.TestMessageOrError{Msg: TestHeartbeat()}
			// should not pull more
			<-respChan.CloseIn
			respChan.CloseOut <- nil
			close(wait)
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		Expect(sched.Close()).To(Succeed())
		<-wait
		close(done)
	})

	It("Return on pull request context cancelled", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))
		ctx := context.Background()

		go func() {
			defer GinkgoRecover()

			<-cl.ReqIn
			respChan := mesos.NewTestChanResponse("1")
			cl.ReqOut <- &mesos.TestClientResponseOrError{Resp: respChan}
			// simulate reader blocked to allow testing of context cancelled
			<-respChan.ReadIn
		}()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Succeed())
		c, cancel := context.WithCancel(ctx)
		cancel()
		_, e := sched.Pull(c)
		Expect(e).To(Equal(context.Canceled))

		close(done)
	})

	It("Return on push request context cancelled", func(done Done) {
		cl := mesos.NewTestChanClient()
		sched := New(cl, WithBufferSize(0))

		go func() {
			defer GinkgoRecover()
			req := <-cl.ReqIn
			<-req.Ctx.Done()
			cl.ReqOut <- &mesos.TestClientResponseOrError{Err: req.Ctx.Err()}
		}()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		Expect(sched.Push(Subscribe(mesos.FrameworkInfo{User: "test"}), ctx)).To(Equal(context.Canceled))
		close(done)
	})
})
