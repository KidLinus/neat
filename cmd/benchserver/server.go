package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kidlinus/neat"
)

func main() {
	srv, err := neat.NewServer("tcp", ":1111")
	if err != nil {
		panic(err)
	}
	clients := []*neat.Client{}
	var clientsL sync.RWMutex
	go func() {
		for {
			cli, err := srv.Accept(true)
			if err != nil {
				panic(err)
			}
			clientsL.Lock()
			clients = append(clients, cli)
			clientsL.Unlock()
		}
	}()
	var step int64
	go func() {
		i := 0
		for range time.NewTicker(time.Second / 30).C {
			m := neat.NewBuffer().Write("This is step", uint64(i))
			clientsL.RLock()
			for _, c := range clients {
				c.Write(m)
				msg, err := c.Read(false)
				if err != nil {
					fmt.Println("err", err)
					continue
				}
				if msg == nil {
					continue
				}
				// fmt.Println("msg", len(msg.Buffer))
			}
			clientsL.RUnlock()
			i++
			atomic.AddInt64(&step, 1)
		}
	}()
	for {
		time.Sleep(time.Second)
		fmt.Println("ops/s", atomic.SwapInt64(&step, 0))
	}
}
