package main

import (
	"bufio"
	"os"

	paxos "github.com/myzWILLmake/paxos-go/multi"
)

func main() {
	ch := make(chan string)
	go paxos.Run(ch)
	reader := bufio.NewReader(os.Stdin)
	for true {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		ch <- text
	}
}
