package sweep

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type Sweep[C any, R any] struct {
	Generator  func(chan C)
	Worker     func(config C) R
	MaxWorkers int
}

/** Dispatches configurations to ready workers. **/
func (s Sweep[C, R]) dispatcher(configs chan C, results chan R) {
	sem := semaphore.NewWeighted(int64(s.MaxWorkers))

	for config := range configs {
		sem.Acquire(context.TODO(), 1)

		go func(config C) {
			results <- s.Worker(config)
			sem.Release(1)
		}(config)
	}

	// When all workers are complete, close the results channel.
	sem.Acquire(context.TODO(), int64(s.MaxWorkers))
	close(results)
}

/** Collects and buffers results from workers **/
func (s Sweep[C, R]) collector(results chan R) []R {
	var results_buffer []R
	for r := range results {
		results_buffer = append(results_buffer, r)
	}
	return results_buffer
}

/** Complete all generated work units in parallel. **/
func (s Sweep[C, R]) Run() []R {
	configs := make(chan C, 100)
	results := make(chan R, 100)

	go s.Generator(configs)
	go s.dispatcher(configs, results)
	return s.collector(results)
}
