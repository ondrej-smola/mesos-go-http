package selectors

import (
	"github.com/ondrej-smola/mesos-go-http/codec/framing"
	"github.com/ondrej-smola/mesos-go-http/codec/framing/recordio"
	"github.com/ondrej-smola/mesos-go-http/codec/framing/single"
	"net/http"
)

// Selects framing.Provider based on whether Context-Length header is set
func ByContentLength(r *http.Response) framing.Provider {
	if r.Header.Get("Content-Length") == "" {
		// user recordio for chunked responses
		return recordio.NewProvider()
	} else {
		// use single frame reader for others
		return single.NewProvider()
	}
}
