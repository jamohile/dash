package main

import (
	"fmt"
	"time"

	"dash/sweep"
)

func main() {
	s := sweep.Sweep[int, int]{
		Generator: func(c chan int, manager sweep.Manager) {
			for i := 0; i < 100; i++ {
				c <- i
			}
			close(c)
		},
		Worker: func(c int, r chan int, manager sweep.Manager) {
			if c >= 3 {
				fmt.Println("cancel")
				manager.Cancel()
			}
			time.Sleep(100 * time.Millisecond)
			if !manager.IsDone() {
				r <- c
			}
		},
		MaxWorkers: 1,
	}

	fmt.Println(s.Run())
}
