// simple RPC demo function which run a calculation app
// exec: go run client.go -op 1 -a 10 -b 20
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"sync"

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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var buff bytes.Buffer
		written, err := io.CopyN(&buff, conn, protocol.CalcResponseLen)
		if err != nil {
			panic(err)
		}
		log.Printf("%v bytes received", written)
		var res protocol.CalcResponse
		if err := protocol.Unmarshal(buff.Bytes(), &res); err != nil {
			panic(err)
		}
		log.Printf("%#v", res)
	}()
	data, err := protocol.Marshal(protocol.CalcRequest{Op: protocol.Operator(*op), A: *a, B: *b})
	if err != nil {
		panic(err)
	}
	written, err := conn.Write(data)
	if err != nil {
		panic(err)
	}
	log.Printf("%v bytes sent", written)
	wg.Wait()
	if err := conn.(io.WriteCloser).Close(); err != nil {
		panic(err)
	}
}
