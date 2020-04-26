package paxos

const (
	ProposerIdBase = 100
	AcceptorIdBase = 200
	LearnerIdBase  = 300

	ChannelLength = 2048

	SeqShift = 32

	MsgRecvTimeout  = 1
	DeniedSleepTime = 3
	PhaseTimeOut    = 10
)

var MsgLossRate float64 = 0.5
