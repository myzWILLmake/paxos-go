package paxos

import (
	"math/rand"
	"sync"
)

func Init(proposerNum, acceptorNum, learnerNum int, wg *sync.WaitGroup) (*network, []int, []int, []int) {
	nt := NewNetwork(proposerNum, acceptorNum, learnerNum)

	proposerIds := make([]int, proposerNum)
	acceptorIds := make([]int, acceptorNum)
	learnerIds := make([]int, learnerNum)

	for i := 0; i < proposerNum; i++ {
		proposerIds[i] = ProposerIdBase + i
	}
	for i := 0; i < acceptorNum; i++ {
		acceptorIds[i] = AcceptorIdBase + i
	}
	for i := 0; i < learnerNum; i++ {
		learnerIds[i] = LearnerIdBase + i
	}

	for i := 0; i < proposerNum; i++ {
		p := NewProposer(ProposerIdBase+i, nt)
		p.setAcceptors(acceptorIds)
		wg.Add(1)
		go p.run(wg)
	}

	for i := 0; i < acceptorNum; i++ {
		a := NewAcceptor(AcceptorIdBase+i, nt)
		a.setLearners(learnerIds)
		a.setProposers(proposerIds)
		wg.Add(1)
		go a.run(wg)
	}

	for i := 0; i < learnerNum; i++ {
		l := NewLearner(LearnerIdBase+i, nt)
		wg.Add(1)
		go l.run(wg)
	}

	return nt, proposerIds, acceptorIds, learnerIds
}

func NewRequests(nt *network, vals []int, proposerIds []int) {
	n := len(proposerIds)
	id := proposerIds[rand.Int()%n]
	msg := NewRequestMsg(id, 0)
	for i := 0; i < len(vals); i++ {
		msg.pv = vals[i]
		nt.send(msg)
	}
}

func NewRequest(nt *network, val int, proposerIds []int) {
	n := len(proposerIds)
	id := proposerIds[rand.Int()%n]
	msg := NewRequestMsg(id, val)
	nt.send(msg)
}

func NewBroadcastRequest(nt *network, val int, proposerIds []int) {
	n := len(proposerIds)
	msg := NewRequestMsg(0, val)
	for i := 0; i < n; i++ {
		msg.to = proposerIds[i]
		nt.send(msg)
	}
}

func NewPrint(nt *network, learnerIds []int) {
	n := len(learnerIds)
	id := learnerIds[rand.Int()%n]
	msg := NewPrintMsg(id)
	nt.send(msg)
}

func SetMsgLossRate(val float64) {
	if val > 1 {
		MsgLossRate = 1
	} else if val < 0 {
		MsgLossRate = 0
	} else {
		MsgLossRate = val
	}
}
