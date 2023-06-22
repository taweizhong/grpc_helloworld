package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pd "grpc_helloworld/client/proto"
	"log"
	"time"
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

	stream, err := client.SayDoubleStream(context.Background())
	var count int32 = 22
	for {
		request := &pd.HelloReq{
			Name: "zzz",
			Age:  count,
		}
		err = stream.Send(request)
		count++
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		recv, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		//websocket
		fmt.Println("客户端收到的流信息", recv.Say)
	}
}
