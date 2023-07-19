package iserver

import "net/rpc"

const HelloServiceName = "HelloService"

type IHelloService = interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(svc IHelloService) error {
	return rpc.RegisterName(HelloServiceName, svc)
}
