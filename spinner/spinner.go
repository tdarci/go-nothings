package main

import (
	"fmt"
	"time"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // nice and slow, so we can see our spinner
	fmt.Printf("\nFibonacci(%d) = %d\n\n", n, fibN)
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
