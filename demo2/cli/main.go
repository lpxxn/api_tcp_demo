package main

import (
	"fmt"
	"net"
	"time"

	"github.com/api_tcp_demo/demo2"
)

func main() {
	data := []byte("~测试数据：一二三四五~")
	conn, err := net.Dial("tcp", ":8899")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 20000; i++ {
		go func() {
			if _, err = demo2.ConnRW.WriteTo(conn, data); err != nil {
				fmt.Printf("write failed , err : %v\n", err)
				panic(err)
		}
		}()
	}
	time.Sleep(time.Second * 4)
}
