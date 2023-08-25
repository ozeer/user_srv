package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"user_srv/handler"
	"user_srv/proto"

	"google.golang.org/grpc"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "ip地址")
	port := flag.Int("port", 8088, "端口号")
	flag.Parse()

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	tcp := fmt.Sprintf("%s:%d", *ip, *port)
	lis, err := net.Listen("tcp", tcp)

	log.Printf("Start rpc server! Listen: %s", tcp)

	if err != nil {
		panic("fail to listen: " + err.Error())
	}

	err = server.Serve(lis)

	if err != nil {
		panic("fail to start grpc: " + err.Error())
	}
}
