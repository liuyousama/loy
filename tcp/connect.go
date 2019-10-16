package tcp

import (
	"fmt"
	"net"
)

type connect struct {
	conn *net.TCPConn
	isClosed bool
	exitChan chan bool
	messageChan chan []byte

}

func NewConnect(conn *net.TCPConn) *connect {
	nc := &connect{
		conn:conn,
		isClosed:false,
		exitChan:make(chan bool, 1),
		messageChan:make(chan []byte, 0),
	}

	return nc
}

func (c *connect)Serve() {
	go c.startReader()
	go c.startWriter()
}

func (c *connect)startReader()  {
	defer func() {_ = c.conn.Close()}()

	for {
		msg, err := acceptMessage(c.conn)
		if err != nil {
			fmt.Println("read msg header error ", err)
			break
		}

		msg = msg
	}
}

func (c *connect)startWriter() {
	for  {
		select {
		case message := <- c.messageChan:
			_, err := c.conn.Write(message)
			if err != nil {
				fmt.Println("Send data to connection error!", err)
				continue
			}
		case <- c.exitChan:
			return

		}
	}
}