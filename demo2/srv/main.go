package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/api_tcp_demo/demo2"
)

func main() {
	l, err := net.Listen("tcp", ":4044")
	if err != nil {
		panic(err)
	}
	fmt.Println("listen to 4044")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("conn err:", err)
		} else {
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client: ", conn.RemoteAddr())

	//result := bytes.NewBuffer(nil)
	start := time.Now()
	idx := 0
	for {
		buf, err := demo2.ConnRW.ReadResponse(conn)
		//result.Write(buf[0:n])
		if err != nil {
			// closed
			if err == io.EOF {
				break
			} else {
				fmt.Println("read err:", err)
				break
			}
		} else {
			idx++
			fmt.Printf("len %d recv: %s idx: %d\n", len(buf), string(buf), idx)
		}
		//result.Reset()
	}
	println(time.Since(start).Seconds())
}
