package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
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

// 不好 丢数据呀~，可能是我用错了
func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client：", conn.RemoteAddr())

	result := bytes.NewBuffer(nil)
	var buf [1024]byte
	idx := 0
	start := time.Now()
	for {
		n, err := conn.Read(buf[:])
		result.Write(buf[:n])
		if err != nil {
			// closed
			if err == io.EOF {
				break
			} else {
				fmt.Println("read err:", err)
				break
			}
		} else {
			scanner := bufio.NewScanner(result)
			scanner.Split(bufSplit)
			for scanner.Scan() {
				idx++
				msgBuf := scanner.Bytes()
				fmt.Printf("len %d recv: %s count: %d \n", len(msgBuf), string(msgBuf), idx)
			}
		}
	}

	println(time.Since(start).Seconds())
}

func bufSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && len(data) > 4 {
		var msgSize int32
		// 读出 数据包中 实际数据 的长度(大小为 0 ~ 2^16)
		if err := binary.Read(bytes.NewReader(data), binary.BigEndian, &msgSize); err != nil {
			return 0, nil, err
		}
		if msgSize < 0 {
			return
		}
		totalMsgSize := int(msgSize) + 4
		if totalMsgSize <= len(data) {
			return totalMsgSize, data[:totalMsgSize], nil
		}
	}
	return
}
