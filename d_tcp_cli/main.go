package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/api_tcp_demo/common"
)

var (
	tcpServer string
)

func init() {
	flag.StringVar(&tcpServer, "tcp_server", common.DefaultServer+":"+common.DefaultTcpPort, "tcp server")
}
func main() {
	conn, err := net.Dial("tcp", tcpServer)
	if err != nil {
		fmt.Println("please run tcp server first")
		panic(err)
	}
	client := &Conn{conn: conn.(*net.TCPConn), ID: time.Now().Unix()}
	client.Loop()
	fmt.Println("end...")
}

type Conn struct {
	ID        int64
	conn      *net.TCPConn
	CloseFlag int32
}

func (c *Conn) Loop() {
	for {
		if atomic.LoadInt32(&c.CloseFlag) == 1 {
			break
		}
		msg, err := common.ReadResponse(c.conn)
		if err != nil {
			if os.IsTimeout(err) {
				continue
			}
			if err == io.EOF && atomic.LoadInt32(&c.CloseFlag) == 1 {
				break
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("IO error - %s", err)
			}
			break
		}
		fmt.Printf("client: %d receive msg: %s \n", c.ID, msg.Data)
		msg.Data = fmt.Sprintf("client: %d receive msg: %s, currentTime: %d", c.ID, msg.Data, time.Now().Unix())
		if _, err := common.WriteTo(c, msg); err != nil {
			break
		}
	}
	c.Close()
}

func (c *Conn) Read(p []byte) (int, error) {
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	return c.conn.Read(p)
}

func (c *Conn) Write(p []byte) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	return c.conn.Write(p)
}

func (c *Conn) Close() error {
	atomic.StoreInt32(&c.CloseFlag, 1)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
