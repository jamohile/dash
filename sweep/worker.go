package sweep

import "time"

type WorkerDescription[C any] struct {
	ID        int
	Config    C
	Name      string
	StartTime time.Time
	EndTime   time.Time
}

type WorkerResult[C any, R any] struct {
	Description WorkerDescription[C]
	Result      R
	Time        time.Time
}

type Event struct {
	Type string
	Data interface{}
}

type WorkerEvent[C any] struct {
	Description WorkerDescription[C]
	Event       Event
	Time        time.Time
}

const (
	WORKER_STARTED   = "dash_started"
	WORKER_COMPLETED = "dash_complete"
)
