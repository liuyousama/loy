package tcp

import (
	"github.com/pkg/errors"
	"net"
)

type Server struct {
	Addr    string
	Network string
	router  *Router
}

func NewServer(network, addr string, router *Router) *Server {
	return &Server{addr, network, router}
}

func (s *Server) start() (err error) {
	addr, err := net.ResolveTCPAddr(s.Network, s.Addr)
	if err != nil {
		return errors.Wrap(err, "resolve tcp addr error before start the tcp server")
	}

	listener, err := net.ListenTCP(s.Network, addr)
	if err != nil {
		return errors.Wrap(err, "listen tcp addr failed")
	}

	for {
		//accept tcp connection from the listener, and wrap it into custom tcp connect
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		connect := NewConnect(conn)
		//start to serve the tcp connect
		connect.Serve()
	}

}
