package paxos

import (
	"sync"
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

func StopOne(nt *network, id int) {
	msg := NewHaltMsg(id)
	nt.send(msg)
	nt.setStat(id, false)
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
