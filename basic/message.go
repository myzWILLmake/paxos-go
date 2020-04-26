package paxos

type msgType int

const (
	Request msgType = 1 << iota

	Prepare
	Promise
	Accept
	Accepted
	Nack

	Print
	Halt
)

type message struct {
	t    msgType
	from int
	to   int
	pn   int
	pv   int
	apn  int
	apv  int
}

func newMessage(t msgType, from, to, pn, pv, apn, apv int) message {
	msg := message{t, from, to, pn, pv, apn, apv}
	return msg
}

func NewRequestMsg(to, v int) message {
	msg := newMessage(Request, 0, to, 0, v, 0, 0)
	return msg
}

func NewPrePareMsg(from, to, pn int) message {
	msg := newMessage(Prepare, from, to, pn, 0, 0, 0)
	return msg
}

func NewPromiseMsg(from, to, pn, apn, apv int) message {
	msg := newMessage(Promise, from, to, pn, 0, apn, apv)
	return msg
}

func NewAcceptMsg(from, to, pn, pv int) message {
	msg := newMessage(Accept, from, to, pn, pv, 0, 0)
	return msg
}

func NewAcceptedMsg(from, to, apn, apv int) message {
	msg := newMessage(Accepted, from, to, 0, 0, apn, apv)
	return msg
}

func NewNackMsg(from, to, pn int) message {
	msg := newMessage(Nack, from, to, pn, 0, 0, 0)
	return msg
}

func NewPrintMsg(to int) message {
	msg := newMessage(Print, 0, to, 0, 0, 0, 0)
	return msg
}

func NewHaltMsg(to int) message {
	msg := newMessage(Halt, 0, to, 0, 0, 0, 0)
	return msg
}

func (m message) getPNSeq() int {
	return m.pn >> SeqShift
}

func (m message) getPNId() int {
	return m.pn & (1<<SeqShift - 1)
}

func (m message) getAPNSeq() int {
	return m.apn >> SeqShift
}

func (m message) getAPNId() int {
	return m.apn & (1<<SeqShift - 1)
}
