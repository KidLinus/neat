package main

import (
	"fmt"

	"github.com/kidlinus/neat"
)

func main() {
	srv, err := neat.NewServer("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	cli, err := srv.Accept(true)
	if err != nil {
		panic(err)
	}
	msg, err := cli.Read(true)
	str := msg.ReadStr()
	fmt.Println("client handshake", str)
	cli.Write(neat.NewBuffer().Write(fmt.Sprintf("Hello, %s!", str)))
	for {
		msg, err := cli.Read(false)
		if err != nil {
			fmt.Println("client err", err)
			break
		}
		if msg == nil {
			continue
		}
		fmt.Println("client msg", len(msg.Buffer))
	}
	srv.Close()
}
