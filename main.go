package main

import (
	"bufio"
	"os"
	"sync"

	basic "github.com/myzWILLmake/paxos-go/basic"
	multi "github.com/myzWILLmake/paxos-go/multi"
)

func main() {
	var wg sync.WaitGroup
	reader := bufio.NewReader(os.Stdin)
	ch := make(chan string)
	running := true
	for running {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		switch text {
		case "basic":
			wg.Add(1)
			go basic.Run(ch, &wg)
		case "multi":
			wg.Add(1)
			go multi.Run(ch, &wg)
		case "exit":
			running = false
			ch <- text
		default:
			ch <- text
		}
	}
	wg.Wait()
}
