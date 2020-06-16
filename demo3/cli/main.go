package main

import (
	"fmt"
	"net"
	"time"

	"github.com/api_tcp_demo/demo3"
)

func main() {
	data := []byte("[这里才是一个完整的数据包]")
	conn, err := net.DialTimeout("tcp", "localhost:4044", time.Second*30)
	if err != nil {
		fmt.Printf("connect failed, err : %v\n", err.Error())
		return
	}
	for i := 0; i < 20000; i++ {
		go func() {
			if _, err = demo3.ConnRW.WriteTo(conn, data); err != nil {
				fmt.Printf("write failed , err : %v\n", err)
				panic(err)
		}
		}()
	}
	time.Sleep(time.Second * 5)
}
