package multi

import (
	"math/rand"
	"sync"
	"time"
)

type proposer struct {
	id        int
	leaderId  int
	seq       int
	round     int
	maxapn    int
	pn        int
	pv        int
	acceptors map[int]*message
	nt        *network

	learderTime time.Time
	prepared    bool
	proposed    bool
	requests    []int
}

func NewProposer(id int, nt *network) *proposer {
	p := new(proposer)
	p.id = id
	p.leaderId = 0
	p.seq = 0
	p.round = 0
	p.maxapn = 0
	p.pn = 0
	p.pv = 0
	p.nt = nt
	p.acceptors = map[int]*message{}

	p.learderTime = time.Now()
	p.prepared = false
	p.requests = []int{}

	for i := 0; i < nt.acceptorNum; i++ {
		p.acceptors[i+AcceptorIdBase] = nil
	}
	return p
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
				if p.leaderId == 0 || p.leaderId == p.id {
					p.requests = append(p.requests, msg.pv)
				} else {
					NewRequest(p.nt, p.leaderId, msg.pv)
				}
			case Promise:
				if !denied && msg.getPNSeq() == p.seq {
					p.checkPromise(msg)
				}
			case Nack:
				if !denied && msg.getPNSeq() > p.seq {
					p.seq = msg.getPNSeq()
					denied = true
				}
			case Accepted:
				if !denied && msg.apn == p.pn {
					p.reset(false)
				}

				if msg.apn > p.maxapn {
					p.maxapn = msg.apn
					leaderId := msg.getAPNId()
					if p.leaderId != leaderId {
						p.leaderId = leaderId
						denied = true
					}
				}
			case Halt:
				running = false
			}
		}

		if denied {
			p.reset(false)
			time.Sleep(DeniedSleepTime * time.Millisecond)
			denied = false
		}

		if p.isLearderTimeTimeout() {
			p.reset(true)
		}

		if !p.prepared && p.leaderId == 0 {
			p.prepare()
		} else if p.leaderId == p.id || p.checkMajority() {
			p.propose()
		}
	}
	wg.Done()
}

func (p *proposer) prepare() {
	maxseq := p.maxapn >> SeqShift >> RoundShift
	if maxseq > p.seq {
		p.seq = maxseq
	}
	p.seq++
	cnt := 0
	for aid := range p.acceptors {
		msg := NewPrePareMsg(p.id, aid, p.getProposeNum())
		p.nt.send(msg)
		cnt++
	}
	p.prepared = true
}

func (p *proposer) propose() {
	if p.proposed {
		return
	}

	if p.pv == 0 {
		if p.leaderId == p.id {
			if len(p.requests) > 0 {
				p.pv = p.requests[0]
				p.requests = p.requests[1:]
				p.round++
			} else {
				return
			}
		}
	}
	p.proposed = true

	p.pn = p.getProposeNum()
	cnt := 0
	for aid := range p.acceptors {
		pmsg := NewAcceptMsg(p.id, aid, p.pn, p.pv)
		p.nt.send(pmsg)
		cnt++
	}
}

func (p *proposer) isLearderTimeTimeout() bool {
	t := time.Now()
	d := t.Sub(p.learderTime)
	ms := d.Milliseconds() + rand.Int63n(500)
	return ms > LeaderTimeOut
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

	return cnt >= p.getMajority()
}

func (p *proposer) getMajority() int {
	n := len(p.acceptors)
	return n/2 + 1
}

func (p *proposer) getSeq() int {
	return p.seq<<RoundShift | p.round
}

func (p *proposer) getProposeNum() int {
	p.pn = p.getSeq()<<SeqShift | p.id
	return p.pn
}

func (p *proposer) reset(newLeader bool) {
	p.pn = 0
	p.pv = 0
	p.learderTime = time.Now()

	for aid := range p.acceptors {
		p.acceptors[aid] = nil
	}

	p.proposed = false

	if newLeader {
		p.round = 0
		p.leaderId = 0
		p.prepared = false
	}
}
