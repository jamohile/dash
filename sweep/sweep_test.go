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
			close(c)
		},
		Worker: func(c int, r chan int, manager Manager) {
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
			close(c)
		},
		Worker: func(c int, r chan int, manager Manager) {
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
