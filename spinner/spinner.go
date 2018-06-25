package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond) // kick off a goroutine to show a spinner
	const startVal = 45
	fibN := fib(45) // calculate fibonacci, nice and slowly... NOT a separate routine
	fmt.Printf("\nFibonacci(%d) = %d\n\n", startVal, fibN)
}

func spinner(delay time.Duration) {
	fmt.Println("")
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c   ", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
