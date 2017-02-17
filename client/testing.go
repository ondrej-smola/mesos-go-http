package client

import (
	"bytes"
	"github.com/ondrej-smola/mesos-go-http"
	"github.com/ondrej-smola/mesos-go-http/codec"
	"io"
	"net/http"
	"strconv"
)

func NewRespBodyFromBytes(body []byte) io.ReadCloser {
	return &dummyReadCloser{bytes.NewReader(body)}
}

func NewProtoBytesResponse(status int, body []byte) *http.Response {
	return &http.Response{
		ContentLength: int64(len(body)),
		Status:        strconv.Itoa(status),
		StatusCode:    status,
		Body:          NewRespBodyFromBytes(body),
		Header:        http.Header{"Content-Type": []string{codec.ProtobufCodec.DecoderContentType}},
	}
}

type dummyReadCloser struct {
	body io.ReadSeeker
}

func (d *dummyReadCloser) Read(p []byte) (int, error) {
	return d.body.Read(p)
}

func (d *dummyReadCloser) Close() error {
	return nil
}

type TestChanClientProvider struct {
	NewIn  chan string
	NewOut chan mesos.Client
}

func NewTestClientProvider() *TestChanClientProvider {
	return &TestChanClientProvider{
		NewIn:  make(chan string),
		NewOut: make(chan mesos.Client),
	}
}

func (t *TestChanClientProvider) New(endpoint string, opt ...Opt) mesos.Client {
	t.NewIn <- endpoint
	return <-t.NewOut
}
