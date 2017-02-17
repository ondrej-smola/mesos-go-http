package client_test

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http"
	. "github.com/ondrej-smola/mesos-go-http/client"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"github.com/ondrej-smola/mesos-go-http/scheduler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"testing"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = Describe("Client", func() {

	endpoint := "http://myhost/test"
	msg := scheduler.Subscribe(mesos.FrameworkInfo{User: "test"})

	It("Send POST", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			Expect(r.Method).To(Equal("POST"))
			return NewProtoBytesResponse(http.StatusOK, nil), nil
		}))

		_, err := cl.Do(msg, context.Background())
		Expect(err).To(Succeed())
	})

	It("Set Context-Type from codec", func() {
		cl := New(endpoint, WithCodec(codec.ProtobufCodec), WithDoFunc(func(r *http.Request) (*http.Response, error) {
			Expect(r.Header.Get("Content-Type")).To(Equal(codec.ProtobufCodec.EncoderContentType))
			return NewProtoBytesResponse(http.StatusOK, nil), nil
		}))
		_, err := cl.Do(msg, context.Background())
		Expect(err).To(Succeed())
	})

	It("Set Content-Length", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			Expect(r.ContentLength).To(BeEquivalentTo(msg.Size()))
			return NewProtoBytesResponse(http.StatusOK, nil), nil
		}))
		_, err := cl.Do(msg, context.Background())
		Expect(err).To(Succeed())
	})

	It("Add request header", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			Expect(r.Header.Get("TEST")).To(Equal("TEST"))
			return NewProtoBytesResponse(http.StatusOK, nil), nil
		}))
		_, err := cl.Do(msg, context.Background(), mesos.WithHeader("TEST", "TEST"))
		Expect(err).To(Succeed())
	})

	It("Map errors - 503", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			return NewProtoBytesResponse(http.StatusServiceUnavailable, nil), nil
		}))
		_, err := cl.Do(msg, context.Background())
		Expect(err).To(Equal(UnavailableError))
	})

	It("Framing - single", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			body, err := proto.Marshal(msg)
			Expect(err).To(Succeed())
			return NewProtoBytesResponse(http.StatusOK, body), nil
		}), WithSingleFraming())

		resp, err := cl.Do(msg, context.Background())
		Expect(err).To(Succeed())

		m := &scheduler.Call{}
		Expect(resp.Read(m)).To(Succeed())
		Expect(resp.Read(m)).To(Equal(io.EOF))
		Expect(m).To(Equal(msg))

		Expect(resp.Close()).To(Succeed())
	})

	It("Framing - recordio", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			body, err := proto.Marshal(msg)
			Expect(err).To(Succeed())
			res := append([]byte(fmt.Sprintf("%v\n", len(body))), body...)
			return NewProtoBytesResponse(http.StatusOK, res), nil
		}), WithRecordIOFraming())

		resp, err := cl.Do(msg, context.Background())
		Expect(err).To(Succeed())

		m := &scheduler.Call{}
		Expect(resp.Read(m)).To(Succeed())
		Expect(m).To(Equal(msg))
		Expect(resp.Read(m)).To(Equal(io.EOF))

		Expect(resp.Close()).To(Succeed())
	})

	It("Handle redirect", func() {
		cl := New(endpoint, WithDoFunc(func(r *http.Request) (*http.Response, error) {
			resp := NewProtoBytesResponse(http.StatusTemporaryRedirect, nil)
			resp.Header.Set("Location", "http://localhost:5050/254")
			return resp, nil
		}))
		_, err := cl.Do(msg, context.Background())

		r, ok := err.(RedirectError)
		Expect(ok).To(BeTrue())
		Expect(r.LeaderHostPort).To(Equal("localhost:5050"))
	})
})
