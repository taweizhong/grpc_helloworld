package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pd "grpc_helloworld/client/proto"
	"log"
)

func main() {
	// 添加证书
	file, err2 := credentials.NewClientTLSFromFile("../server/server.pem", "*.test.example.com")
	if err2 != nil {
		log.Fatal("证书错误", err2)
	}
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(file))

	if err != nil {
		fmt.Print(err)
	}
	defer conn.Close()
	client := pd.NewHelloServerClient(conn)
	feature, err := client.Say(context.Background(), &pd.HelloReq{
		Name: "hello",
		Age:  22,
	})
	if err != nil {
		log.Print(err)
	}
	fmt.Print(feature)
}
