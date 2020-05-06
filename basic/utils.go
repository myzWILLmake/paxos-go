package paxos

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Init(proposerNum, acceptorNum, learnerNum int, wg *sync.WaitGroup) *network {
	nt := NewNetwork(proposerNum, acceptorNum, learnerNum)

	for i := 0; i < proposerNum; i++ {
		p := NewProposer(ProposerIdBase+i, nt)
		wg.Add(1)
		go p.run(wg)
	}

	for i := 0; i < acceptorNum; i++ {
		a := NewAcceptor(AcceptorIdBase+i, nt)
		wg.Add(1)
		go a.run(wg)
	}

	for i := 0; i < learnerNum; i++ {
		l := NewLearner(LearnerIdBase+i, nt)
		wg.Add(1)
		go l.run(wg)
	}

	return nt
}

func NewRequest(nt *network, proposerId, val int) {
	msg := NewRequestMsg(proposerId, val)
	nt.send(msg)
}

func NewPrint(nt *network, learnerId int) {
	msg := NewPrintMsg(learnerId)
	nt.send(msg)
}

func MultiRequest(nt *network, valStart, num, interval int) {
	for i := 0; i < num; i++ {
		r := rand.Intn(nt.proposerNum)
		NewRequest(nt, r+ProposerIdBase, i+valStart)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func StopOne(nt *network, id int) {
	msg := NewHaltMsg(id)
	nt.send(msg)
	nt.setStat(id, false)
}

func StopAll(nt *network) {
	msg := NewHaltMsg(100)
	for i := 0; i < nt.proposerNum; i++ {
		msg.to = i + ProposerIdBase
		nt.send(msg)
	}

	for i := 0; i < nt.acceptorNum; i++ {
		msg.to = i + AcceptorIdBase
		nt.send(msg)
	}

	for i := 0; i < nt.learnerNum; i++ {
		msg.to = i + LearnerIdBase
		nt.send(msg)
	}
}

func RestartOne(nt *network, id int, wg *sync.WaitGroup) {
	role := id - id%100
	switch role {
	case ProposerIdBase:
		p := NewProposer(id, nt)
		wg.Add(1)
		go p.run(wg)
	case AcceptorIdBase:
		a := NewAcceptor(id, nt)
		wg.Add(1)
		go a.run(wg)
	case LearnerIdBase:
		l := NewLearner(id, nt)
		wg.Add(1)
		go l.run(wg)
	}
	nt.setStat(id, true)
}

func Run(ch chan string, pWg *sync.WaitGroup) {
	var wg sync.WaitGroup
	var nt *network

	running := true

	for running {
		var msg string
		select {
		case msg = <-ch:
		case <-time.After(100 * time.Millisecond):
		}

		if msg != "" {
			args := strings.Split(msg, " ")
			switch args[0] {
			case "init":
				proposerNum, _ := strconv.Atoi(args[1])
				acceptorNum, _ := strconv.Atoi(args[2])
				learnerNum, _ := strconv.Atoi(args[3])
				fmt.Println("Init with", proposerNum, "proposers,",
					acceptorNum, "acceptors,", learnerNum, "learners.")
				nt = Init(proposerNum, acceptorNum, learnerNum, &wg)
			case "r":
				proposerId, _ := strconv.Atoi(args[1])
				val, _ := strconv.Atoi(args[2])
				fmt.Println("Request to proposer", proposerId, "with value", val)
				NewRequest(nt, proposerId, val)
			case "mr":
				num, _ := strconv.Atoi(args[1])
				valStart, _ := strconv.Atoi(args[2])
				interval, _ := strconv.Atoi(args[3])
				fmt.Println("Request", num, "requests starting at", valStart, "with interval", interval, "ms")
				go MultiRequest(nt, valStart, num, interval)
			case "stop":
				id, _ := strconv.Atoi(args[1])
				if len(args) == 3 {
					id2, _ := strconv.Atoi(args[2])
					for i := id; id <= id2; id++ {
						StopOne(nt, i)
					}
				} else {
					StopOne(nt, id)
				}
				fmt.Println("ok")
			case "restart":
				id, _ := strconv.Atoi(args[1])
				RestartOne(nt, id, &wg)
				fmt.Println("ok")
			case "exit":
				fmt.Println("Exit...")
				StopAll(nt)
				running = false
			}
		}
	}

	wg.Wait()
	pWg.Done()
}
