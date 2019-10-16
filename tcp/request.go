package tcp

type request struct {
	conn *connect
	msg  *message
}

func newRequesst(conn *connect, msg *message) *request {
	return &request{conn, msg}
}

func (r *request)GetTcpConnection() *connect {
	return r.conn
}

func (r *request)GetRequestMsg() *message {
	return r.msg
}

func (r *request)GetRequestData() []byte {
	return r.msg.data
}
