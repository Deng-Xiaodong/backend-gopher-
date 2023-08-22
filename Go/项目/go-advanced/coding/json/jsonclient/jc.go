package jsonclient

import (
	"encoding/json"
	json2 "go-advanced/coding/json"
	"log"
	"net"
)

func DialAndSend(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	item := &json2.Item{
		Age:  18,
		Name: "dong",
	}
	if err = json.NewEncoder(conn).Encode(item); err != nil {
		log.Fatalln(err)
	}
}
