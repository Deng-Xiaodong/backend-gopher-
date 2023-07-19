package main

import (
	"go-advanced/rpc/hellorpc/interface/iserver"
	"go-advanced/rpc/hellorpc/server/helloservice"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {

	iserver.RegisterHelloService(new(helloservice.HelloService))

	listen, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("Accept error", err)
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
