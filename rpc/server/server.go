package main

import (
	"bytes"
	"io"
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
		var buff bytes.Buffer
		written, err := io.CopyN(&buff, conn, protocol.CalcRequestLen)
		if err != nil {
			panic(err)
		}
		log.Printf("%v bytes received", written)
		var req protocol.CalcRequest
		if err := protocol.Unmarshal(buff.Bytes(), &req); err != nil {
			panic(err)
		}
		res := Exec(req)
		data, err := protocol.Marshal(res)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		write, err := conn.Write(data)
		if err != nil {
			panic(err)
		}
		log.Printf("%v bytes sent", write)
	}
}

func Exec(req protocol.CalcRequest) protocol.CalcResponse {
	switch req.Op {
	case protocol.Add:
		return protocol.CalcResponse{Err: protocol.Nil, Result: req.A + req.B}
	case protocol.Sub:
		return protocol.CalcResponse{Err: protocol.Nil, Result: req.A - req.B}
	case protocol.Mul:
		return protocol.CalcResponse{Err: protocol.Nil, Result: req.A * req.B}
	case protocol.Div:
		return protocol.CalcResponse{Err: protocol.Nil, Result: req.A / req.B}
	default:
		return protocol.CalcResponse{Err: protocol.NoMethod, Result: 0}
	}
}
