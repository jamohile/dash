package main

import (
	"fmt"
	"time"

	"dash/sweep"
)

func main() {
	s := sweep.Sweep[int, int]{
		Generator: func(c chan int, manager sweep.SweepManager) {
			for i := 0; i < 100; i++ {
				c <- i
			}
			close(c)
		},
		Worker: func(c int, manager sweep.SweepManager) int {
			if c > 3 {
				manager.Cancel()
			}
			time.Sleep(100 * time.Millisecond)
			return c
		},
		MaxWorkers: 1,
	}

	fmt.Println(s.Run())
}
