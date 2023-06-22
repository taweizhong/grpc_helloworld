package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pd "grpc_helloworld/server/proto"
	"io/ioutil"
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
	cert, err := tls.LoadX509KeyPair("./server.pem", "./server.key")
	if err != nil {
		log.Fatal("证书读取错误", err)
	}
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../encryption/ca.crt")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})
	listen, _ := net.Listen("tcp", ":8080")
	// 创建服务并添加证书
	newServer := grpc.NewServer(grpc.Creds(creds))

	pd.RegisterHelloServerServer(newServer, new(Server))
	err = newServer.Serve(listen)
	if err != nil {
		log.Print("err:", err)
	}
}
