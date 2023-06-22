package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	pd "grpc_helloworld/server/proto"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	pd.UnimplementedHelloServerServer
}

func (s *Server) Say(ctx context.Context, q *pd.HelloReq) (*pd.HelloRep, error) {
	return &pd.HelloRep{
		Say: fmt.Sprintf("%s已经%d岁了", q.Name, q.Age),
	}, nil
}
func (s *Server) SayStream(stream pd.HelloServer_SayStreamServer) error {
	for {
		//源源不断的去接收客户端发来的信息
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Println("服务端接收到的流", recv.Name, recv.Age)
	}
	return stream.SendAndClose(&pd.HelloRep{Say: "结束"})
}
func (s *Server) SayStreamServer(req *pd.HelloReq, stream pd.HelloServer_SayStreamServerServer) error {
	count := 0
	fmt.Printf(req.Name, req.Age)
	for {
		rsp := &pd.HelloRep{Say: "server send"}
		err := stream.Send(rsp)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
		count++
		if count > 10 {
			return nil
		}
	}
}
func (s *Server) SayDoubleStream(stream pd.HelloServer_SayDoubleStreamServer) error {
	for {
		recv, err := stream.Recv()
		if err != nil {
			return nil
		}
		fmt.Println("服务端收到客户端的消息", recv.Name, recv.Age)
		time.Sleep(time.Second)
		rsp := &pd.HelloRep{Say: fmt.Sprintf("%s的年龄：%d", recv.Name, recv.Age)}
		err = stream.Send(rsp)
		if err != nil {
			return nil
		}
	}
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
