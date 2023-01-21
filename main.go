package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

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
			time.Sleep(2 * time.Second)
			cmd := exec.Command("echo", strconv.Itoa(c))
			output, _ := cmd.Output()
			fmt.Println("Done!")
			return string(output)
		},
		MaxWorkers: 5,
	}

	fmt.Println(s.Run())
}
