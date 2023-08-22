package tcp

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type server struct {
}

const defaultRPCPath = "/rpc"

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	log.Println("test ok")
	w.Write([]byte("success"))
}

func Dial() {
	conn, _ := net.Dial("tcp", ":8888")
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", defaultRPCPath))

	resp, err := http.ReadResponse(bufio.NewReader(conn), nil)
	log.Println("readResponse finished")
	if err != nil {
		log.Println(err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)

	println(string(bytes))
}

func Server() {
	l, _ := net.Listen("tcp", ":8888")
	http.Handle(defaultRPCPath, new(server))
	http.Serve(l, nil)

}
