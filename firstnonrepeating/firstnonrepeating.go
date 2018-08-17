package main

import (
	"fmt"
	"os"
)

func main() {

	in := os.Args[1]
	foo := firstNonRepeating(in)
	if foo == nil {
		fmt.Printf("NO First non repeating for %s\n", in)
	} else {
		fmt.Printf("First non repeating of %s is %#U\n", in, *foo)
	}

}

func firstNonRepeating(in string) *rune {

	charMap := make(map[rune]int)
	for _, c := range in {
		charMap[c]++
	}

	for _, c := range in {
		if charMap[c] == 1 {
			return &c
		}
	}

	return nil
}
