package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
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
	listen, _ := net.Listen("tcp", ":8080")
	newServer := grpc.NewServer()
	pd.RegisterHelloServerServer(newServer, new(Server))
	err := newServer.Serve(listen)
	if err != nil {
		log.Print("err:", err)
	}
}
