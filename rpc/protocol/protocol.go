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

type Method int64

const (
	Unkonwn Method = iota
	CalculatorMethod
)

type Codec interface {
	Marshal(c interface{}) ([]byte, error)
	Unmarshal(b []byte, x interface{}) error
	CreateRequest() interface{}
	CreateResponse() interface{}
}

const (
	HeaderLen       = 16
	CalcRequestLen  = 24
	CalcResponseLen = 16
)

var MethodMap = map[Method]Codec{
	CalculatorMethod: Calculator{
		RequestLen:  CalcResponseLen,
		ResponseLen: CalcResponseLen,
	},
}

type Header struct {
	Method Method
	Len    int64
}

type CalcRequest struct {
	Op Operator
	A  int64
	B  int64
}

type CalcResponse struct {
	Err    CalcError
	Result int64
}
