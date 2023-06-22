package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pd "grpc_helloworld/client/proto"
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

	stream, err := client.SayStream(context.Background())
	err = stream.Send(&pd.HelloReq{
		Name: "111",
		Age:  3,
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	stream.Send(&pd.HelloReq{
		Name: "111",
		Age:  4,
	})
	stream.Send(&pd.HelloReq{
		Name: "111",
		Age:  5,
	})
	recv, err := stream.CloseAndRecv()
	if err != nil {
		return
	}
	fmt.Printf("接受：%s", recv.Say)
}
