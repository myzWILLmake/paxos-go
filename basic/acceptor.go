package paxos

import (
	"sync"
)

type acceptor struct {
	id        int
	pn        int
	apn       int
	apv       int
	proposers []int
	learners  []int
	nt        *network
}

func NewAcceptor(id int, nt *network) *acceptor {
	a := new(acceptor)
	a.id = id
	a.pn = 0
	a.apn = 0
	a.apv = 0
	a.nt = nt
	a.proposers = []int{}
	a.learners = []int{}
	return a
}

func (a *acceptor) setProposers(proposers []int) {
	a.proposers = proposers
}

func (a *acceptor) setLearners(learners []int) {
	a.learners = learners
}

func (a *acceptor) recvMsgs() []*message {
	return a.nt.recv(a.id)
}

func (a *acceptor) run(wg *sync.WaitGroup) {
	running := true
	for running {
		msgs := a.recvMsgs()
		for _, msg := range msgs {
			// fmt.Println(msg)
			switch msg.t {
			case Prepare:
				a.promise(msg)
			case Accept:
				a.accept(msg)
			case Halt:
				running = false
			}
		}
	}
	wg.Done()
}

func (a *acceptor) promise(msg *message) {
	if msg.pn > a.pn {
		a.pn = msg.pn
		amsg := NewPromiseMsg(a.id, msg.from, msg.pn, a.apn, a.apv)
		a.nt.send(amsg)
	} else if msg.pn <= a.pn {
		amsg := NewNackMsg(a.id, msg.from, a.pn)
		a.nt.send(amsg)
	}
}

func (a *acceptor) accept(msg *message) {
	if msg.pn == a.pn && msg.pn > a.apn {

		a.apn = msg.pn
		a.apv = msg.pv

		for _, id := range a.proposers {
			amsg := NewAcceptedMsg(a.id, id, a.apn, a.apv)
			a.nt.send(amsg)
		}

		for _, id := range a.learners {
			amsg := NewAcceptedMsg(a.id, id, a.apn, a.apv)
			a.nt.send(amsg)
		}
	}
}
