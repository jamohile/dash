package sweep

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

type Sweep[C any, R any] struct {
	Generator  func(config chan C, manager Manager)
	Worker     func(config C, results chan R, manager Manager)
	MaxWorkers int
}

/** Dispatches configurations to ready workers. **/
func (s Sweep[C, R]) dispatcher(configs chan C, results chan R, manager Manager) {
	sem := semaphore.NewWeighted(int64(s.MaxWorkers))

	// When all workers are complete, close the results channel.
	defer func() {
		sem.Acquire(context.TODO(), int64(s.MaxWorkers))
		close(results)
	}()

	// As long as there are configurations for workers, keep spinning up workers.
	for config := range configs {
		if manager.IsDone() {
			return
		}

		sem.Acquire(context.Background(), 1)

		go func(config C) {
			fmt.Println("dispatch")
			fmt.Println(config)
			s.Worker(config, results, manager.Child())
			sem.Release(1)
		}(config)
	}
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
	configs := make(chan C, s.MaxWorkers)
	results := make(chan R, 1000)
	manager := CreateManager()

	go s.Generator(configs, manager)
	go s.dispatcher(configs, results, manager)
	return s.collector(results)
}
