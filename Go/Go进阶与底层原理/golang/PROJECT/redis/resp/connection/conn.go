package connection

import (
	"errors"
	"net"
	"redis/lib/sync/wait"
	"time"
)

type Connection struct {
	conn net.Conn
	//加入等待组是因为防止还没响应完给客户端，连接就被关闭
	waitingReply wait.Wait
	selectDB     int
}

func (c *Connection) Close() error {
	timeout := c.waitingReply.WaitWithTimeout(10 * time.Second)
	c.conn.Close()
	if timeout {
		return errors.New("timeout")
	}
	return nil
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	c.waitingReply.Add(1)
	defer c.waitingReply.Done()
	_, err := c.conn.Write(bytes)
	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.selectDB = dbNum
}
