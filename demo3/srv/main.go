package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
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
	defer fmt.Println("关闭")
	fmt.Println("新连接：", conn.RemoteAddr())

	result := bytes.NewBuffer(nil)
	var buf [1024]byte
	for {
		n, err := conn.Read(buf[:])
		result.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				fmt.Println("read err:", err)
				break
			}
		} else {
			for {
				if result.Len() < 0 || result.Len() < 4 {
					fmt.Println("not enough 1111111111")
					break
				}
				var msgSize int32
				// message size
				err := binary.Read(bytes.NewReader(result.Bytes()), binary.BigEndian, &msgSize)
				if err != nil {
					panic(err)
				}
				//  4 字节的数据长度+具体数据
				lenBuf := result.Len()
				if int32(lenBuf) < msgSize+4 {
					fmt.Println("not enough-------------", string(result.Bytes()))
					break
				}
				// message binary data
				buf := make([]byte, msgSize+4)
				_, err = io.ReadFull(result, buf)
				if err != nil {
					fmt.Println(lenBuf)
					panic(err)
				}
				fmt.Printf("len %d recv: %s \n", len(buf), string(buf))
			}
		}
	}
}
