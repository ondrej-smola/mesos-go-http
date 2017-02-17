package leader_test

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/client"
	. "github.com/ondrej-smola/mesos-go-http/client/leader"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leader Suite")
}

var _ = Describe("Leader", func() {

	endpoints := []string{"http://master1/test", "http://master2/test", "http://master3/test"}
	msg := scheduler.Subscribe(mesos.FrameworkInfo{User: "test"})

	It("Follow redirects", func(done Done) {
		tProv := client.NewTestClientProvider()

		cl := New(endpoints, WithClientProvider(tProv))
		go func() {
			defer GinkgoRecover()
			Expect(<-tProv.NewIn).To(Equal(endpoints[0]))
			tProv.NewOut <- mesos.DoFunc(func(proto.Message, context.Context, ...mesos.RequestOpt) (resp mesos.Response, err error) {
				return nil, client.RedirectError{LeaderHostPort: "host2:5050"}
			})

			Expect(<-tProv.NewIn).To(Equal("http://host2:5050/test"))
			tProv.NewOut <- mesos.DoFunc(func(proto.Message, context.Context, ...mesos.RequestOpt) (resp mesos.Response, err error) {
				return &mesos.TestEmptyResponse{}, nil
			})
		}()

		_, err := cl.Do(msg, context.Background())

		Expect(err).To(Succeed())
		close(done)
	})

	It("Try all endpoints before failing", func(done Done) {
		tProv := client.NewTestClientProvider()

		cl := New(endpoints, WithClientProvider(tProv))
		wait := make(chan struct{})

		go func() {
			defer GinkgoRecover()

			for i := 0; i < len(endpoints); i++ {
				Expect(<-tProv.NewIn).To(Equal(endpoints[i]))
				tProv.NewOut <- mesos.DoFunc(func(proto.Message, context.Context, ...mesos.RequestOpt) (resp mesos.Response, err error) {
					return nil, client.NotFoundError
				})
			}
			close(wait)
		}()

		_, err := cl.Do(msg, context.Background())
		Expect(err).To(HaveOccurred())
		<-wait
		close(done)
	})

	It("Max redirects", func(done Done) {
		tProv := client.NewTestClientProvider()

		cl := New(endpoints, WithClientProvider(tProv), WithMaxRedirects(5))

		stop := make(chan struct{})
		go func() {
			defer GinkgoRecover()

			for {
				select {
				case <-stop:
					break
				case <-tProv.NewIn:
					tProv.NewOut <- mesos.DoFunc(func(proto.Message, context.Context, ...mesos.RequestOpt) (resp mesos.Response, err error) {
						return nil, client.RedirectError{LeaderHostPort: "host2:5050"}
					})
				}
			}
		}()

		_, err := cl.Do(msg, context.Background())
		Expect(err).To(HaveOccurred())
		close(stop)
		close(done)
	})

	It("Context cancel", func(done Done) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tProv := client.NewTestClientProvider()
		cl := New(endpoints, WithClientProvider(tProv), WithMaxRedirects(5))

		stop := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			<-tProv.NewIn
			tProv.NewOut <- mesos.DoFunc(func(proto.Message, context.Context, ...mesos.RequestOpt) (resp mesos.Response, err error) {
				cancel()
				return nil, client.NotFoundError
			})

			select {
			case <-stop:
				break
			case <-tProv.NewIn:
				Fail("Should not create new client when cancelled")
			}

		}()

		_, err := cl.Do(msg, ctx)
		Expect(err).To(HaveOccurred())
		close(stop)
		close(done)
	})
})
