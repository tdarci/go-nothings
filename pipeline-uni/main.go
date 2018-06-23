package main

import (
	"fmt"
	"time"
)

const maxNaturals = 5

func generator(out chan<- int) {
	for x := 0; x <= maxNaturals; x++ {
		time.Sleep(time.Second)
		out <- x
	}
	// we are done producing numbers
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for x := range in {
		out <- x * x
	}
	close(out)
}

func printer(in <-chan int) {
	for s := range in {
		fmt.Println(s)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go generator(naturals)
	go squarer(squares, naturals)
	printer(squares)
}
