package main

import (
	"fmt"
	"net"

	"github.com/kidlinus/neat"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	l, err := net.Listen("tcp", ":1777")
	panicOnErr(err)
	srv := neat.NewServer(l)
	for {
		c, err := srv.Accept(true)
		if err != nil {
			fmt.Println("error", err)
			break
		}
		fmt.Println("Client joined", c.RemoteAddr())
		go func() {
			for {
				msg, err := c.Read(true)
				if err != nil {
					fmt.Println("client error", err)
					return
				}
				fmt.Println("client message", msg)
			}
		}()
	}
	srv.Close()
}
