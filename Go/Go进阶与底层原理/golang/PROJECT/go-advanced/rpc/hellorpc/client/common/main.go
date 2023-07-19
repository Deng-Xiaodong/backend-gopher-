package main

import (
	"go-advanced/rpc/hellorpc/interface/iclient"
	"log"
)

func main() {
	client, err := iclient.DialHelloService("tcp", ":1234")
	if err != nil {
		log.Fatalln(err)
	}
	var reply string
	err = client.Hello("dong", &reply)
	if err != nil {
		log.Fatalln(err)
	}
	println(reply)
}
