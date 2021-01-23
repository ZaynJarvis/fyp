package main

import (
	"log"
	"net"

	"github.com/zaynjarvis/fyp/rpc/protocol"
)

func main() {
	listen, err := net.Listen("tcp", ":9700")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("%v\n", err)
		}
		Serve(conn)
	}
}

func Serve(conn net.Conn) {
	req, err := protocol.ReadRequest(conn)
	if err != nil {
		panic(err)
	}
	var data []byte
	switch req := req.(type) {
	case *protocol.CalcRequest:
		res := Exec(req)
		data, err = protocol.MethodMap[protocol.CalculatorMethod].Marshal(res)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	write, err := conn.Write(data)
	if err != nil {
		panic(err)
	}
	log.Printf("%v bytes sent", write)
}

func Exec(req *protocol.CalcRequest) *protocol.CalcResponse {
	switch req.Op {
	case protocol.Add:
		return &protocol.CalcResponse{Err: protocol.Nil, Result: req.A + req.B}
	case protocol.Sub:
		return &protocol.CalcResponse{Err: protocol.Nil, Result: req.A - req.B}
	case protocol.Mul:
		return &protocol.CalcResponse{Err: protocol.Nil, Result: req.A * req.B}
	case protocol.Div:
		return &protocol.CalcResponse{Err: protocol.Nil, Result: req.A / req.B}
	default:
		return &protocol.CalcResponse{Err: protocol.NoMethod, Result: 0}
	}
}
