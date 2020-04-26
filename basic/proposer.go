package paxos

import (
	"sync"
	"time"
)

type proposer struct {
	id        int
	seq       int
	maxapn    int
	pn        int
	pv        int
	acceptors map[int]*message
	nt        *network

	phaseTime time.Time
	prepared  bool
	requests  []int
}

func NewProposer(id int, nt *network) *proposer {
	p := new(proposer)
	p.id = id
	p.seq = 0
	p.maxapn = 0
	p.pn = 0
	p.pv = 0
	p.nt = nt
	p.acceptors = map[int]*message{}

	p.phaseTime = time.Now()
	p.prepared = false
	p.requests = []int{}
	return p
}

func (p *proposer) setAcceptors(acceptorIds []int) {
	for _, id := range acceptorIds {
		p.acceptors[id] = nil
	}
}

func (p *proposer) recvMsgs() []*message {
	return p.nt.recv(p.id)
}

func (p *proposer) run(wg *sync.WaitGroup) {
	running := true
	for running {
		msgs := p.recvMsgs()
		denied := false
		for _, msg := range msgs {
			switch msg.t {
			case Request:
				p.requests = append(p.requests, msg.pv)
			case Promise:
				if !denied && msg.pn == p.pn {
					p.checkPromise(msg)
				}
			case Nack:
				if !denied && msg.pn > p.pn {
					p.seq = msg.getPNSeq()
					denied = true
				}
			case Accepted:
				if !denied && msg.apn > p.maxapn {
					p.reset()
					p.maxapn = msg.apn
					if len(p.requests) > 0 && msg.apv == p.requests[0] {
						p.requests = p.requests[1:]
					}
				}
			case Halt:
				running = false
			}
		}

		if denied || p.isPhaseTimeout() {
			p.reset()
			if denied {
				time.Sleep(DeniedSleepTime * time.Millisecond)
				denied = false
			}
		}

		if !p.prepared {
			if len(p.requests) > 0 {
				p.prepare()
			}
		}

		if p.prepared && p.checkMajority() {
			p.propose()
		}
	}
	wg.Done()
}

func (p *proposer) prepare() {
	p.seq++
	cnt := 0
	for aid := range p.acceptors {
		msg := NewPrePareMsg(p.id, aid, p.getProposeNum())
		p.nt.send(msg)
		cnt++

		// if cnt >= p.getMajority() {
		// 	break
		// }
	}
	p.prepared = true
}

func (p *proposer) propose() {
	if p.pv == 0 {
		p.pv = p.requests[0]
	}

	cnt := 0
	for aid, msg := range p.acceptors {
		if msg != nil {
			pmsg := NewAcceptMsg(p.id, aid, p.pn, p.pv)
			p.nt.send(pmsg)
			cnt++

			// if cnt >= p.getMajority() {
			// 	break
			// }
		}
	}
}

func (p *proposer) isPhaseTimeout() bool {
	t := time.Now()
	d := t.Sub(p.phaseTime)
	return d.Milliseconds() > PhaseTimeOut
}

func (p *proposer) checkPromise(msg *message) {
	prevmsg := p.acceptors[msg.from]
	if prevmsg == nil || msg.apn > prevmsg.apn {
		p.acceptors[msg.from] = msg
		if msg.apn > p.maxapn {
			p.maxapn = msg.apn
			p.pv = msg.apv
		}
	}
}

func (p *proposer) checkMajority() bool {
	cnt := 0
	for _, msg := range p.acceptors {
		if msg != nil {
			cnt++
		}
	}

	return cnt > p.getMajority()
}

func (p *proposer) getMajority() int {
	n := len(p.acceptors)
	return n/2 + 1
}

func (p *proposer) getProposeNum() int {
	p.pn = p.seq<<SeqShift | p.id
	return p.pn
}

func (p *proposer) reset() {
	p.pn = 0
	p.pv = 0
	p.phaseTime = time.Now()

	for aid := range p.acceptors {
		p.acceptors[aid] = nil
	}

	p.prepared = false
}
