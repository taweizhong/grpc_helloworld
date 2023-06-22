package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pd "grpc_helloworld/server/proto"
	"log"
	"net"
)

type Server struct {
	pd.UnimplementedHelloServerServer
}

func (s *Server) Say(ctx context.Context, q *pd.HelloReq) (*pd.HelloRep, error) {
	return &pd.HelloRep{
		Say: fmt.Sprintf("%s已经%d岁了", q.Name, q.Age),
	}, nil
}
func main() {
	// 读取证书
	file, _ := credentials.NewServerTLSFromFile("./server.pem", "./server.key")

	listen, _ := net.Listen("tcp", ":8080")
	// 创建服务并添加证书
	newServer := grpc.NewServer(grpc.Creds(file))
	pd.RegisterHelloServerServer(newServer, new(Server))
	err := newServer.Serve(listen)
	if err != nil {
		log.Print("err:", err)
	}
}
