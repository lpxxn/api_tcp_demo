package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/api_tcp_demo/common"
	"github.com/api_tcp_demo/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	httpPort  string
	rpcServer string

	rpcClient *grpc.ClientConn
)

func init() {
	flag.StringVar(&httpPort, "tcp_port", common.DefaultHttpPort, "tcp server port")
	flag.StringVar(&rpcServer, "rpc_server", common.DefaultServer+":"+common.DefaultRpcPort, "rpc server")

	flag.Parse()
}

func main() {
	if err := ConnRpc(); err != nil {
		panic(err)
	}
	r := gin.New()
	r.POST("progressData", ProgressData)
	fmt.Printf("http server running at http_port [%s] \n", httpPort)
	fmt.Println(r.Run(":" + httpPort))
}

func ProgressData(c *gin.Context) {
	param := &protos.Msg{}
	if err := c.ShouldBindJSON(param); err != nil || param.Type == "" {
		c.JSON(http.StatusOK, "invalid parameter")
		return
	}

	client := protos.NewHandlerMsgClient(rpcClient)
	resp, err := client.ProgressMsg(context.Background(), param)
	if err != nil {
		c.JSON(http.StatusOK, err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

func ConnRpc() (err error) {
	ctx := context.Background()
	rpcClient, err = grpc.DialContext(ctx, rpcServer, grpc.WithInsecure())
	return err
}
