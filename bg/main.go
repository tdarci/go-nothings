package main

import "fmt"
import "context"
// import "golang.org/x/net/context"

func main() {
	fmt.Println("I am getting the background context...")
	c := context.Background()
	fmt.Printf("Context: %v\n", c)
	fmt.Println("Done.")
}

