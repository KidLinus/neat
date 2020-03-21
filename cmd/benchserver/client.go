package main

import (
	"fmt"
	"time"

	"github.com/kidlinus/neat"
)

func main() {
	for i := 0; i < 3000; i++ {
		spawnClient(i)
	}
	fmt.Println("Spawned all")
	for {
		time.Sleep(time.Second)
	}
}

func spawnClient(ii int) {
	cli, err := neat.NewClient("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	i := 0
	go func() {
		for range time.NewTicker(time.Second / 30).C {
			cli.Read(false)
			cli.Write(neat.NewBuffer().Write("Im here:", uint64(i)))
			i++
		}
	}()
}
