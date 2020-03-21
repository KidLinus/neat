package neat

import (
	"net"
)

type Server struct {
	listener net.Listener
	ch       chan serverAccept
}

type serverAccept struct {
	conn net.Conn
	err  error
}

func NewServer(network, address string) (*Server, error) {
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: l,
		ch:       make(chan serverAccept, 100),
	}
	go server.runtime()
	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) runtime() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			s.ch <- serverAccept{err: err}
			return
		}
		s.ch <- serverAccept{conn: c}
	}
}

func (s *Server) Accept(block bool) (*Client, error) {
	var v serverAccept
	if block {
		v = <-s.ch
	} else {
		select {
		case v = <-s.ch:
		default:
			return nil, nil
		}
	}
	if v.err != nil {
		return nil, v.err
	}
	cli := &Client{
		conn:   v.conn,
		buffer: make([]byte, 0),
		ch:     make(chan *clientPayload, 30),
	}
	go cli.runtime()
	return cli, nil
}
