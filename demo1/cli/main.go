package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	data := []byte("~测试数据：一二三四五~")
	conn, err := net.DialTimeout("tcp", "localhost:4044", time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 2000; i++ {
		if _, err = conn.Write(data); err != nil {
			fmt.Printf("write failed , err : %v\n", err)
			break
		}
	}
}
