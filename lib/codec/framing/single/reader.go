package single

import (
	"io"

	"github.com/ondrej-smola/mesos-go-http/lib/codec/framing"
)

type Reader struct {
	r io.Reader
}

func New(r io.Reader) *Reader {
	return &Reader{r: r}
}

// Single frame reader is used for non chunked response.
// Emits whole response body as single frame
func NewProvider() framing.Provider {
	return func(r io.Reader) framing.Reader {
		return New(r)
	}
}

func (rr *Reader) ReadFrame(p []byte) (endOfFrame bool, n int, err error) {
	n, err = rr.r.Read(p)

	if err == io.EOF {
		endOfFrame = true
	}

	return
}

func (rr *Reader) Read(p []byte) (int, error) {
	return rr.r.Read(p)
}
