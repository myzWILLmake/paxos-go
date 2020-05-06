package multi

import (
	"fmt"
	"sort"
	"sync"
)

type result struct {
	pn    int
	seq   int
	round int
	id    int
	pv    int
}

type learner struct {
	id      int
	nt      *network
	results map[int]*result
}

func NewLearner(id int, nt *network) *learner {
	l := new(learner)
	l.id = id
	l.nt = nt
	l.results = map[int]*result{}
	return l
}

func (l *learner) recvMsgs() []*message {
	return l.nt.recv(l.id)
}

func (l *learner) run(wg *sync.WaitGroup) {
	running := true
	for running {
		msgs := l.recvMsgs()
		for _, msg := range msgs {
			switch msg.t {
			case Accepted:
				l.getResult(msg)
			case Print:
				l.printResults()
			case Halt:
				running = false
			}
		}
	}
	wg.Done()
}

func (l *learner) getResult(msg *message) {
	r := new(result)
	r.id = msg.getAPNId()
	r.seq = msg.getAPNSeq()
	r.round = msg.getAPNRound()
	r.pn = msg.apn
	r.pv = msg.apv
	seqround := r.seq<<RoundShift | r.round
	if l.results[seqround] == nil {
		fmt.Println("Learner: seq", r.seq, "round", r.round, "proposerId", r.id, "value", r.pv)
	}
	l.results[seqround] = r
}

func (l *learner) printResults() {
	fmt.Println("Learner's Results:")
	fmt.Printf("     seq      id      pv\n")

	keys := []int{}
	for k := range l.results {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		r := l.results[k]
		fmt.Printf("%8d%8d%8d\n", r.seq, r.id, r.pv)
	}
	fmt.Println()
}
