package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func Marshal(c interface{}) ([]byte, error) {
	switch c := c.(type) {
	case CalcRequest:
		b := make([]byte, 24)
		binary.LittleEndian.PutUint64(b[:8], uint64(c.Op))
		binary.LittleEndian.PutUint64(b[8:16], uint64(c.A))
		binary.LittleEndian.PutUint64(b[16:], uint64(c.B))
		return b, nil
	case CalcResponse:
		b := make([]byte, 16)
		binary.LittleEndian.PutUint64(b[:8], uint64(c.Err))
		binary.LittleEndian.PutUint64(b[8:], uint64(c.Result))
		return b, nil
	}
	return nil, errors.New("no method")
}

func Unmarshal(b []byte, x interface{}) error {
	switch x := x.(type) {
	case *CalcRequest:
		if len(b) != CalcRequestLen {
			return errors.New("not match")
		}
		x.Op = Operator(binary.LittleEndian.Uint64(b[:8]))
		x.A = int64(binary.LittleEndian.Uint64(b[8:16]))
		x.B = int64(binary.LittleEndian.Uint64(b[16:]))
		return nil
	case *CalcResponse:
		if len(b) != CalcResponseLen {
			return errors.New("no match")
		}
		x.Err = CalcError(binary.LittleEndian.Uint64(b[:8]))
		x.Result = int64(binary.LittleEndian.Uint64(b[8:]))
		return nil
	default:
		return fmt.Errorf("%#v", x)
	}
}
