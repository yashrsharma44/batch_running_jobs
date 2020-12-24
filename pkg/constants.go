package pkg

// ENUM for defining possible states for a worker
type STATE int

const (
	NOT_RUNNING STATE = iota
	PLAY
	PAUSE
	STOP
	COMPLETED
)

// The total concurrent jobs we want to handle
const MAXWORKER uint64 = 2

const WORKERIDEXPIRE uint64 = 5

var MapStatetoMsg = map[STATE]string{
	NOT_RUNNING: "Not Running",
	PLAY:        "Running",
	PAUSE:       "Paused",
	STOP:        "Terminated",
	COMPLETED:   "Completed",
}

var MapParamtoState = map[string]STATE{
	"PAUSE":  PAUSE,
	"STOP":   STOP,
	"RESUME": PLAY,
}
