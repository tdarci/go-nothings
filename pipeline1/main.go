package main

import (
	"fmt"
	"time"
)

const maxNaturals = 5

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Number Generator
	go func() {
		for x := 0; x <= maxNaturals; x++ {
			time.Sleep(time.Second)
			naturals <- x
		}
		// we are done producing numbers
		close(naturals)
	}()

	// Squarer
	go func() {
		for x := range naturals { // when naturals closes, this loop exits
			squares <- x * x
		}
		close(squares) // this will shut down our program
	}()

	// Printer (in main goroutine)
	for s := range squares { // when squares closes, this loop exits
		fmt.Println(s)
	}
}
