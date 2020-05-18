package main

import (
	"net"
	"fmt"
	"time"
	"github.com/happyxcj/gosocket"
	"github.com/happyxcj/gosocket/protocol"
	"github.com/happyxcj/gosocket/route"
)

const addr = ":8080"

func main() {
	initHandlers()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("listen error: %v", err.Error()))
	}
	for {
		//等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				time.Sleep(time.Second)
				fmt.Println("net temporary error when accepting new connections: ", err.Error())
			} else {
				fmt.Println("other error when accepting new connections: ", err.Error())
			}
			continue
		}
		c := &Client{}
		ec := gosocket.NewEasyConn(conn)
		qc := gosocket.NewQueueConn(ec, gosocket.PktHandler(c.handlePacket))
		c.QueueConn = qc
	}
}

type Client struct {
	*gosocket.QueueConn
}

func (c *Client) handlePacket(p protocol.Packet) {
	route.HandlePacket(c, p)
}
