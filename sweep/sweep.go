package sweep

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Sweep[C any, R any] struct {
	Generator     func(config chan C, manager Manager)
	Worker        func(config C, results chan R, events chan Event, manager Manager)
	OnEvent       func(event WorkerEvent[C])
	GetWorkerName func(config C) string
	MaxWorkers    int
}

/** Generates configurations for workers. **/
func (s Sweep[C, R]) generate(configs chan C, manager Manager) {
	s.Generator(configs, manager)
	close(configs)
}

/** Dispatches configurations to ready workers. **/
func (s Sweep[C, R]) dispatch(configs chan C, results chan WorkerResult[C, R], events chan WorkerEvent[C], manager Manager) {
	sem := semaphore.NewWeighted(int64(s.MaxWorkers))

	// When all workers are complete, close the results channel.
	defer func() {
		sem.Acquire(context.TODO(), int64(s.MaxWorkers))
		close(results)
	}()

	// As long as there are configurations for workers, keep spinning up workers.
	// Uniquely identify each worker with an index.
	workerId := 0
	for config := range configs {
		if manager.IsDone() {
			return
		}

		sem.Acquire(context.Background(), 1)

		description := WorkerDescription[C]{
			ID:     workerId,
			Config: config,
		}
		if s.GetWorkerName != nil {
			description.Name = s.GetWorkerName(config)
		}

		go func(config C) {
			worker_results := make(chan R, 100)
			worker_events := make(chan Event)

			go func() {
				for worker_event := range worker_events {
					event := WorkerEvent[C]{
						Description: description,
						Event:       worker_event,
					}
					events <- event
				}
			}()

			worker_events <- Event{Type: WORKER_STARTED}
			s.Worker(config, worker_results, worker_events, manager.Child())
			worker_events <- Event{Type: WORKER_COMPLETED}

			close(worker_results)
			close(worker_events)

			for worker_result := range worker_results {
				results <- WorkerResult[C, R]{
					Description: description,
					Result:      worker_result,
				}
			}
			sem.Release(1)
		}(config)
	}
}

/** Collects and buffers results from workers **/
func (s Sweep[C, R]) collect(results chan WorkerResult[C, R]) []WorkerResult[C, R] {
	var results_buffer []WorkerResult[C, R]
	for r := range results {
		results_buffer = append(results_buffer, r)
	}
	return results_buffer
}

/** Complete all generated work units in parallel. **/
func (s Sweep[C, R]) Run() []WorkerResult[C, R] {
	configs := make(chan C, s.MaxWorkers)
	results := make(chan WorkerResult[C, R], 1000)
	events := make(chan WorkerEvent[C], 1000)
	manager := CreateManager()

	go s.generate(configs, manager)
	go s.dispatch(configs, results, events, manager)
	go func() {
		for e := range events {
			if s.OnEvent != nil {
				s.OnEvent(e)
			}
		}
	}()
	return s.collect(results)
}
