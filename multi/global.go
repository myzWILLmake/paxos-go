package multi

const (
	ProposerIdBase = 100
	AcceptorIdBase = 200
	LearnerIdBase  = 300

	ChannelLength = 2048

	SeqShift   = 32
	RoundShift = 16

	MsgRecvTimeout  = 5
	DeniedSleepTime = 300
	LeaderTimeOut   = 10000
)

var MsgLossRate float64 = 0
