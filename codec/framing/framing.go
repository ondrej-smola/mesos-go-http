package framing

import (
	"io"
)

type (
	Reader interface {
		ReadFrame(buf []byte) (endOfFrame bool, n int, err error)
	}

	ReaderFunc func(buf []byte) (endOfFrame bool, n int, err error)

	// Map io.Reader to framing.Reader
	Provider func(r io.Reader) Reader
)

// Functional adaptation of Reader interface
func (f ReaderFunc) ReadFrame(buf []byte) (bool, int, error) {
	return f(buf)
}
