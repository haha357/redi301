package http

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"net"
	"redi301/config"
	"redi301/redirect"
	"syscall"
)

func Start() {
	cfg := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				//syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				//syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_RCVBUF, 1)
				linger := syscall.Linger{
					Onoff:  1,
					Linger: 0,
				}
				syscall.SetsockoptLinger(int(fd), syscall.SOL_SOCKET, unix.SO_LINGER, &linger)
				syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, unix.TCP_NODELAY, 1)
			})
		},
	}
	listener, err := cfg.Listen(context.Background(), "tcp", config.HttpAddr)
	//listener, err := net.Listen("tcp", ":"+config.HttpPort)
	if err != nil {
		logrus.Fatalf("[http] listen fail, err: %v\n", err)
	}
	for {
		// 等待连接
		conn, err := listener.Accept()
		if err != nil {
			logrus.Debugf("[http] accept fail, err: %v\n", err)
			continue
		}
		// 对每个新连接创建一个协程进行收发数据
		go redirect.Process(conn)
	}
}
