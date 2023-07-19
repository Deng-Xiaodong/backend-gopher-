package main

import (
	"go-advanced/rpc/hellorpc/interface/iserver"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {

	conn, err := net.Dial("tcp", ":2345")
	if err != nil {
		log.Fatal(err)
	}

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	var reply string
	err = client.Call(iserver.HelloServiceName+".Hello", "dongdong ", &reply)
	if err != nil {
		log.Fatal(err)
	}

	println(reply)
}
