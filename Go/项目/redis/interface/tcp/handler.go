package tcp

import (
	"context"
	"net"
)

type HandleFunc func(ctx context.Context, conn net.Conn)

//tcp请求处理接口
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
