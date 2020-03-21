package main

import (
	"fmt"

	"github.com/kidlinus/neat"
)

func main() {
	client, err := neat.NewClient("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	client.Write(neat.NewBuffer().Write("Kalle ANKA"))
	msg, err := client.Read(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(msg.ReadStr())
	client.Close()
}
