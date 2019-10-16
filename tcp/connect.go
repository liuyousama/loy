package tcp

import (
	"fmt"
	"net"
)

type connect struct {
	conn           *net.TCPConn
	isClosed       bool
	exitChan       chan bool
	messageChan    chan []byte
	addRequestTask func(req *request)
}

func NewConnect(conn *net.TCPConn, addRequestTask func(req *request)) *connect {
	nc := &connect{
		conn:           conn,
		isClosed:       false,
		exitChan:       make(chan bool, 1),
		messageChan:    make(chan []byte, 0),
		addRequestTask: addRequestTask,
	}

	return nc
}

func (c *connect) Serve() {
	go c.startReader()
	go c.startWriter()
}

func (c *connect) write(data []byte) {
	select {
	case c.messageChan <- data:
		return
	}
}

func (c *connect) startReader() {
	defer func() { _ = c.conn.Close() }()

	for {
		msg, err := acceptMessage(c.conn)
		if err != nil {
			fmt.Println("read msg header error ", err)
			break
		}

		request := newRequesst(c, msg)

		c.addRequestTask(request)
	}
}

func (c *connect) startWriter() {
	for {
		select {
		case message := <-c.messageChan:
			_, err := c.conn.Write(message)
			if err != nil {
				fmt.Println("Send data to connection error!", err)
				continue
			}
		case <-c.exitChan:
			return

		}
	}
}
