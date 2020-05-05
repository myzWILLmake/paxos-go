package paxos

import (
	"math/rand"
	"sync"
	"time"
)

type network struct {
	proposerNum int
	acceptorNum int
	learnerNum  int
	ch          map[int]chan message
	stat        map[int]bool
	rw          *sync.RWMutex
}

func NewNetwork(proposerNum, acceptorNum, learnerNum int) *network {
	nt := new(network)
	nt.proposerNum = proposerNum
	nt.acceptorNum = acceptorNum
	nt.learnerNum = learnerNum
	nt.ch = map[int]chan message{}
	nt.stat = map[int]bool{}
	nt.rw = new(sync.RWMutex)

	for i := 0; i < proposerNum; i++ {
		id := ProposerIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
		nt.stat[id] = true
	}

	for i := 0; i < acceptorNum; i++ {
		id := AcceptorIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
		nt.stat[id] = true
	}

	for i := 0; i < learnerNum; i++ {
		id := LearnerIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
		nt.stat[id] = true
	}

	return nt
}

func (nt *network) setChannel(id int, ch chan message) {
	nt.ch[id] = ch
}

func (nt *network) setStat(id int, val bool) {
	nt.rw.Lock()
	defer nt.rw.Unlock()

	nt.stat[id] = val
}

func (nt *network) send(msg message) {
	if msg.t&Request|msg.t&Print|msg.t&Halt == 0 {
		r := rand.Float64()
		if r <= MsgLossRate {
			return
		}
	}

	nt.rw.RLock()
	defer nt.rw.RUnlock()

	if !nt.stat[msg.to] {
		return
	}

	nt.ch[msg.to] <- msg

}

func (nt *network) recv(id int) []*message {
	msgs := []*message{}

	haveMsg := true
	for haveMsg {
		select {
		case msg := <-nt.ch[id]:
			msgs = append(msgs, &msg)
		case <-time.After(MsgRecvTimeout * time.Millisecond):
			haveMsg = false
		}
	}

	return msgs
}
