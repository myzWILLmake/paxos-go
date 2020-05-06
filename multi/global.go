package multi

const (
	ProposerIdBase = 100
	AcceptorIdBase = 200
	LearnerIdBase  = 300

	ChannelLength = 2048

	SeqShift = 32

	MsgRecvTimeout  = 5
	DeniedSleepTime = 300
	PhaseTimeOut    = 1000
)

var MsgLossRate float64 = 0
