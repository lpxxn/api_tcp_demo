package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8899")
	if err != nil {
		panic(err)
	}
	fmt.Println("listen to 8899")
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		} else {
			go handleConn(conn)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client：", conn.RemoteAddr())

	//result := bytes.NewBuffer(nil)
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[:])
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
			fmt.Printf("len %d recv: %s \n", n, string(buf[0:n]))
		}
		//result.Reset()
	}
}
