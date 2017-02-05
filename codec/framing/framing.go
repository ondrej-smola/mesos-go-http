package framing

import (
	"io"
	"net/http"
)

type (
	Reader interface {
		ReadFrame(buf []byte) (endOfFrame bool, n int, err error)
	}

	ReaderFunc func(buf []byte) (endOfFrame bool, n int, err error)

	// Map io.Reader to framing.Reader
	Provider func(r io.Reader) Reader

	// Select provider based on response properties
	// Should not modify response
	Selector func(r *http.Response) Provider
)

// Functional adaptation of Reader interface
func (f ReaderFunc) ReadFrame(buf []byte) (bool, int, error) {
	return f(buf)
}
