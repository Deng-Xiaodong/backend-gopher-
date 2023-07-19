package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"redis/interface/tcp"
	"redis/lib/logger"
	"sync"
	"syscall"
	"time"
)

//tcp服务配置
type TcpConfig struct {
	Address string
	MaxConn uint32
	Timeout time.Duration
}

//一个接收系统中断信号的tcp服务
// ListenAndServeWithSignal binds port and handle requests, blocking until receive stop signal
func ListenAndServerWithSignal(cfg *TcpConfig, handler tcp.Handler) error {
	//开启两个协程，一个监听系统中断信号，一个监听关闭信号
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}

		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s start listening...", cfg.Address))
	ListenAndServer(listener, handler, closeChan)
	return nil
}

func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan chan struct{}) {
	//开启协程去负责监听cloasechan，并在接到关闭信号后关闭资源
	go func() {
		<-closeChan
		logger.Info("shutting down")
		_ = listener.Close()
		_ = handler.Close()
	}()

	//除了信号中断，可能还有一些未知错误导致没有关闭资源，所有defer recover 保证一定做了关闭
	defer func() {
		// close during unexpected error
		_ = listener.Close()
		_ = handler.Close()
	}()

	//开始正常tcp连接服务
	ctx := context.Background()
	//为了防止一些未知错误导致主程序退出，而用户连接还没有完成业务就被迫退出，所有使用了等待组
	var wg sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		logger.Info(fmt.Sprintf("accept ip: " + conn.RemoteAddr().String()))
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	wg.Wait()
}
