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
	c, err := net.Dial("tcp", "192.168.0.105:1777")
	panicOnErr(err)
	client := neat.NewClient(c)
	for {
		msg, err := client.Read(true)
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		fmt.Println("Got message!", string(msg))
	}
	client.Close()
}
