package sweep

type WorkerDescription[C any] struct {
	ID     int
	Config C
	Name   string
}

type WorkerResult[C any, R any] struct {
	Description WorkerDescription[C]
	Result      R
}

type Event struct {
	Key  string
	Data interface{}
}

type WorkerEvent[C any] struct {
	Description WorkerDescription[C]
	Event       Event
}
