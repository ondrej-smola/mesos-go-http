package codec

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/codec/framing"
	"io"
)

type (
	NewEncoderFunc func(w io.Writer) *Encoder
	NewDecoderFunc func(r framing.Reader) *Decoder

	UnmarshalFunc func(src []byte, m proto.Message) error
	MarshalFunc   func(m proto.Message) ([]byte, error)

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
	marsh := jsonpb.Marshaler{EmitDefaults: true}

	return &Encoder{
		w: w,
		mf: func(m proto.Message) ([]byte, error) {
			b := &bytes.Buffer{}
			err := marsh.Marshal(b, m)
			return b.Bytes(), err
		},
	}
}

func NewProtobufEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:  w,
		mf: proto.Marshal,
	}
}

func NewProtobufDecoder(r framing.Reader) *Decoder {
	return &Decoder{
		r:   r,
		buf: make([]byte, DECODER_BUFFER_SIZE),
		uf:  proto.Unmarshal,
	}
}

func NewJsonDecoder(r framing.Reader) *Decoder {
	return &Decoder{
		r:   r,
		buf: make([]byte, DECODER_BUFFER_SIZE),
		uf: func(src []byte, m proto.Message) error {
			return jsonpb.Unmarshal(bytes.NewBuffer(src), m)
		},
	}
}

func (e *Encoder) Encode(m proto.Message) error {
	bs, err := e.mf(m)
	if err != nil {
		return err
	}
	_, err = e.w.Write(bs)
	return err
}

// Decode reads the next message from its input and stores it in the value pointed to by m.
func (d *Decoder) Decode(m proto.Message) error {
	var (
		buf     = d.buf
		readlen = 0
	)
	for {
		eof, nr, err := d.r.ReadFrame(buf)
		readlen += nr
		if readlen > MAX_SIZE_BYTES {
			return ErrSize
		}

		if eof && readlen > 0 {
			return d.uf(d.buf[:readlen], m)
		} else if err != nil {
			return err
		} else if len(buf) == nr {
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
