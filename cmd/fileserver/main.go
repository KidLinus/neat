package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kidlinus/neat"
)

func main() {
	f, err := ioutil.ReadFile("THX-DeepNote-96Khz.wav")
	if err != nil {
		panic(err)
	}
	total := len(f)
	srv, err := neat.NewServer("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	defer srv.Close()
	for {
		cli, err := srv.Accept(true)
		if err != nil {
			panic(err)
		}
		cli.Write(neat.NewBuffer().Write("THX-DeepNote-96Khz.wav"))
		start := time.Now()
		for i := 0; i < len(f); i += 1024 {
			top := i + 1024
			if top > len(f) {
				top = len(f)
			}
			if i/1024%1024 == 0 {
				fmt.Printf("%.1f (%d MB / %d MB)\n", float64(i)/float64(total)*100, i/MB, total/MB)
			}
			cli.Write(neat.NewBuffer(f[i:top]))
		}
		cli.Close()
		end := time.Now()
		fmt.Printf("Transfer complete. time: %v, speed %.2f mbit/s\n", end.Sub(start), float64(total/MB)/end.Sub(start).Seconds())
	}
}

const (
	_      = iota // ignore first value by assigning to blank identifier
	KB int = 1 << (10 * iota)
	MB
)
