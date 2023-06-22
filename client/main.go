package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pd "grpc_helloworld/client/proto"
	"io/ioutil"
	"log"
)

func main() {
	// 添加证书
	cert, _ := tls.LoadX509KeyPair("./client.pem", "./client.key")
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("../encryption/ca.crt")
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ServerName: "*.test.example.com",
		RootCAs:    certPool,
	})
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(creds))

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
