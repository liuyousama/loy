package tcp

import "net"

type Context struct {
	Value interface{}
	request *request
}

func newContext(r *request) *Context {
	return &Context{request:r}
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