package neat

import (
	"errors"
	"net"
	"time"
)

var (
	ErrTimeout = errors.New("timeout")
)

type Server struct {
	listener net.Listener
	ch       chan serverAccept
}

type serverAccept struct {
	err    error
	client *Client
}

func NewServer(listener net.Listener) *Server {
	s := &Server{
		listener: listener,
		ch:       make(chan serverAccept, 10),
	}
	go s.runtime()
	return s
}

func (s *Server) Accept(blocking bool) (*Client, error) {
	if blocking {
		v := <-s.ch
		return v.client, v.err
	}
	select {
	case v := <-s.ch:
		return v.client, v.err
	default:
		return nil, nil
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) runtime() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.ch <- serverAccept{err: err}
			return
		}
		c := &Client{
			conn: conn,
			ch:   make(chan clientMessage, 10),
		}
		go c.runtime()
		s.ch <- serverAccept{client: c}
	}
}

type Client struct {
	conn net.Conn
	ch   chan clientMessage
}

func NewClient(conn net.Conn) *Client {
	c := &Client{
		conn: conn,
		ch:   make(chan clientMessage, 10),
	}
	go c.runtime()
	return c
}

type clientMessage struct {
	err  error
	data *BufferReadable
}

func (c *Client) runtime() {
	reader := make([]byte, 8192)
	buffer := []byte{}
	for {
		c.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		l, err := c.conn.Read(reader)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				c.onError(ErrTimeout)
				return
			}
			c.onError(err)
			return
		}
		buffer = append(buffer, reader[:l]...)
	parse:
		for {
			if len(buffer) < 2 {
				break parse
			}
			messageSize := 2 + int(uint16(buffer[0])|uint16(buffer[1])<<8)
			if len(buffer) < messageSize {
				break parse
			}
			c.onMessage(buffer[2:messageSize])
			if len(buffer) == messageSize {
				buffer = []byte{}
				break parse
			}
			buffer = buffer[messageSize:]
		}
	}
}

func (c *Client) onMessage(v []byte) {
	if len(v) == 0 {
		c.conn.Write([]byte{0, 0})
		return
	}
	c.ch <- clientMessage{data: &BufferReadable{Buffer: v}}
}

func (c *Client) onError(err error) {
	c.ch <- clientMessage{err: err}
}

func (c *Client) Read(blocking bool) (*BufferReadable, error) {
	if blocking {
		v := <-c.ch
		return v.data, v.err
	}
	select {
	case v := <-c.ch:
		return v.data, v.err
	default:
		return nil, nil
	}
}

func (c *Client) Write(v []byte) {
	l := uint16(len(v))
	c.conn.Write([]byte{byte(l), byte(l >> 8)})
	c.conn.Write(v)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Client) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
