package main

import (
	"fmt"
	"strings"

	"dash/sweep"
)

func main() {
	s := sweep.Sweep[int, string]{
		Generator: func(c chan int) {
			for i := 0; i < 20; i++ {
				c <- i
			}
			close(c)
		},
		Worker: func(c int) string {
			return strings.Repeat("yo", c)
		},
		MaxWorkers: 5,
	}

	fmt.Println(s.Run())
}
