package protocol

type Operator int64

const (
	Nop Operator = iota
	Add
	Sub
	Mul
	Div
)

type CalcError int64

const (
	Nil CalcError = iota
	NoMethod
)

const (
	CalcRequestLen  = 24
	CalcResponseLen = 16
)

type CalcRequest struct {
	Op Operator
	A  int64
	B  int64
}

type CalcResponse struct {
	Err    CalcError
	Result int64
}
