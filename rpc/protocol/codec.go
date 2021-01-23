package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func ReadRequest(data io.Reader) (interface{}, error) {
	var buff bytes.Buffer
	_, _ = io.CopyN(&buff, data, HeaderLen)
	h := Header{}
	if err := UnmarshalHeader(buff.Bytes(), &h); err != nil {
		return 0, err
	}
	buff.Reset()
	_, _ = io.CopyN(&buff, data, h.Len)
	x := MethodMap[h.Method].CreateRequest()
	if err := MethodMap[h.Method].Unmarshal(buff.Bytes(), x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadResponse(data io.Reader) (interface{}, error) {
	var buff bytes.Buffer
	_, _ = io.CopyN(&buff, data, HeaderLen)
	h := Header{}
	if err := UnmarshalHeader(buff.Bytes(), &h); err != nil {
		return 0, err
	}
	buff.Reset()
	_, _ = io.CopyN(&buff, data, h.Len)
	x := MethodMap[h.Method].CreateResponse()
	if err := MethodMap[h.Method].Unmarshal(buff.Bytes(), x); err != nil {
		return 0, err
	}
	return x, nil
}

func MarshalHeader(h Header, b []byte) {
	binary.LittleEndian.PutUint64(b[:8], uint64(h.Method))
	binary.LittleEndian.PutUint64(b[8:], uint64(h.Len))
}

func UnmarshalHeader(b []byte, h *Header) error {
	if len(b) != HeaderLen {
		return errors.New("not header")
	}
	h.Method = Method(binary.LittleEndian.Uint64(b[:8]))
	h.Len = int64(binary.LittleEndian.Uint64(b[8:]))
	return nil
}

type Calculator struct {
	RequestLen  int64
	ResponseLen int64
}

func (Calculator) Marshal(c interface{}) ([]byte, error) {
	switch c := c.(type) {
	case *CalcRequest:
		t := make([]byte, HeaderLen+CalcRequestLen)
		MarshalHeader(Header{Method: CalculatorMethod, Len: CalcRequestLen}, t[:HeaderLen])
		b := t[HeaderLen:]
		binary.LittleEndian.PutUint64(b[:8], uint64(c.Op))
		binary.LittleEndian.PutUint64(b[8:16], uint64(c.A))
		binary.LittleEndian.PutUint64(b[16:], uint64(c.B))
		return t, nil
	case *CalcResponse:
		t := make([]byte, HeaderLen+CalcResponseLen)
		MarshalHeader(Header{Method: CalculatorMethod, Len: CalcResponseLen}, t[:HeaderLen])
		b := t[HeaderLen:]
		binary.LittleEndian.PutUint64(b[:8], uint64(c.Err))
		binary.LittleEndian.PutUint64(b[8:], uint64(c.Result))
		return t, nil
	}
	return nil, errors.New("no method")
}

func (Calculator) Unmarshal(b []byte, x interface{}) error {
	switch x := x.(type) {
	case *CalcRequest:
		if len(b) != CalcRequestLen {
			return errors.New("not calcReq")
		}
		x.Op = Operator(binary.LittleEndian.Uint64(b[:8]))
		x.A = int64(binary.LittleEndian.Uint64(b[8:16]))
		x.B = int64(binary.LittleEndian.Uint64(b[16:]))
		return nil
	case *CalcResponse:
		if len(b) != CalcResponseLen {
			return errors.New("not calcRes")
		}
		x.Err = CalcError(binary.LittleEndian.Uint64(b[:8]))
		x.Result = int64(binary.LittleEndian.Uint64(b[8:]))
		return nil
	default:
		return fmt.Errorf("%#v", x)
	}
}

func (Calculator) CreateRequest() interface{} {
	return &CalcRequest{}
}

func (Calculator) CreateResponse() interface{} {
	return &CalcResponse{}
}
