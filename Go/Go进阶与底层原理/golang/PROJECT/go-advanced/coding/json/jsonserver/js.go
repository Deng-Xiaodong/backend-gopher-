package jsonserver

import (
	"encoding/json"
	"fmt"
	json2 "go-advanced/coding/json"
	"log"
	"net"
)

func Accept(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := listen.Accept()
	println(conn)
	if err != nil {
		log.Fatalln(err)
	}
	item := new(json2.Item)
	if err = json.NewDecoder(conn).Decode(item); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v", item)
}
