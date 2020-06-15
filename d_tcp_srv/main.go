package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/api_tcp_demo/common"
	"github.com/api_tcp_demo/protos"
	"google.golang.org/grpc"
)

var (
	tcpPort string
	rpcPort string
	TcpSrv  *tcpServer
)

func init() {
	flag.StringVar(&tcpPort, "tcp_port", common.DefaultTcpPort, "tcp server port")
	flag.StringVar(&rpcPort, "rpc_port", common.DefaultRpcPort, "rpc server port")
	flag.Parse()
}

func main() {
	go StartRpc()
	tcpListener, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		panic(err)
	}
	TcpSrv = &tcpServer{MsgChan: make(chan *ReplayMsg, 500)}
	if err := TcpSrv.TCPServer(tcpListener); err != nil {
		panic(err)
	}
}

type tcpServer struct {
	conns   sync.Map
	MsgChan chan *ReplayMsg
}

func (s *tcpServer) TCPServer(listener net.Listener) error {
	fmt.Printf("tcp run at  %s \n", listener.Addr())
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				fmt.Printf("temporary Accept() failure - %s\n", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return fmt.Errorf("listener.Accept() error - %s", err)
			}
			break
		}

		go func() {
			s.HandleTcpClientConn(clientConn)
		}()
	}
	fmt.Println("TCP closing ...")
	return nil
}

func (s *tcpServer) HandleTcpClientConn(clientConn net.Conn) {
	client := &ConnClient{conn: clientConn}
	fmt.Println("got a connection")
	s.conns.Store(clientConn.RemoteAddr(), clientConn)
	go s.MessageProcess(client)
	for {
		if atomic.LoadInt32(&client.CloseFlag) == 1 {
			break
		}
		msg, err := common.ReadResponse(client.conn)
		if err != nil {
			if os.IsTimeout(err) {
				continue
			}
			if err == io.EOF && atomic.LoadInt32(&client.CloseFlag) == 1 {
				break
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("IO error - %s", err)
			}
			break
		}
		if client.CurrentReplayMsg == nil {
			continue
		}
		select {
		case <-client.CurrentReplayMsg.Ctx.Done():
			fmt.Println("client replay msg ctx done")
		case client.CurrentReplayMsg.RepMsg <- msg:
		}
		client.CurrentReplayMsg = nil
	}
	client.Close()
	s.conns.Delete(clientConn.RemoteAddr())
}

func (s *tcpServer) MessageProcess(client *ConnClient) {
	for {
		if client.MsgChan == nil && client.CurrentReplayMsg == nil {
			client.MsgChan = s.MsgChan
		} else {
			client.MsgChan = nil
		}
		select {
		// 收到一条信息
		case msg := <-client.MsgChan:
			client.CurrentReplayMsg = msg
			if _, err := common.WriteTo(client, msg.SendMsg); err != nil {
				break
			}
		default:

		}
	}
}

type ReplayMsg struct {
	Ctx     context.Context
	SendMsg *protos.Msg
	RepMsg  chan *protos.Msg
}

func NewReplayMsg(ctx context.Context, sendMsg *protos.Msg) *ReplayMsg {
	return &ReplayMsg{Ctx: ctx, SendMsg: sendMsg, RepMsg: make(chan *protos.Msg)}
}

type ConnClient struct {
	conn net.Conn
	// 处理服务器信息
	MsgChan chan *ReplayMsg
	// 服务器返回给api的信息
	CurrentReplayMsg *ReplayMsg
	CloseFlag        int32
}

func (c *ConnClient) Read(p []byte) (int, error) {
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	return c.conn.Read(p)
}

// Write performs a deadlined write on the underlying TCP connection
func (c *ConnClient) Write(p []byte) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	return c.conn.Write(p)
}

func (c *ConnClient) Close() error {
	atomic.StoreInt32(&c.CloseFlag, 1)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type RpcHandler struct {
}

func (r *RpcHandler) ProgressMsg(ctx context.Context, msg *protos.Msg) (*protos.Msg, error) {
	tCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	rMsg := NewReplayMsg(tCtx, msg)
	TcpSrv.MsgChan <- rMsg
	revMsg := &protos.Msg{}
	select {
	case <-tCtx.Done():
		revMsg.Type = "-1"
		revMsg.Data = "time out"
	case revMsg = <-rMsg.RepMsg:
	}
	return revMsg, nil
}

func StartRpc() {
	address := fmt.Sprintf(":%s", rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("start rpc provider error:%v", err))
	}

	fmt.Printf("rpc running at port [%s] \n", rpcPort)
	rpcSrv := grpc.NewServer()
	protos.RegisterHandlerMsgServer(rpcSrv, new(RpcHandler))
	if err := rpcSrv.Serve(listener); err != nil {
		panic(err)
	}
}
