package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pd "grpc_helloworld/client/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
