// simple RPC demo function which run a calculation app
// exec: go run client.go -op 1 -a 10 -b 20
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/zaynjarvis/fyp/rpc/protocol"
)

var (
	op = flag.Int64("op", 0, "operation for calculation")
	a  = flag.Int64("a", 0, "first operand")
	b  = flag.Int64("b", 0, "second operand")
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", "localhost:9700")
	if err != nil {
		panic(err)
	}

	res, err := Call(conn, protocol.CalculatorMethod, &protocol.CalcRequest{Op: protocol.Operator(*op), A: *a, B: *b})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res)
	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func Call(conn net.Conn, method protocol.Method, req interface{}) (interface{}, error) {
	errCh := make(chan error)
	resCh := make(chan interface{})
	go func() {
		res, err := protocol.ReadResponse(conn)
		if err != nil {
			errCh <- err
			return
		}
		resCh <- res
	}()
	data, err := protocol.MethodMap[method].Marshal(req)
	if err != nil {
		return nil, err
	}
	written, err := conn.Write(data)
	if err != nil {
		return nil, err
	}
	log.Printf("%v bytes sent", written)
	select {
	case err := <-errCh:
		return nil, err
	case res := <-resCh:
		return res, nil
	}
}
