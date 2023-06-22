package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	var authInterceptor grpc.UnaryServerInterceptor
	authInterceptor = func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		//拦截普通方法请求，验证 Token
		err = Auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}

	listen, _ := net.Listen("tcp", ":8080")
	// 创建服务并添加证书
	newServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))

	pd.RegisterHelloServerServer(newServer, new(Server))
	err := newServer.Serve(listen)
	if err != nil {
		log.Print("err:", err)
	}
}
func Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var user string
	var password string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if user != "admin" || password != "123456" {
		return status.Errorf(codes.Unauthenticated, "token不合法")
	}
	return nil
}
