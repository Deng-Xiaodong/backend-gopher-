package handler

import (
	"context"
	"io"
	"net"
	database2 "redis/database"
	"redis/interface/database"
	"redis/lib/logger"
	"redis/lib/sync/atomic"
	"redis/resp/connection"
	"redis/resp/parser"
	"redis/resp/reply"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

//处理器维护所有连接和数据库
type RespHandler struct {
	activeConn sync.Map // *client -> placeholder
	mdb        database.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	return &RespHandler{
		mdb: database2.NewDatabase(),
	}
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	client.Close()
	h.mdb.AfterClientClose(client)
	h.activeConn.Delete(client)
}
func (h *RespHandler) Handle(ctx context.Context, conn net.Conn) {

	//1 判断处理器是否处于关闭状态，如果是则拒绝新的连接
	if h.closing.Get() {
		// closing handler refuse new connection
		conn.Close()
	}
	//2 将原始连接包装成redis连接,并将新连接加入连接字典
	println(conn.RemoteAddr())
	client := connection.NewConn(conn)
	h.activeConn.Store(client, 1)

	//3 将redis连接送给解析器，并监听数据通道

	ch := parser.ParseStream(conn)

	for payload := range ch {
		go func() {
			h.exec(payload, client)
		}()
	}

}
func (h *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// TODO: concurrent wait
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.mdb.Close()
	return nil
}
func (h *RespHandler) exec(payload *parser.Payload, client *connection.Connection) {

	if payload.Err != nil {
		if payload.Err == io.EOF ||
			payload.Err == io.ErrUnexpectedEOF ||
			strings.Contains(payload.Err.Error(), "use of closed network connection") {
			// connection closed
			h.closeClient(client)
			logger.Info("connection closed: " + client.RemoteAddr().String())
			return
		}
		// protocol err
		errReply := reply.MakeErrReply(payload.Err.Error())
		err := client.Write(errReply.ToBytes())
		if err != nil {
			h.closeClient(client)
			logger.Info("connection closed: " + client.RemoteAddr().String())
			return
		}

	}
	if payload.Data == nil {
		logger.Error("empty payload")

	}
	r, ok := payload.Data.(*reply.MultiBulkReply)
	if !ok {
		logger.Error("require multi bulk reply")

	}
	result := h.mdb.Exec(client, r.Args)
	if result != nil {
		client.Write(result.ToBytes())
	} else {
		client.Write(unknownErrReplyBytes)
	}
}
