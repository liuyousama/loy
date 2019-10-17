package tcp

import "net"

type Context struct {
	Value interface{}
	request *request
}

func newContext(r *request) *Context {
	return &Context{request:r}
}

func (c *Context)GetMessageData() []byte {
	return c.request.msg.data
}

func (c *Context)GetMessage() *message {
	return c.request.msg
}

func (c *Context)Write(data []byte)  {
	c.request.conn.write(data)
}

func (c *Context)SaveConnect(key interface{})  {
	if key == nil {
		return
	}
	
	saveConnect(key, c.request.conn)
}

func (c *Context)Send(key interface{}, data []byte) bool {
	connect := getConnect(key)
	if connect == nil {
		return false
	}

	connect.messageChan <- data
	return true
}

func (c *Context)Broadcast(data []byte)  {
	for _, val := range getAllConnect() {
		val.messageChan <- data
	}
}

func (c *Context)Close()  {
	key := c.request.conn.key
	if key != nil {
		delConnect(key)
	}
	
	c.request.conn.exitChan <- true
	
}

func (c *Context) Set(value interface{})  {
	c.Value = value
}

func (c *Context) Get() interface{} {
	return c.Value
}

func (c *Context) GetRemoteAddr() net.Addr {
	return c.request.conn.conn.RemoteAddr()
}