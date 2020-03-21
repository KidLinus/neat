package main

import (
	"fmt"
	"os"

	"github.com/kidlinus/neat"
)

func main() {
	cli, err := neat.NewClient("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	msg, err := cli.Read(true)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(msg.ReadStr())
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for {
		msg, err = cli.Read(true)
		if err != nil {
			fmt.Println("closed", err)
			return
		}
		_, err := f.Write(msg.Buffer)
		if err != nil {
			panic(err)
		}
	}
}
