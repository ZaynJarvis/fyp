package protocol

import (
	"testing"
	"unsafe"
)

func TestStruct(t *testing.T) {
	c := CalcRequest{Op: Add, A: 10, B: 10}
	if size := unsafe.Sizeof(c); size != CalcRequestLen {
		t.Errorf("size of the Calc Obj incorrect, got %v want %v", size, CalcRequestLen)
	}
	r := CalcResponse{Err: NoMethod, Result: 0}
	if size := unsafe.Sizeof(r); size != CalcResponseLen {
		t.Errorf("size of the Calc Obj incorrect, got %v want %v", size, CalcResponseLen)
	}
}
