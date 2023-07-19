package main

import (
	"go-advanced/rpc/hellorpc/interface/iserver"
	"go-advanced/rpc/hellorpc/server/helloservice"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {

	iserver.RegisterHelloService(new(helloservice.HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		conn := struct {
			io.Writer
			io.ReadCloser
		}{
			Writer:     w,
			ReadCloser: r.Body,
		}

		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":3456", nil)
}
