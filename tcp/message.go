package tcp

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"net"
)

type message struct {
	length uint32
	pathId uint32
	data   []byte
}

func (m *message)GetMessageLength() uint32 {
	return m.length
}

func (m *message)GetMessagePath() uint32 {
	return m.pathId
}

func (m *message)GetMessageData() []byte {
	return m.data
}

func acceptMessage(conn *net.TCPConn) (msg *message, err error) {
	msg = &message{}
	b := make([]byte, 8)

	//read message header from tcp connection
	_, err = io.ReadFull(conn, b)
	if err != nil {
		err = errors.Wrap(err, "read message header form tcp connection failed. ")
		return
	}
	r := bytes.NewReader(b)
	//read message length
	err = binary.Read(r, binary.LittleEndian, &msg.length)
	if err != nil {
		err = errors.Wrap(err, "read message length from buffer failed.")
	}
	//read message path id
	err = binary.Read(r, binary.LittleEndian, &msg.pathId)
	if err != nil {
		err = errors.Wrap(err, "read message path id from buffer failed.")
	}

	//read the main body of data
	data := make([]byte, msg.length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		err = errors.Wrap(err, "receive the data from tcp connection failed.")
		return
	}

	msg.data = data

	return
}