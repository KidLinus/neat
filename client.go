package neat

import (
	"net"
)

type Client struct {
	conn   net.Conn
	buffer []byte
	ch     chan *clientPayload
}

type clientPayload struct {
	buffer *BufferReadable
	err    error
}

func NewClient(network, address string) (*Client, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	cli := &Client{
		conn:   c,
		buffer: make([]byte, 0),
		ch:     make(chan *clientPayload, 30),
	}
	go cli.runtime()
	return cli, nil
}

func (c *Client) runtime() {
	b := make([]byte, 8092)
	for {
		l, err := c.conn.Read(b)
		if err != nil {
			c.ch <- &clientPayload{err: err}
			return
		}
		c.buffer = append(c.buffer, b[:l]...)
	parse:
		for {
			if len(c.buffer) < 2 {
				break parse
			}
			target := int(uint16(c.buffer[0])|uint16(c.buffer[1])<<8) + 2
			if len(c.buffer) < target {
				break parse
			}
			c.ch <- &clientPayload{buffer: NewBuffer(c.buffer[2:target]).Readable()}
			if len(c.buffer)-target == 0 {
				c.buffer = make([]byte, 0)
				break parse
			}
			c.buffer = c.buffer[target:]
		}
	}
}

func (c *Client) Read(block bool) (*BufferReadable, error) {
	if block {
		v := <-c.ch
		return v.buffer, v.err
	}
	select {
	case v := <-c.ch:
		return v.buffer, v.err
	default:
		return nil, nil
	}
}

func (c *Client) Write(data Buffer) error {
	l := len(data)
	_, err := c.conn.Write([]byte{byte(l), byte(l >> 8)})
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
