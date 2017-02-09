package master

import (
	"bytes"
	"github.com/gogo/protobuf/jsonpb"
)

var marshaller = jsonpb.Marshaler{EmitDefaults: true}

func (e *Call) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := marshaller.Marshal(buf, e)
	return buf.Bytes(), err
}

func (e *Event) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := marshaller.Marshal(buf, e)
	return buf.Bytes(), err
}

func (e *Call) UnmarshalJSON(in []byte) error {
	return jsonpb.Unmarshal(bytes.NewBuffer(in), e)
}

func (e *Event) UnmarshalJSON(in []byte) error {
	return jsonpb.Unmarshal(bytes.NewBuffer(in), e)
}

func (e *Response) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := marshaller.Marshal(buf, e)
	return buf.Bytes(), err
}

func (e *Response) UnmarshalJSON(in []byte) error {
	return jsonpb.Unmarshal(bytes.NewBuffer(in), e)
}
