package sweep

import (
	"testing"
)

func TestRunsFullSweep(t *testing.T) {
	s := Sweep[int, int]{
		Generator: func(c chan int, manager Manager) {
			for i := 0; i < 100; i++ {
				c <- i
			}
		},
		Worker: func(c int, r chan int, e chan Event, manager Manager) {
			r <- c
		},
		MaxWorkers: 10,
	}

	results := s.Run()

	if len(results) != 100 {
		t.Fatalf("Invalid results length: %d", len(results))
	}

	// Make sure all values are in the result.
	for i := 0; i < 100; i++ {
		found := false
		for j := range results {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Value not found in results: %d", i)
			break
		}
	}
}

func TestExitsSweepEarly(t *testing.T) {
	s := Sweep[int, int]{
		Generator: func(c chan int, manager Manager) {
			for i := 0; i < 100; i++ {
				if manager.IsDone() {
					break
				}
				c <- i
			}
		},
		Worker: func(c int, r chan int, e chan Event, manager Manager) {
			if c >= 10 {
				manager.Cancel()
				return
			}
			r <- c
		},

		// If we can have more than one worker, this may be non deterministic.
		MaxWorkers: 1,
	}

	results := s.Run()

	if len(results) != 10 {
		t.Fatalf("Invalid results length: %d", len(results))
	}
}

func TestFiresEvents(t *testing.T) {
	event_map := map[int]bool{}

	s := Sweep[int, int]{
		Generator: func(c chan int, manager Manager) {
			for i := 0; i < 100; i++ {
				c <- i
			}
		},
		Worker: func(c int, r chan int, e chan Event, manager Manager) {
			r <- c
			e <- Event{
				Type: "data",
				Data: c,
			}
		},
		OnEvent: func(e WorkerEvent[int]) {
			if e.Event.Type != "data" {
				return
			}
			event_map[e.Event.Data.(int)] = true
			if e.Description.Config != e.Event.Data.(int) {
				t.Fatalf("Invalid data. Got %d, expected %d", e.Description.Config, e.Event.Data.(int))
			}
		},
		MaxWorkers: 10,
	}
	s.Run()

	// Check that all 100 events have been called.
	for i := 0; i < 100; i++ {
		_, ok := event_map[i]
		if !ok {
			t.Fatalf("Event not fired: %d", i)
			break
		}
	}
}

func TestNestedSweep(t *testing.T) {
	s := Sweep[int, []int]{
		Generator: func(c chan int, m Manager) {
			for i := 0; i < 100; i++ {
				c <- i
			}
		},
		Worker: func(c int, r chan []int, e chan Event, m Manager) {
			s1 := Sweep[int, int]{
				Generator: func(c1 chan int, m1 Manager) {
					for i := 0; i < 100; i++ {
						c1 <- i
					}
				},
				Worker: func(c1 int, r1 chan int, e chan Event, m1 Manager) {
					r1 <- c1
				},
				MaxWorkers: 1,
			}
			results := []int{}
			for _, result := range s1.Run() {
				results = append(results, result.Result)
			}
			r <- results
		},
		MaxWorkers: 1,
	}

	results := s.Run()

	if len(results) != 100 {
		t.Fatalf("Invalid results length: %d", len(results))
	}

	for i := 0; i < 100; i++ {
		subResult := results[i]
		if len(subResult.Result) != 100 {
			t.Fatalf("Invalid sub results length: %d", len(subResult.Result))
		}
	}
}
