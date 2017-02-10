package codec

import (
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/codec/framing"
	"io"
)

type (
	Message interface {
		json.Marshaler
		json.Unmarshaler
		proto.Marshaler
		proto.Unmarshaler
	}

	NewEncoderFunc func(w io.Writer) *Encoder
	NewDecoderFunc func(r framing.Reader) *Decoder

	UnmarshalFunc func(src []byte, m Message) error
	MarshalFunc   func(m Message) ([]byte, error)

	Codec struct {
		Name               string
		DecoderContentType string
		EncoderContentType string
		NewEncoder         NewEncoderFunc
		NewDecoder         NewDecoderFunc
	}

	Decoder struct {
		r   framing.Reader
		buf []byte
		uf  UnmarshalFunc
	}

	Encoder struct {
		w  io.Writer
		mf MarshalFunc
	}
)

const (
	MAX_SIZE_MB         = 4
	MAX_SIZE_BYTES      = MAX_SIZE_MB * 1024 * 1024
	DECODER_BUFFER_SIZE = 4 * 1024
)

var (
	// ErrSize is returned by Decode calls when a message would exceed the maximum allowed size.
	ErrSize = fmt.Errorf("proto: message exceeds %dMB", MAX_SIZE_MB)

	// BROKEN till https://issues.apache.org/jira/browse/MESOS-5995 is fixed
	// problem with uint64 <-> string conversion
	JsonCodec = &Codec{
		Name:               "json",
		DecoderContentType: "application/json",
		EncoderContentType: "application/json",
		NewDecoder: func(r framing.Reader) *Decoder {
			return NewJsonDecoder(r)
		},
		NewEncoder: func(w io.Writer) *Encoder {
			return NewJsonEncoder(w)
		},
	}

	ProtobufCodec = &Codec{
		Name:               "protobuf",
		DecoderContentType: "application/x-protobuf",
		EncoderContentType: "application/x-protobuf",
		NewDecoder: func(r framing.Reader) *Decoder {
			return NewProtobufDecoder(r)
		},
		NewEncoder: func(w io.Writer) *Encoder {
			return NewProtobufEncoder(w)
		},
	}
)

func NewJsonEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
		mf: func(m Message) ([]byte, error) {
			return m.MarshalJSON()
		},
	}
}

func NewProtobufEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
		mf: func(m Message) ([]byte, error) {
			return m.Marshal()
		},
	}
}

func NewProtobufDecoder(r framing.Reader) *Decoder {
	return &Decoder{
		r:   r,
		buf: make([]byte, DECODER_BUFFER_SIZE),
		uf: func(src []byte, m Message) error {
			return m.Unmarshal(src)
		},
	}
}

func NewJsonDecoder(r framing.Reader) *Decoder {
	return &Decoder{
		r:   r,
		buf: make([]byte, DECODER_BUFFER_SIZE),
		uf: func(src []byte, m Message) error {
			return m.UnmarshalJSON(src)
		},
	}
}

func (e *Encoder) Encode(m Message) error {
	bs, err := e.mf(m)
	if err != nil {
		return err
	}
	_, err = e.w.Write(bs)
	return err
}

// Decode reads the next message from its input and stores it in the value pointed to by m.
func (d *Decoder) Decode(m Message) error {
	var (
		buf     = d.buf
		readlen = 0
	)
	for {
		eof, nr, err := d.r.ReadFrame(buf)
		if err != nil {
			return err
		}

		readlen += nr
		if readlen > MAX_SIZE_BYTES {
			return ErrSize
		}

		if eof {
			return d.uf(d.buf[:readlen], m)
		}

		if len(buf) == nr {
			// readlen and len(d.buf) are the same here
			newbuf := make([]byte, readlen+4096)
			copy(newbuf, d.buf)
			d.buf = newbuf
			buf = d.buf[readlen:]
		} else {
			buf = buf[nr:]
		}
	}
}
