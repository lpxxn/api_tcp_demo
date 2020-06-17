package main

import (
	"fmt"
	"net"
)

func main() {
	data := []byte("~测试数据：一二三四五~")
	conn, err := net.Dial("tcp", ":8899")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 2000; i++ {
		if _, err = conn.Write(data); err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}
	}
}
