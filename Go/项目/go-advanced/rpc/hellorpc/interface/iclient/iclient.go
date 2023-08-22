package iclient

import (
	"go-advanced/rpc/hellorpc/interface/iserver"
	"net/rpc"
)

type HelloServiceClient struct {
	//组合的形式，初始化或者使用是用最后面的大写结构体Client
	*rpc.Client
}

func DialHelloService(network string, address string) (*HelloServiceClient, error) {
	client, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: client}, nil
}

func (c *HelloServiceClient) Hello(request string, reply *string) error {
	return c.Client.Call(iserver.HelloServiceName+".Hello", request, reply)
}
