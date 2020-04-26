package paxos

import (
	"math/rand"
	"time"
)

type network struct {
	ch map[int]chan message
}

func NewNetwork(proposerNum, acceptorNum, learnerNum int) *network {
	nt := new(network)
	nt.ch = map[int]chan message{}

	for i := 0; i < proposerNum; i++ {
		id := ProposerIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
	}

	for i := 0; i < acceptorNum; i++ {
		id := AcceptorIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
	}

	for i := 0; i < learnerNum; i++ {
		id := LearnerIdBase + i
		nt.ch[id] = make(chan message, ChannelLength)
	}

	return nt
}

func (nt *network) setChannel(id int, ch chan message) {
	nt.ch[id] = ch
}

func (nt *network) send(msg message) {
	if msg.t&Request|msg.t&Print|msg.t&Halt == 0 {
		r := rand.Float64()
		if r <= MsgLossRate {
			return
		}
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
