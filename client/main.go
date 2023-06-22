package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pd "grpc_helloworld/client/proto"
	"io"
	"log"
)

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
func main() {
	user := &Authentication{
		User:     "admin",
		Password: "123456",
	}
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(user))

	if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()
	client := pd.NewHelloServerClient(conn)

	stream, err := client.SayStreamServer(context.Background(), &pd.HelloReq{
		Name: "111",
		Age:  0,
	})
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端数据接收完成")
				err := stream.CloseSend()
				if err != nil {
					log.Fatal(err)
				}
				break
			}
			log.Fatal(err)
		}
		fmt.Println("客户端收到的流", recv.Say)
	}
}
